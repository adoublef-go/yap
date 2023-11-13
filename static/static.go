package static

import (
	"embed"
	"html/template"
	"path/filepath"

	"github.com/benbjohnson/hashfs"
)

//go:embed all:*.css
var embedFS embed.FS
var hashFS = hashfs.NewFS(embedFS)

var FuncMap = template.FuncMap{
	"static": func(filename string) string {
		return filepath.Join("static", hashFS.HashName(filename))
	},
}

var Handler = hashfs.FileServer(hashFS)
