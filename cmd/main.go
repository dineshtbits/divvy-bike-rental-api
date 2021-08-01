package main

import (
	"fmt"

	"github.com/dineshtbits/divvy-bike-rental-api/pkg/rental"
	"github.com/dineshtbits/divvy-bike-rental-api/pkg/station"
	"github.com/gin-gonic/gin"
)

const basicAuthUsername = "admin"
const basicAuthPassword = "admin"

func main() {
	fmt.Println("Running Divvys Bike Rental API")
	fmt.Println("Initial load of data")
	stations := station.LoadStationsData()
	trips, _ := rental.LoadTripsData()
	fmt.Println("Initial load of data - Complete")

	router := gin.Default()

	// Setting up basic auth for all api routes
	authenticatedRoutes := router.Group("/api/v1/", gin.BasicAuth(gin.Accounts{
		basicAuthUsername: basicAuthPassword,
	}))

	authenticatedRoutes.GET("/stations/:id", func(c *gin.Context) {
		station.GetStationById(c, stations)
	})
	authenticatedRoutes.POST("/trips/riders/summary", func(c *gin.Context) {
		rental.GetRidersSummary(c, trips)
	})
	authenticatedRoutes.POST("/trips/summary", func(c *gin.Context) {
		rental.GetTripsSummary(c, trips)
	})
	router.Run("localhost:8081")
}
