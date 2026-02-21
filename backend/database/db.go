package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stephenafamo/bob"
)

func GetDbConnection() (bob.DB, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load env file: %v\n", err)
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	// defer pool.Close()

	pool.Config().MaxConns = 100

	// Wrap pool for Bob
	db := bob.NewDB(stdlib.OpenDBFromPool(pool))

	// defer func() {
	// 	_ = db.Close()
	// }()

	return db, err
}
