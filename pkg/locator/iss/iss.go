package iss

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

type locationProvider struct{}

func init() {
	common.Providers[TYPE] = func() common.LocationProvider { return &locationProvider{} }
	log.Println("✅ Provider: International Space Station positions available.")
}

func (_ *locationProvider) SetParams(func(interface{}) error) error { return nil }

func (_ *locationProvider) Location() (math.LLACoords, string, bool) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://api.open-notify.org/iss-now.json", nil)
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return math.LLACoords{}, "", false
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var response serviceResponse
	err = decoder.Decode(&response)
	if err != nil {
		return math.LLACoords{}, "", false
	}

	latitude, err := strconv.ParseFloat(response.Position.Latitude, 64)
	if err != nil {
		return math.LLACoords{}, "", false
	}
	longitude, err := strconv.ParseFloat(response.Position.Longitude, 64)
	if err != nil {
		return math.LLACoords{}, "", false
	}

	return math.LLACoords{
		Latitude:  math.Degrees(latitude),
		Longitude: math.Degrees(longitude),
		Altitude:  408000,
	}, "The International Space Station", true
}
