package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Prefs are the user's settings, saved between runs.
type Prefs struct {
	Dark       bool     `json:"dark"`
	DataURL    string   `json:"dataURL,omitempty"`
	Compare    []string `json:"compare,omitempty"`
	Recent     []string `json:"recent,omitempty"`
	LastSynced int64    `json:"lastSynced,omitempty"`
}

func loadPrefs(path string) Prefs {
	p := Prefs{Dark: true}
	if path == "" {
		return p
	}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &p)
	}
	return p
}

func savePrefs(path string, p Prefs) {
	if path == "" {
		return
	}
	b, err := json.Marshal(p)
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, b, 0o644)
}
