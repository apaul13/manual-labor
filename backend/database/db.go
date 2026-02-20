package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func GetDbConnection() (*pgxpool.Pool, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load env file: %v\n", err)
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	pool.Config().MaxConns = 100

	return pool, err
}
