package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest("GET", "/api/v1/stations/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")
	router := gin.Default()

	s := Station{StationId: "1", Name: "Test", ShortName: "Short Test Name"}
	stations := &StationsData{Data: Stations{Stations: []Station{s}}}

	router.GET("/api/v1/stations/:id", func(c *gin.Context) {
		GetStationById(c, stations)
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var station Station
	err := json.Unmarshal([]byte(w.Body.String()), &station)
	assert.Nil(t, err)
	assert.Equal(t, station.Name, "Test")
	assert.Equal(t, station.ShortName, "Short Test Name")
}
