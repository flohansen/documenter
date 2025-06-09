package web

import "embed"

var (
	//go:embed templates/*.gohtml
	Templates embed.FS
)
