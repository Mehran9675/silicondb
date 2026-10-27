// Package assets embeds the files the app ships with: fonts, icons and
// the app icon. Chip data is never embedded; it is downloaded at runtime.
package assets

import "embed"

// Fonts holds the Inter typeface in the weights the UI uses.
//
//go:embed fonts/*.ttf
var Fonts embed.FS

// Icons holds the UI's SVG icons, one file per icon.
//
//go:embed icons/*.svg
var Icons embed.FS

// AppIcon is the window and launcher icon.
//
//go:embed icon.png
var AppIcon []byte
