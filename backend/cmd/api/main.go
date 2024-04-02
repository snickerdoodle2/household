package main

import (
	"log/slog"
	"net/http"
	"os"
)

type App struct {
	logger *slog.Logger
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := App{
		logger: logger,
	}

	app.logger.Info("starting server", "addr", "http://127.0.0.1:8000")
	http.ListenAndServe(":8000", app.routes())
}
