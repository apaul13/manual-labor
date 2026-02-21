package api

import (
	"time"

	"github.com/apaul13/manual-labor/api/cars"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunRouter() {
	router := gin.Default()
	// router.Use(cors.Default()) // Allows all origins

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))

	router.GET("/cars", cars.GetCars)
	router.POST("/cars", cars.PostCars)
	router.GET("/health", Health)

	router.Run(":8080")
}

func Health(c *gin.Context) {

}
