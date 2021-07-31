package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest("GET", "/api/v1/PersonId/Id456", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")
	router := gin.Default()
	// router.GET("/api/v1/stations/:id", GetStationById)

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)
	fmt.Println(w.Body)
}
