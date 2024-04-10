package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func openDB(cfg Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(cfg.db.maxOpenConns)
	config.MaxConnIdleTime = cfg.db.maxIdleTime

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
