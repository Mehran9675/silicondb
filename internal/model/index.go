package model

// FileRef points to a single file from the index.
type FileRef struct {
	Path string `json:"path"`
	Rev  string `json:"rev"`
}

// IndexEntry lists a manufacturer, copied from its manufacturer.json, and
// its shards.
type IndexEntry struct {
	Manufacturer
	Shards []ShardRef `json:"shards"`
}

// Index is the table of contents of the remote database. Clients fetch it
// first and then download only the files whose revision changed.
type Index struct {
	Format        int          `json:"format"`
	Updated       string       `json:"updated"`
	Schema        FileRef      `json:"schema"`
	Manufacturers []IndexEntry `json:"manufacturers"`
}

// Files returns every file the index refers to, keyed by path, with its
// revision.
func (ix *Index) Files() map[string]string {
	files := map[string]string{ix.Schema.Path: ix.Schema.Rev}
	for _, m := range ix.Manufacturers {
		for _, s := range m.Shards {
			files[s.Path] = s.Rev
		}
	}
	return files
}
