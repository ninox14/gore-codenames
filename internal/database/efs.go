package database

import (
	"embed"
)

//go:embed "migrations"
var EmbeddedFiles embed.FS
