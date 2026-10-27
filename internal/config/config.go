// Package config holds the app's build-time settings.
package config

// DataURL is where the database lives: a folder holding index.json,
// schema.json and the manufacturer shards. Override it at build time with
//
//	-ldflags "-X github.com/mehran9675/silicondb/internal/config.DataURL=https://example.com/data/"
//
// or at runtime in Settings.
var DataURL = "https://raw.githubusercontent.com/Mehran9675/silicondb/main/data/"

// Version is the app version shown in Settings.
var Version = "1.0.0"

// MaxCompare is how many chips the comparison holds at once.
const MaxCompare = 4
