package cars

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/apaul13/manual-labor/database"
	"github.com/gin-gonic/gin"
)

type Car struct {
	ID    int64
	Make  string
	Model string
	Year  string
	Trim  string
}

func GetCars(c *gin.Context) {
	db, err := database.GetDbConnection()

	var cars []Car

	// var pathParams = c.Params
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

func PostCars(c *gin.Context) {

	db, err := database.GetDbConnection()

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
	db, err := database.GetDbConnection()

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
