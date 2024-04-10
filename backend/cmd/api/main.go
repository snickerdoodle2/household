package main

import (
	"flag"
	"log/slog"
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
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "Database DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer db.Close()

	logger.Info("DB connection established")

	app := App{
		logger: logger,
		config: cfg,
	}

	err = app.serve()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
