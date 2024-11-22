package main

import (
	"flag"
	"inzynierka/internal/data"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type Config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	cors struct {
		trustedOrigins []string
	}
}

type App struct {
	config    Config
	logger    *log.Logger
	models    data.Models
	listeners data.SensorListeners
	rules     struct {
		channel      chan data.ValidRuleAction
		stopChannels map[uuid.UUID]chan struct{}
	}
}

func main() {
	var cfg Config
	flag.IntVar(&cfg.port, "port", 8080, "API Server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "Database DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated) (eg. http://localhost:5173)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
	})

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer db.Close()

	logger.Info("DB connection established")

	app := App{
		logger:    logger,
		config:    cfg,
		models:    data.NewModels(db),
		listeners: make(data.SensorListeners),
		rules: struct {
			channel      chan data.ValidRuleAction
			stopChannels map[uuid.UUID]chan struct{}
		}{
			channel:      make(chan data.ValidRuleAction, 1),
			stopChannels: make(map[uuid.UUID]chan struct{}),
		},
	}

	err = app.serve()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
