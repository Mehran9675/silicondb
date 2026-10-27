package data

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// revsFile records the revision of every cached file.
const revsFile = "revs.json"

// Cache keeps downloaded files on disk, mirroring the remote layout. With
// an empty directory it keeps them in memory only, for platforms without
// a writable data directory.
type Cache struct {
	dir string

	mu    sync.Mutex
	revs  map[string]string
	files map[string][]byte // memory-only mode
}

// OpenCache opens the cache in dir, creating it if needed.
func OpenCache(dir string) *Cache {
	c := &Cache{dir: dir, revs: map[string]string{}, files: map[string][]byte{}}
	if dir == "" {
		return c
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		c.dir = ""
		return c
	}
	if b, err := os.ReadFile(filepath.Join(dir, revsFile)); err == nil {
		_ = json.Unmarshal(b, &c.revs)
	}
	return c
}

// Rev returns the revision of the cached copy of path, or "".
func (c *Cache) Rev(path string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.revs[path]
}

// Read returns the cached content of path.
func (c *Cache) Read(path string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dir == "" {
		if b, ok := c.files[path]; ok {
			return b, nil
		}
		return nil, os.ErrNotExist
	}
	p, err := c.local(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(p)
}

// Write stores content as path, recording its revision.
func (c *Cache) Write(path string, content []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dir == "" {
		c.files[path] = content
		c.revs[path] = Rev(content)
		return nil
	}
	p, err := c.local(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, content, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, p); err != nil {
		return err
	}
	c.revs[path] = Rev(content)
	return c.saveRevs()
}

// Retain deletes every cached file not in keep.
func (c *Cache) Retain(keep map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for path := range c.revs {
		if _, ok := keep[path]; ok || path == IndexPath {
			continue
		}
		delete(c.revs, path)
		delete(c.files, path)
		if c.dir != "" {
			if p, err := c.local(path); err == nil {
				os.Remove(p)
			}
		}
	}
	if c.dir != "" {
		_ = c.saveRevs()
	}
}

// Clear deletes everything cached.
func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.revs = map[string]string{}
	c.files = map[string][]byte{}
	if c.dir == "" {
		return nil
	}
	if err := os.RemoveAll(c.dir); err != nil {
		return err
	}
	return os.MkdirAll(c.dir, 0o755)
}

func (c *Cache) saveRevs() error {
	b, err := json.Marshal(c.revs)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.dir, revsFile), b, 0o644)
}

// local maps a remote path to a file in the cache directory, refusing
// paths that would escape it.
func (c *Cache) local(path string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(path))
	if filepath.IsAbs(clean) || clean == "." || strings.HasPrefix(clean, "..") {
		return "", errors.New("invalid path " + path)
	}
	return filepath.Join(c.dir, clean), nil
}
