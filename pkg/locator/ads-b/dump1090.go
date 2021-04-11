// +build !sbs

package ads_b

import (
	"encoding/json"
	"fmt"
	"github.com/jphastings/jan-poka/pkg/locator/common"
	"github.com/jphastings/jan-poka/pkg/math"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type locationProvider struct {
	host        string
	focusFlight string
}

func init() {
	positionQuery, err := newClient("192.168.86.137:8090")
	if err != nil {
		log.Printf("❌ Provider: ADS-B unavailable, %s\n", err.Error())
	} else {
		common.Providers[TYPE] = func() common.LocationProvider { return positionQuery }
		log.Println("✅ Provider: ADS-B airplane positions available.")
	}
}

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	loc := &params{}
	err := decodeInto(loc)
	if err != nil {
		return err
	}
	lp.focusFlight = loc.Flight
	return nil
}

func (lp *locationProvider) Location() (math.LLACoords, string, bool) {
	resp, err := http.Get(fmt.Sprintf("http://%s/dump1090/data.json", lp.host))
	if err != nil {
		return math.LLACoords{}, "", false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return math.LLACoords{}, "", false
	}

	var results []dataLine
	err = json.Unmarshal(body, &results)
	if err != nil {
		panic(err)
		return math.LLACoords{}, "", false
	}

	var flights []string
	for _, flight := range results {
		if flight.Flight == "" {
			continue
		}
		f := strings.Trim(flight.Flight, " ")
		flights = append(flights, f)
		if f == lp.focusFlight {
			return math.LLACoords{
				Altitude:  math.Meters(flight.Altitude * 0.3048),
				Latitude:  math.Degrees(flight.Lat),
				Longitude: math.Degrees(flight.Lon),
			}, "Flight " + f, true
		}
	}

	fmt.Println(strings.Join(flights, ", "))

	return math.LLACoords{}, "", false
}

type dataLine struct {
	Flight   string  `json:"flight"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Altitude float64 `json:"altitude"`
}

func newClient(host string) (*locationProvider, error) {
	return &locationProvider{host: host}, nil
}
