package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const stationInfoUrl = "https://gbfs.divvybikes.com/gbfs/en/station_information.json"
const tripsFilePath = "../resources/Divvy_Trips_2019_Q2"
const basicAuthUsername = "admin"
const basicAuthPassword = "admin"

type Rental struct {
	ID                 int
	StartTime          time.Time
	EndTime            time.Time
	BikeID             int
	Duration           float64
	StartStationID     int
	StartStationName   string
	EndStationID       int
	EndStationName     string
	UserType           string
	Gender             string
	MemberBirthdayYear int
}

type StationsData struct {
	Data        Stations `json:"data"`
	LastUpdated int      `json:"last_updated"`
	Ttl         int      `json:"ttl"`
}

type Stations struct {
	Stations []Station `json:"stations"`
}

type Uris struct {
	Android string `json:"android"`
	Ios     string `json:"ios"`
}

type RiderSummaryRequest struct {
	Filters RiderSummaryRequestFilters `json:"filters"`
}

type RiderSummaryRequestFilters struct {
	StationIds []int `json:"station_ids"`
}

type RiderSummaryResponse struct {
	Data []StationSummary
}

type StationSummary struct {
	StationId string
	Report    map[string]map[string]int
}

type Station struct {
	Lon                         float64  `json:"lon"`
	LegacyId                    string   `json:"legacy_id"`
	ExternalId                  string   `json:"external_id"`
	Eightd_station_services     []string `json:"eightd_station_services"`
	StationType                 string   `json:"station_type"`
	RentalUris                  []Uris   `json:"rental_uris"`
	ElectricBikeSurchargeWaiver bool     `json:"electric_bike_surcharge_waiver"`
	StationId                   string   `json:"station_id"`
	Capacity                    int      `json:"capacity"`
	HasKiosk                    bool     `json:"has_kiosk"`
	Eightd_has_key_dispenser    bool     `json:"eightd_has_key_dispenser"`
	RentalMethods               []string `json:"rental_methods"`
	Name                        string   `json:"name"`
	ShortName                   string   `json:"short_name"`
	Lat                         float64  `json:"lat"`
}

func main() {
	fmt.Println("Running Divvys Bike Rental API")
	fmt.Println("Initial load of data")
	stations := loadStationsData()
	trips, _ := loadTripsData()
	fmt.Println("Initial load of data - Complete")

	router := gin.Default()

	// Setting up basic auth for all api routes
	authenticatedRoutes := router.Group("/api/v1/", gin.BasicAuth(gin.Accounts{
		basicAuthUsername: basicAuthPassword,
	}))

	authenticatedRoutes.GET("/stations/:id", func(c *gin.Context) {
		GetStationById(c, stations)
	})
	authenticatedRoutes.POST("/trips/riders/summary", func(c *gin.Context) {
		getRidersSummary(c, trips)
	})
	authenticatedRoutes.POST("/trips/summary", func(c *gin.Context) {
		getTripsSummary(c, trips)
	})
	router.Run("localhost:8081")
}

func getTripsSummary(c *gin.Context, trips *[]Rental) {
	m := make(map[string]map[string][]Rental)
	var riderSummaryRequest RiderSummaryRequest
	if err := c.BindJSON(&riderSummaryRequest); err != nil {
		return
	}
	filterd := []Rental{}
	for _, trip := range *trips {
		if contains(riderSummaryRequest.Filters.StationIds, trip.EndStationID) {
			m[strconv.Itoa(trip.EndStationID)] = make(map[string][]Rental)
			filterd = append(filterd, trip)
		}
	}
	sort.Slice(filterd, func(i, j int) bool { return filterd[i].EndTime.After(filterd[j].EndTime) })
	for _, r := range filterd {
		val, _ := m[strconv.Itoa(r.EndStationID)]
		date := r.EndTime.Format("2006-01-02")
		_, ok := val[date]
		if ok {
			if len(val[date]) < 20 {
				val[date] = append(val[date], r)
			}
		} else {
			val[date] = []Rental{r}
		}
	}
	c.IndentedJSON(http.StatusOK, m)
}

func getRidersSummary(c *gin.Context, trips *[]Rental) {

	var riderSummaryRequest RiderSummaryRequest
	if err := c.BindJSON(&riderSummaryRequest); err != nil {
		return
	}

	m := make(map[string]map[string]map[string]int)
	for _, trip := range *trips {
		if contains(riderSummaryRequest.Filters.StationIds, trip.EndStationID) {
			date := trip.EndTime.Format("2006-01-02")
			ageGroup := getAgeGroup(trip.MemberBirthdayYear)
			val, ok := m[strconv.Itoa(trip.EndStationID)]
			if ok {
				val2, ok := val[date]
				if ok {
					val2[ageGroup] += 1
				} else {
					val[date] = map[string]int{ageGroup: 1}
				}
			} else {
				age := map[string]int{ageGroup: 1}
				d := map[string]map[string]int{date: age}
				m[strconv.Itoa(trip.EndStationID)] = d
			}
		}
	}
	c.IndentedJSON(http.StatusOK, m)
}

func GetStationById(c *gin.Context, stations *StationsData) {
	for i := 0; i < len(stations.Data.Stations); i++ {
		if stations.Data.Stations[i].StationId == c.Param("id") {
			c.IndentedJSON(http.StatusOK, stations.Data.Stations[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "station not found"})
}

func loadStationsData() *StationsData {
	req, err := http.NewRequest(http.MethodGet, stationInfoUrl, nil)
	if err != nil {
		fmt.Println("Unable to create new request to get json")
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Unable to get json")
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	byteValue, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}
	var stationsData StationsData
	json.Unmarshal(byteValue, &stationsData)
	return &stationsData
}

func loadTripsData() (*[]Rental, error) {
	rentals := []Rental{}
	f, err := os.Open(tripsFilePath)
	if err != nil {
		return &rentals, err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return &rentals, err
	}

	for i := 1; i < len(lines); i++ {
		startTime, _ := time.Parse("2006-01-02 15:04:05", lines[i][1])
		endTime, _ := time.Parse("2006-01-02 15:04:05", lines[i][2])
		id, _ := strconv.Atoi(lines[i][0])
		bikeId, _ := strconv.Atoi(lines[i][3])
		duration, _ := strconv.ParseFloat(lines[i][4], 5)
		startStationID, _ := strconv.Atoi(lines[i][5])
		endStationID, _ := strconv.Atoi(lines[i][7])
		memberBirthdayYear, _ := strconv.Atoi(lines[i][11])
		r := Rental{
			ID:                 id,
			StartTime:          startTime,
			EndTime:            endTime,
			BikeID:             bikeId,
			Duration:           duration,
			StartStationID:     startStationID,
			StartStationName:   lines[i][6],
			EndStationID:       endStationID,
			EndStationName:     lines[i][8],
			UserType:           lines[i][9],
			Gender:             lines[i][10],
			MemberBirthdayYear: memberBirthdayYear,
		}
		rentals = append(rentals, r)
	}
	return &rentals, nil
}

func contains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func getAgeGroup(birthYear int) string {
	age := time.Now().Year() - birthYear
	// [0-20,21-30,31-40,41-50,51+, unknown]
	switch {
	case age >= 0 && age <= 20:
		return "0-20"
	case age >= 21 && age <= 30:
		return "21-30"
	case age >= 31 && age <= 40:
		return "31-40"
	case age >= 41 && age <= 50:
		return "41-50"
	case age >= 51:
		return "51+"
	default:
		return "unknown"
	}
}
