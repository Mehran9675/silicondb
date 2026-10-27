package store

import "time"

// Phase is what the sync is doing.
type Phase int

const (
	// Idle: no sync has run yet.
	Idle Phase = iota
	// Syncing: downloading the index or changed files.
	Syncing
	// Synced: the offline copy matches the remote database.
	Synced
	// Failed: the last sync failed; the offline copy, if any, is shown.
	Failed
)

// SyncStatus describes the latest sync.
type SyncStatus struct {
	Phase      Phase
	Done       int
	Total      int
	Files      int
	Bytes      int
	Err        string
	LastSynced time.Time
}
