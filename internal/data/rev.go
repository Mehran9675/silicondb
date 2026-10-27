package data

import (
	"crypto/sha256"
	"encoding/hex"
)

// Rev returns the revision of a file's content: the first 16 hex digits of
// its SHA-256. The index generator computes revisions the same way.
func Rev(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:8])
}
