package iss

import (
  "encoding/json"
  "net/http"
  "strconv"

  "github.com/jphastings/corviator/pkg/transforms"
)

type serviceResponse struct {
  Position struct {
    Latitude  string `json:"latitude"`
    Longitude string `json:"longitude"`
  } `json:"iss_position"`
}

type locationProvider struct{}
type params struct{}

func NewLocationProvider() *locationProvider {
  return &locationProvider{}
}

func (_ *locationProvider) EmptyParams() interface{} {
  return &params{}
}

func (_ *locationProvider) Location(_ interface{}) (transforms.LLACoords, bool) {
  client := &http.Client{}

  req, _ := http.NewRequest("GET", "http://api.open-notify.org/iss-now.json", nil)
  req.Header.Add("Accept", "application/json")
  res, err := client.Do(req)
  if err != nil {
    return transforms.LLACoords{}, false
  }

  defer res.Body.Close()
  decoder := json.NewDecoder(res.Body)
  var response serviceResponse
  err = decoder.Decode(&response)
  if err != nil {
    return transforms.LLACoords{}, false
  }

  φ, err := strconv.ParseFloat(response.Position.Latitude, 64)
  if err != nil {
    return transforms.LLACoords{}, false
  }
  λ, err := strconv.ParseFloat(response.Position.Longitude, 64)
  if err != nil {
    return transforms.LLACoords{}, false
  }

  return transforms.LLACoords{
    Φ: φ,
    Λ: λ,
    A: 408000,
  }, true
}
