package cars

import (
	"net/http"

	"github.com/apaul13/manual-labor/database"
	"github.com/gin-gonic/gin"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/scan"
)

type Car struct {
	ID    int64 // `db:",pk"`
	Make  string
	Model string
	Year  string
	Trim  string
}

// var carTable = psql.NewTable[any, Car, CarSetter]("public", "car")

func GetCars(c *gin.Context) {
	db, err := database.GetDbConnection()

	makeStr := c.Query("Make")
	modelStr := c.Query("model")
	yearStr := c.Query("year")
	trimStr := c.Query("trim")

	// Incredibly beautiful dynamic query
	q := psql.Select(sm.From("car"))
	if len(makeStr) > 0 {
		q.Apply(sm.Where(psql.Quote("make").In(psql.Arg(makeStr))))
	}
	if len(modelStr) > 0 {
		q.Apply(sm.Where(psql.Quote("model").In(psql.Arg(modelStr))))
	}
	if len(yearStr) > 0 {
		q.Apply(sm.Where(psql.Quote("year").In(psql.Arg(yearStr))))
	}
	if len(trimStr) > 0 {
		q.Apply(sm.Where(psql.Quote("trim").In(psql.Arg(trimStr))))
	}

	cars, err := bob.All(c, db, q, scan.StructMapper[Car]())

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	q := psql.Insert(im.Into("car", "make", "model", "year", "trim"),
		im.Values(psql.Arg(newCar.Make), psql.Arg(newCar.Model), psql.Arg(newCar.Year), psql.Arg(newCar.Trim)))

	result, err := bob.Exec(c, db, q)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsInserted, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsInserted > 0 {
		c.IndentedJSON(http.StatusCreated, newCar)
	} else {
		c.IndentedJSON(http.StatusNoContent, "No new cars created.")
	}

}
