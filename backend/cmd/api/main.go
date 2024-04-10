package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type App struct {
	config Config
	logger *slog.Logger
}

func main() {
	var cfg Config
	flag.IntVar(&cfg.port, "port", 8080, "API Server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "Database DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := App{
		logger: logger,
	}

	app.logger.Info("starting server", "addr", "http://127.0.0.1:8000")
	http.ListenAndServe(":8000", app.routes())
}
