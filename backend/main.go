package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Car struct {
	ID    int64
	Make  string
	Model string
	Year  string
	Trim  string
}

var db *pgx.Conn

func main() {
	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load env file: %v\n", err)
	}

	db, err = pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close(context.Background())

	// Example query run before starting server; ListenAndServe below blocks.
	cars, err := carsByMake("Honda")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Cars found: %v\n", cars)

	router := gin.Default()
	router.Use(cors.Default()) // Allows all origins

	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost.com"}
	// config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	// config.AllowCredentials = true
	// config.MaxAge = 12 * time.Hour

	// router.Use(cors.New(config))

	router.GET("/cars", getCars)
	router.POST("/cars", postCars)

	router.Run("localhost:8080")
}

func getCars(c *gin.Context) {
	var cars []Car

	rows, err := db.Query(context.Background(), "SELECT * FROM car")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Make, &car.Model, &car.Year, &car.Trim); err != nil {
			fmt.Fprintf(os.Stderr, "Scan for car failed: %v\n", err)
		}

		cars = append(cars, car)
	}
	c.IndentedJSON(http.StatusOK, cars)
}

func postCars(c *gin.Context) {
	var newCar Car
	if err := c.BindJSON(&newCar); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec(context.Background(), "INSERT INTO car (make, model, year, trim) VALUES ($1, $2, $3, $4)", newCar.Make, newCar.Model, newCar.Year, newCar.Trim)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result.Insert() {
		c.IndentedJSON(http.StatusCreated, newCar)
	}
}

func carsByMake(make string) ([]Car, error) {
	var cars []Car

	rows, err := db.Query(context.Background(), "SELECT * FROM car WHERE make = $1", make)

	if err != nil {
		return nil, fmt.Errorf("carsByMake %q: %v", make, err)
	}

	defer rows.Close()
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.Make, &car.Model, &car.Year, &car.Trim); err != nil {
			return nil, fmt.Errorf("carsByMake %q: %v", make, err)
		}

		cars = append(cars, car)
	}
	return cars, nil

}
