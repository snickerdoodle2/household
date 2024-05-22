package static

import "embed"

//go:embed "_app" "favicon.png" "index.html"
var StaticFiles embed.FS
