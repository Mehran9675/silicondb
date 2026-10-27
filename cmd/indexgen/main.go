// Command indexgen validates the database under data/ and rewrites
// data/index.json with a revision for every file, so apps download only
// what changed.
//
//	go run ./cmd/indexgen            # from the repository root
//	go run ./cmd/indexgen -check     # validate only, don't write
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mehran9675/silicondb/internal/data"
	"github.com/mehran9675/silicondb/internal/model"
)

func main() {
	dir := flag.String("dir", "data", "database directory")
	check := flag.Bool("check", false, "validate without writing index.json")
	gaps := flag.Bool("report", false, "list each chip's missing specs instead of writing index.json")
	prefix := flag.String("only", "", "with -report, only chips whose id starts with this")
	flag.Parse()

	ix, report, err := build(*dir, *gaps, *prefix)
	if err != nil {
		fmt.Fprintln(os.Stderr, "indexgen:", err)
		os.Exit(1)
	}
	for _, w := range report.warnings {
		fmt.Fprintln(os.Stderr, "warning:", w)
	}
	if len(report.errors) > 0 {
		for _, e := range report.errors {
			fmt.Fprintln(os.Stderr, "error:", e)
		}
		os.Exit(1)
	}
	fmt.Printf("%d manufacturers, %d shards, %d chips\n", len(ix.Manufacturers), report.shards, report.chips)
	if report.complete != nil {
		report.complete.print()
		return
	}
	if *check {
		return
	}
	out := filepath.Join(*dir, data.IndexPath)
	if old, err := os.ReadFile(out); err == nil {
		var prev model.Index
		if json.Unmarshal(old, &prev) == nil && sameFiles(&prev, ix) {
			fmt.Println("index.json is up to date")
			return
		}
	}
	b, err := json.MarshalIndent(ix, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "indexgen:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(out, append(b, '\n'), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "indexgen:", err)
		os.Exit(1)
	}
	fmt.Println("wrote", out)
}

func build(dir string, gaps bool, prefix string) (*model.Index, *report, error) {
	rep := &report{}
	schemaRaw, err := os.ReadFile(filepath.Join(dir, "schema.json"))
	if err != nil {
		return nil, nil, err
	}
	var schema model.Schema
	if err := strictDecode(schemaRaw, &schema); err != nil {
		return nil, nil, fmt.Errorf("schema.json: %w", err)
	}
	if gaps {
		rep.complete = newCompleteness(&schema, prefix)
		rep.gaps = rep.complete.add
	}
	ix := &model.Index{
		Format:  data.SupportedFormat,
		Updated: time.Now().UTC().Format("2006-01-02"),
		Schema:  model.FileRef{Path: "schema.json", Rev: data.Rev(schemaRaw)},
	}

	root := filepath.Join(dir, "manufacturers")
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, nil, err
	}
	ids := map[string]string{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		entry, ok := readMaker(root, e.Name(), &schema, ids, rep)
		if ok {
			ix.Manufacturers = append(ix.Manufacturers, *entry)
		}
	}
	sort.Slice(ix.Manufacturers, func(i, j int) bool {
		return strings.ToLower(ix.Manufacturers[i].Name) < strings.ToLower(ix.Manufacturers[j].Name)
	})
	return ix, rep, nil
}

func readMaker(root, id string, schema *model.Schema, ids map[string]string, rep *report) (*model.IndexEntry, bool) {
	dir := filepath.Join(root, id)
	raw, err := os.ReadFile(filepath.Join(dir, "manufacturer.json"))
	if err != nil {
		rep.errorf("%s: %v", id, err)
		return nil, false
	}
	entry := &model.IndexEntry{}
	if err := strictDecode(raw, &entry.Manufacturer); err != nil {
		rep.errorf("%s/manufacturer.json: %v", id, err)
		return nil, false
	}
	if entry.ID != id {
		rep.errorf("%s/manufacturer.json: id %q doesn't match its folder", id, entry.ID)
	}
	if entry.Name == "" {
		rep.errorf("%s/manufacturer.json: name is empty", id)
	}
	files, _ := os.ReadDir(dir)
	for _, f := range files {
		name := f.Name()
		if f.IsDir() || name == "manufacturer.json" || !strings.HasSuffix(name, ".json") {
			continue
		}
		cat := strings.TrimSuffix(name, ".json")
		rel := path.Join("manufacturers", id, name)
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			rep.errorf("%s: %v", rel, err)
			continue
		}
		n := validateShard(rel, b, id, cat, schema, ids, rep)
		if n == 0 {
			continue
		}
		entry.Shards = append(entry.Shards, model.ShardRef{Path: rel, Category: cat, Count: n, Rev: data.Rev(b)})
		rep.shards++
		rep.chips += n
	}
	if len(entry.Shards) == 0 {
		rep.warnf("%s: no chips", id)
		return nil, false
	}
	return entry, true
}

// sameFiles reports whether two indexes list the same files and
// manufacturers, so the date isn't bumped when nothing changed.
func sameFiles(a, b *model.Index) bool {
	ca, _ := json.Marshal(a.Manufacturers)
	cb, _ := json.Marshal(b.Manufacturers)
	return a.Schema == b.Schema && a.Format == b.Format && bytes.Equal(ca, cb)
}

func strictDecode(b []byte, v any) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	return d.Decode(v)
}
