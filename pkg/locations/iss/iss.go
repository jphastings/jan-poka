package iss

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jphastings/corviator/pkg/math"
)

const TYPE = "iss"

type serviceResponse struct {
	Position struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"iss_position"`
}

type locationProvider struct{}

func NewLocationProvider() *locationProvider {
	return &locationProvider{}
}

func (_ *locationProvider) SetParams(func(interface{}) error) error { return nil }

func (_ *locationProvider) Location() (math.LLACoords, bool) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "http://api.open-notify.org/iss-now.json", nil)
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return math.LLACoords{}, false
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var response serviceResponse
	err = decoder.Decode(&response)
	if err != nil {
		return math.LLACoords{}, false
	}

	φ, err := strconv.ParseFloat(response.Position.Latitude, 64)
	if err != nil {
		return math.LLACoords{}, false
	}
	λ, err := strconv.ParseFloat(response.Position.Longitude, 64)
	if err != nil {
		return math.LLACoords{}, false
	}

	return math.LLACoords{
		Φ: math.Degrees(φ),
		Λ: math.Degrees(λ),
		A: 408000,
	}, true
}
