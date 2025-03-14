package shell

import "embed"

//go:embed templates/*.tmpl
var templates embed.FS
