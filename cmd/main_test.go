package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dineshtbits/divvy-bike-rental-api/pkg/rental"
	"github.com/dineshtbits/divvy-bike-rental-api/pkg/station"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetStation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest("GET", "/api/v1/stations/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")
	router := gin.Default()

	s := station.Station{StationId: "1", Name: "Test", ShortName: "Short Test Name"}
	stations := &station.StationsData{Data: station.Stations{Stations: []station.Station{s}}}

	router.GET("/api/v1/stations/:id", func(c *gin.Context) {
		station.GetStationById(c, stations)
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var station station.Station
	err := json.Unmarshal([]byte(w.Body.String()), &station)
	assert.Nil(t, err)
	assert.Equal(t, station.Name, "Test")
	assert.Equal(t, station.ShortName, "Short Test Name")
}

func TestGetRidersSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := bytes.NewBuffer([]byte("{\"filters\":{\"station_ids\":[176]}}"))
	req, _ := http.NewRequest("POST", "/api/v1/trips/riders/summary", body)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")
	router := gin.Default()

	trips := []rental.Rental{}
	trips = append(trips, rental.Rental{ID: 22178529, StartTime: time.Now(), EndTime: time.Now().Add(time.Hour), BikeID: 1, Duration: 60.8, StartStationID: 123, EndStationID: 176, MemberBirthdayYear: 1987})
	trips = append(trips, rental.Rental{ID: 22178530, StartTime: time.Now(), EndTime: time.Now().Add(time.Hour), BikeID: 1, Duration: 60.8, StartStationID: 123, EndStationID: 176, MemberBirthdayYear: 1998})
	trips = append(trips, rental.Rental{ID: 22178531, StartTime: time.Now(), EndTime: time.Now().Add(time.Hour), BikeID: 1, Duration: 60.8, StartStationID: 123, EndStationID: 176, MemberBirthdayYear: 2010})
	trips = append(trips, rental.Rental{ID: 22178532, StartTime: time.Now(), EndTime: time.Now().Add(time.Hour), BikeID: 1, Duration: 60.8, StartStationID: 123, EndStationID: 176, MemberBirthdayYear: 2021})

	router.POST("/api/v1/trips/riders/summary", func(c *gin.Context) {
		rental.GetRidersSummary(c, &trips)
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body.String())
	assert.Contains(t, w.Body.String(), "2021-07-31")
}
