package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type DB struct {
	*sql.DB
}

func NewConnection(dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// TODO подумать куда вынести, в конфиг/env
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DB{DB: db}, nil
}

func (d *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := d.PingContext(ctx); err != nil {
		zap.L().Error("database health check failed", zap.Error(err))
		return fmt.Errorf("database health check failed: %w", err)
	}

	zap.L().Info("database health check passed")
	return nil
}
