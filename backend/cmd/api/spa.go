package main

import (
	"inzynierka/static"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

func (app *App) spaHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	fi, err := fs.Stat(static.StaticFiles, path)

	app.logger.Info("SPA", "path", path, "exist?", !os.IsNotExist(err))
	if os.IsNotExist(err) || fi.IsDir() {
		// no caching
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		// just serve initial page lool
		http.ServeFileFS(w, r, static.StaticFiles, "index.html")
		return
	}
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	http.ServeFileFS(w, r, static.StaticFiles, path)
}
