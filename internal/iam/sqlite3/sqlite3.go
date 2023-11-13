package sqlite3

import (
	"embed"

	"github.com/adoublef/yap/sqlite3"
)

//go:embed all:*.up.sql
var embedFS embed.FS
var Migrate = sqlite3.NewFS(embedFS)
