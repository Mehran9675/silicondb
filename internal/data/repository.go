// Package data downloads the chip database and keeps an offline copy.
package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mehran9675/silicondb/internal/model"
)

// IndexPath is the database's table of contents.
const IndexPath = "index.json"

// parallelDownloads bounds concurrent requests during a sync.
const parallelDownloads = 6

// SupportedFormat is the newest index format this app understands.
const SupportedFormat = 1

// ErrNoCache means nothing has been downloaded yet.
var ErrNoCache = errors.New("no offline copy yet")

// SyncResult summarizes a sync.
type SyncResult struct {
	// Files is how many files were downloaded; 0 means the offline copy
	// was already current.
	Files int
	// Bytes is how much was downloaded, the index included.
	Bytes int
}

// Repository loads the database from the offline cache and brings it up
// to date from the remote source, downloading only changed files.
type Repository struct {
	cache  *Cache
	client *Client
}

// NewRepository returns a repository caching in dir and downloading from
// base.
func NewRepository(dir, base string) *Repository {
	return &Repository{cache: OpenCache(dir), client: NewClient(base)}
}

// SetSource changes where the database is downloaded from.
func (r *Repository) SetSource(base string) { r.client = NewClient(base) }

// LoadCached assembles the catalog from the offline copy.
func (r *Repository) LoadCached() (*model.Catalog, error) {
	b, err := r.cache.Read(IndexPath)
	if err != nil {
		return nil, ErrNoCache
	}
	var ix model.Index
	if err := json.Unmarshal(b, &ix); err != nil {
		return nil, fmt.Errorf("offline index: %w", err)
	}
	return r.assemble(&ix)
}

// Sync downloads the remote index and every file whose revision differs
// from the offline copy, then returns the updated catalog. progress, if
// not nil, is called from the downloading goroutines.
func (r *Repository) Sync(ctx context.Context, progress func(done, total int)) (*model.Catalog, SyncResult, error) {
	var res SyncResult
	// The index changes on every database update, so bypass caches.
	raw, err := r.client.Fetch(ctx, IndexPath, strconv.FormatInt(time.Now().Unix()/60, 10))
	if err != nil {
		return nil, res, err
	}
	res.Bytes += len(raw)
	var ix model.Index
	if err := json.Unmarshal(raw, &ix); err != nil {
		return nil, res, fmt.Errorf("index: %w", err)
	}
	if ix.Format > SupportedFormat {
		return nil, res, fmt.Errorf("the database needs a newer version of the app (format %d)", ix.Format)
	}

	files := ix.Files()
	var stale []string
	for path, rev := range files {
		if r.cache.Rev(path) != rev {
			stale = append(stale, path)
		}
	}

	var (
		mu       sync.Mutex
		firstErr error
		done     int
		wg       sync.WaitGroup
		slots    = make(chan struct{}, parallelDownloads)
	)
	if progress != nil {
		progress(0, len(stale))
	}
	for _, path := range stale {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			slots <- struct{}{}
			defer func() { <-slots }()
			b, err := r.client.Fetch(ctx, path, files[path])
			if err == nil {
				err = r.cache.Write(path, b)
			}
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			done++
			res.Files++
			res.Bytes += len(b)
			if progress != nil {
				progress(done, len(stale))
			}
		}(path)
	}
	wg.Wait()
	if firstErr != nil {
		return nil, res, firstErr
	}
	// The index is written last, so an interrupted sync leaves the
	// previous, consistent offline copy in place.
	if err := r.cache.Write(IndexPath, raw); err != nil {
		return nil, res, err
	}
	r.cache.Retain(files)
	cat, err := r.assemble(&ix)
	return cat, res, err
}

// ClearCache deletes the offline copy.
func (r *Repository) ClearCache() error { return r.cache.Clear() }

// assemble builds the catalog for ix from cached files. A shard that is
// missing or unreadable is skipped rather than failing the whole catalog.
func (r *Repository) assemble(ix *model.Index) (*model.Catalog, error) {
	b, err := r.cache.Read(ix.Schema.Path)
	if err != nil {
		return nil, fmt.Errorf("schema: %w", err)
	}
	schema := new(model.Schema)
	if err := json.Unmarshal(b, schema); err != nil {
		return nil, fmt.Errorf("schema: %w", err)
	}
	makers := make([]*model.Manufacturer, 0, len(ix.Manufacturers))
	var shards []*model.Shard
	for i := range ix.Manufacturers {
		e := &ix.Manufacturers[i]
		m := e.Manufacturer
		makers = append(makers, &m)
		for _, ref := range e.Shards {
			b, err := r.cache.Read(ref.Path)
			if err != nil {
				continue
			}
			s := new(model.Shard)
			if json.Unmarshal(b, s) != nil {
				continue
			}
			if s.Manufacturer == "" {
				s.Manufacturer = m.ID
			}
			if s.Category == "" {
				s.Category = ref.Category
			}
			shards = append(shards, s)
		}
	}
	return model.NewCatalog(schema, ix.Updated, makers, shards), nil
}
