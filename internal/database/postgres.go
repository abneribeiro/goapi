package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/abneribeiro/goapi/internal/config"
	"github.com/abneribeiro/goapi/internal/pkg/logger"
)

func NewPostgresConnection(cfg *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established", logger.WithFields(map[string]interface{}{
		"host": cfg.Host,
		"port": cfg.Port,
		"name": cfg.Name,
	}))

	return db, nil
}
