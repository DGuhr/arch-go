package html

import "embed"

//go:embed templates/*.tmpl
var templateFiles embed.FS

//go:embed assets/*.css
var styles embed.FS

//go:embed assets/*.png
var images embed.FS
