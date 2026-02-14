package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

type Car struct {
	id    int64
	make  string
	model string
	year  string
	trim  string
}

var db *sql.DB

func main() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(context.Background())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":8080", nil)
}

func carsByName(name string) ([]Car, error) {
	var cars []Car

	rows, err := db.Query("SELECT * FROM car WHERE name = ?", name)

	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return cars, nil

}
