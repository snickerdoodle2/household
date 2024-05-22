package main

import (
	"flag"
	"fmt"
	"inzynierka/internal/data"
	"inzynierka/internal/listener"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
	logger    *slog.Logger
	models    data.Models
	upgrader  websocket.Upgrader
	listeners map[uuid.UUID]listener.ListenerT
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
		listeners: make(map[uuid.UUID]listener.ListenerT),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == fmt.Sprintf("http://localhost:%d", cfg.port) {
					return true
				}
				for _, v := range cfg.cors.trustedOrigins {
					if origin == v {
						return true
					}
				}
				return false
			},
		},
	}

	err = app.serve()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
