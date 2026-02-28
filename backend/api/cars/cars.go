package cars

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/apaul13/manual-labor/database"
	"github.com/gin-gonic/gin"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/scan"
)

type Make struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Year string `db:"year"`
}

type MakePaginated struct {
	Makes             []Make
	PaginationDetails PaginationDetails
}

type Model struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	MakeID int64  `db:"make_id"`
}

type ModelPaginated struct {
	Models            []Model
	PaginationDetails PaginationDetails
}

type Car struct {
	ID    int64
	Make  string
	Model string
	Year  string
	Trim  string
}

type Year struct {
	Year string
}

// type YearPaginated struct {
// 	Years             []Year
// 	PaginationDetails PaginationDetails
// }

// TotalCount is used to scan a COUNT(*) result from the database.
type TotalCount struct {
	Total int64 `db:"total"`
}

type PaginationDetails struct {
	Offset int
	Limit  int
	Total  int64
}

type VINLookupCar struct {
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  string `json:"year"`
	Trim  string `json:"trim"`
}

// var carTable = psql.NewTable[any, Car, CarSetter]("public", "car")

func GetCars(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	makeStr := c.Query("make")
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

	if cars == nil {
		cars = []Car{}
	}

	c.IndentedJSON(http.StatusOK, cars)
}

// Returns vehicle makes ordered by name ascending with pagination limits
func GetMakes(c *gin.Context) {
	err := errors.New("offset and limit are required parameters.")
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	if !CheckPaginationParams(offsetStr, limitStr) {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	db, err := database.GetDB()

	yearStr := c.Query("year")

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Base query used for retrieving paginated makes
	q := psql.Select(sm.From("make"))
	if len(yearStr) > 0 {
		q.Apply(sm.Where(psql.Quote("year").EQ(psql.Arg(yearStr))))
	}

	// Build a total count query (same filters but without limit/offset)
	totalQuery := psql.Select(sm.Columns("COUNT(*) AS total"), sm.From("make"))
	if len(yearStr) > 0 {
		totalQuery.Apply(sm.Where(psql.Quote("year").EQ(psql.Arg(yearStr))))
	}

	var totalRows []TotalCount
	totalRows, err = bob.All(c, db, totalQuery, scan.StructMapper[TotalCount]())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var total int64
	if len(totalRows) > 0 {
		total = totalRows[0].Total
	}

	// Apply pagination and ordering to the list query
	q.Apply(sm.Offset(offsetStr))
	q.Apply(sm.Limit(limitStr))
	q.Apply(sm.OrderBy("name").Asc())

	makes, err := bob.All(c, db, q, scan.StructMapper[Make]())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if makes == nil {
		makes = []Make{}
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	limitInt, err := strconv.Atoi(limitStr)

	paginationDetails := PaginationDetails{
		Offset: offsetInt,
		Limit:  limitInt,
		Total:  total,
	}

	makesPaginated := MakePaginated{
		Makes:             makes,
		PaginationDetails: paginationDetails,
	}

	c.IndentedJSON(http.StatusOK, makesPaginated)
}

// Returns years in descending order with pagination limits
func GetYears(c *gin.Context) {
	// err := errors.New("offset and limit are required parameters.")
	// offsetStr := c.Query("offset")
	// limitStr := c.Query("limit")
	// if !CheckPaginationParams(offsetStr, limitStr) {
	// 	c.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }

	db, err := database.GetDB()
	q := psql.Select(sm.Columns("year"),
		sm.Distinct(),
		sm.From("make"),
		sm.OrderBy("year").Desc())

	years, err := bob.All(c, db, q, scan.StructMapper[Year]())

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, years)
}

func GetModels(c *gin.Context) {
	err := errors.New("offset and limit are required parameters.")
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	if !CheckPaginationParams(offsetStr, limitStr) {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	makeIdStr := c.Query("make")

	fmt.Print(makeIdStr)
	db, err := database.GetDB()
	q := psql.Select(
		sm.From("model"))

	if len(makeIdStr) > 0 {
		q.Apply(sm.Where(psql.Quote("make_id").EQ(psql.Arg(makeIdStr))))
	}

	// Apply ordering and pagination before executing the query
	q.Apply(sm.OrderBy("name").Desc())
	q.Apply(sm.Offset(offsetStr))
	q.Apply(sm.Limit(limitStr))

	var buf strings.Builder
	_, error := q.WriteSQL(c, &buf, q.Dialect, 1)
	if error != nil {
		// log.Fatal(err)
	}
	fmt.Println("SQL:", buf.String())

	models, err := bob.All(c, db, q, scan.StructMapper[Model]())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Normalize nil slice to an empty slice so JSON encodes [] instead of null
	if models == nil {
		models = []Model{}
	}

	totalQuery := psql.Select(sm.Columns("COUNT(*) AS total"), sm.From("model"))
	if len(makeIdStr) > 0 {
		totalQuery.Apply(sm.Where(psql.Quote("make_id").EQ(psql.Arg(makeIdStr))))
	}

	var totalRows []TotalCount
	totalRows, err = bob.All(c, db, totalQuery, scan.StructMapper[TotalCount]())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var total int64
	if len(totalRows) > 0 {
		total = totalRows[0].Total
	}

	// if err != nil {
	// 	c.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }

	offsetInt, err := strconv.Atoi(offsetStr)
	limitInt, err := strconv.Atoi(limitStr)

	paginationDetails := PaginationDetails{
		Offset: offsetInt,
		Limit:  limitInt,
		Total:  total,
	}

	modelsPaginated := ModelPaginated{
		Models:            models,
		PaginationDetails: paginationDetails,
	}

	c.IndentedJSON(http.StatusOK, modelsPaginated)
}

// Not needed as of now, will need a refactor with DB changes
func PostCars(c *gin.Context) {

	db, err := database.GetDB()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var newCar Car
	if err := c.BindJSON(&newCar); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	q := psql.Insert(im.Into("car", "make", "model", "year", "trim"),
		im.Values(psql.Arg(strings.ToUpper(newCar.Make)), psql.Arg(strings.ToUpper(newCar.Model)), psql.Arg(strings.ToUpper(newCar.Year)), psql.Arg(strings.ToUpper(newCar.Trim))))

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

func LookupVIN(c *gin.Context) {
	endpoint := "https://vpic.nhtsa.dot.gov/api/vehicles/decodevin/" + c.Query("vin") + "?format=json"
	if len(c.Query("year")) > 0 {
		endpoint = endpoint + "&modelyear=" + c.Query("year")
	}

	fmt.Println("Calling API with endpoint: " + endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// helper types to unmarshal NHTSA response
	type resultItem struct {
		Variable   string `json:"Variable"`
		VariableId int    `json:"VariableId"`
		Value      any    `json:"Value"`
	}
	var apiResp struct {
		Count   int          `json:"Count"`
		Message string       `json:"Message"`
		Results []resultItem `json:"Results"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	vinLookupCar := &VINLookupCar{}

	// Convert interface{} Value -> string helper
	toStr := func(v any) string {
		if v == nil {
			return ""
		}
		switch t := v.(type) {
		case string:
			return strings.TrimSpace(t)
		case float64:
			// JSON numbers unmarshal to float64
			// Format without fractional part if it's an integer
			if t == float64(int64(t)) {
				return fmt.Sprintf("%d", int64(t))
			}
			return fmt.Sprintf("%v", t)
		default:
			// fallback: marshal back to JSON and use that
			b, _ := json.Marshal(t)
			return strings.TrimSpace(string(b))
		}
	}

	// Map variables
	for _, r := range apiResp.Results {
		val := toStr(r.Value)
		switch r.Variable {
		case "Make":
			if vinLookupCar.Make == "" {
				vinLookupCar.Make = val
			}
		case "Model":
			if vinLookupCar.Model == "" {
				vinLookupCar.Model = val
			}
		case "Trim":
			if vinLookupCar.Trim == "" {
				vinLookupCar.Trim = val
			}
		case "Model Year", "ModelYear":
			if vinLookupCar.Year == "" {
				vinLookupCar.Year = val
			}
		// optional: handle other common names / fallback by VariableId
		case "":
			// no-op
		}
		// fallback by VariableId if needed:
		if vinLookupCar.Make == "" && r.VariableId == 26 {
			vinLookupCar.Make = val
		}
		if vinLookupCar.Model == "" && r.VariableId == 28 {
			vinLookupCar.Model = val
		}
		if vinLookupCar.Year == "" && r.VariableId == 29 {
			vinLookupCar.Year = val
		}
		if vinLookupCar.Trim == "" && r.VariableId == 38 {
			vinLookupCar.Trim = val
		}
	}

	strings.ToUpper(vinLookupCar.Make)
	strings.ToUpper(vinLookupCar.Model)
	strings.ToUpper(vinLookupCar.Trim)

	c.IndentedJSON(http.StatusOK, vinLookupCar)
}

func CheckPaginationParams(offset string, limit string) bool {
	if len(offset) == 0 || len(limit) == 0 {
		return false
	}

	return true
}
