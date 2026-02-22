package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stephenafamo/bob"
)

// This package provides a single shared DB pool for the application.
// Use InitDB at application startup, GetDB in handlers, and CloseDB on shutdown.

var (
	mu         sync.RWMutex
	globalDB   bob.DB
	globalPool *pgxpool.Pool
)

// InitDB initializes the global database pool and bob.DB wrapper.
// It is safe to call multiple times; subsequent calls are no-ops.
func InitDB() error {
	// Try to load an env file if present. In Docker we supply env vars via docker-compose
	// so failure to load is not fatal and will be silently ignored.
	_ = godotenv.Load("../.env")

	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return errors.New("DB_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	// Tune pool as desired
	pool.Config().MaxConns = 100

	db := bob.NewDB(stdlib.OpenDBFromPool(pool))

	mu.Lock()
	globalPool = pool
	globalDB = db
	mu.Unlock()

	return nil
}

// GetDB returns the initialized global DB instance. Do NOT Close() the returned value
// in handlers — CloseDB should be called once at application shutdown.
func GetDB() (bob.DB, error) {
	mu.RLock()
	db := globalDB
	pool := globalPool
	mu.RUnlock()

	// Use the underlying pool pointer to determine initialization state.
	if pool == nil {
		// Return the zero value of bob.DB and an error when uninitialized.
		return bob.DB{}, errors.New("database not initialized")
	}
	return db, nil
}

// CloseDB closes the global DB and underlying pool. Call once at shutdown.
func CloseDB() error {
	mu.Lock()
	defer mu.Unlock()

	var err error
	// Use globalPool as the initialization indicator; when it's non-nil we should perform cleanup.
	if globalPool != nil {
		// Close the bob DB wrapper and capture any error.
		if e := globalDB.Close(); e != nil {
			err = e
		}
		// Reset globalDB to its zero value (cannot assign nil to a concrete bob.DB type).
		globalDB = bob.DB{}

		// Close and nil out the underlying pgx pool.
		globalPool.Close()
		globalPool = nil
	}
	return err
}

// Backwards-compatible helper for existing code that used GetDbConnection.
// Prefer InitDB/GetDB/CloseDB for lifecycle control.
func GetDbConnection() (bob.DB, error) {
	return GetDB()
}
