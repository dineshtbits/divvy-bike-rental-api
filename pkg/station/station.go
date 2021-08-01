package station

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const stationInfoUrl = "https://gbfs.divvybikes.com/gbfs/en/station_information.json"

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

func LoadStationsData() *StationsData {
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

func GetStationById(c *gin.Context, stations *StationsData) {
	for i := 0; i < len(stations.Data.Stations); i++ {
		if stations.Data.Stations[i].StationId == c.Param("id") {
			c.IndentedJSON(http.StatusOK, stations.Data.Stations[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "station not found"})
}
