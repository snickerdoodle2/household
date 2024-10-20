package static

import "embed"

// go:embed "all:generated"
var StaticFiles embed.FS
