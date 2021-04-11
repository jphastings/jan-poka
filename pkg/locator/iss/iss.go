package iss

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jphastings/jan-poka/pkg/locator/common"
	"github.com/jphastings/jan-poka/pkg/math"
)

const TYPE = "iss"

type serviceResponse struct {
	Position struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"iss_position"`
}

var _ common.LocationProvider = (*locationProvider)(nil)

type locationProvider struct{}

func init() {
	common.Providers[TYPE] = func() common.LocationProvider { return &locationProvider{} }
	log.Println("✅ Provider: International Space Station positions available.")
}

func (_ *locationProvider) SetParams(func(interface{}) error) error { return nil }

func (_ *locationProvider) Location() (math.LLACoords, time.Time, string, bool) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://api.open-notify.org/iss-now.json", nil)
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return math.LLACoords{}, time.Time{}, "", false
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var response serviceResponse
	err = decoder.Decode(&response)
	if err != nil {
		return math.LLACoords{}, time.Time{}, "", false
	}

	latitude, err := strconv.ParseFloat(response.Position.Latitude, 64)
	if err != nil {
		return math.LLACoords{}, time.Time{}, "", false
	}
	longitude, err := strconv.ParseFloat(response.Position.Longitude, 64)
	if err != nil {
		return math.LLACoords{}, time.Time{}, "", false
	}

	return math.LLACoords{
		Latitude:  math.Degrees(latitude),
		Longitude: math.Degrees(longitude),
		Altitude:  408000,
	}, time.Now(), "The International Space Station", true
}
