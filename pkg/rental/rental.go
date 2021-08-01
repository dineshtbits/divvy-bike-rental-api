package rental

import (
	"encoding/csv"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const tripsFilePath = "../resources/Divvy_Trips_2019_Q2"

type RiderSummaryRequest struct {
	Filters RiderSummaryRequestFilters `json:"filters"`
}

type RiderSummaryRequestFilters struct {
	StationIds []int `json:"station_ids"`
}

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

func LoadTripsData() (*[]Rental, error) {
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

func GetTripsSummary(c *gin.Context, trips *[]Rental) {
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

func GetRidersSummary(c *gin.Context, trips *[]Rental) {

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
