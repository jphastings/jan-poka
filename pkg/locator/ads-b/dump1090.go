// +build !sbs

package ads_b

import (
	"encoding/json"
	"fmt"
	. "github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/math"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var _ LocationProvider = (*locationProvider)(nil)

type locationProvider struct {
	httpClient  *http.Client
	endpoint    string
	focusFlight string
}

const dataEndpoint = "/dump1090/data.json"

func init() {
	positionQuery, err := newClient("192.168.86.137:8090")
	if err != nil {
		log.Printf("❌ Provider: ADS-B unavailable, %s\n", err.Error())
	} else {
		Providers[TYPE] = func() LocationProvider { return positionQuery }
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

func (lp *locationProvider) Location() (TargetDetails, bool, error) {
	resp, err := http.Get(lp.endpoint)
	if err != nil {
		return TargetDetails{}, false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TargetDetails{}, false, err
	}

	var results []dataLine
	err = json.Unmarshal(body, &results)
	if err != nil {
		return TargetDetails{}, false, err
	}

	var flights []string
	for _, flight := range results {
		if flight.Flight == "" {
			continue
		}
		f := strings.Trim(flight.Flight, " ")
		flights = append(flights, f)
		if f == lp.focusFlight {
			details := TargetDetails{
				Name: "Flight " + f,
				Coords: math.LLACoords{
					Altitude:  math.Meters(flight.Altitude * 0.3048),
					Latitude:  math.Degrees(flight.Lat),
					Longitude: math.Degrees(flight.Lon),
				},
				// TODO: Can I get a more accurate time reading?
				AccurateAt: time.Now(),
			}

			return details, true, nil
		}
	}

	fmt.Println(strings.Join(flights, ", "))

	return TargetDetails{}, true, err
}

type dataLine struct {
	Flight   string  `json:"flight"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Altitude float64 `json:"altitude"`
}

func newClient(host string) (*locationProvider, error) {
	return nil, fmt.Errorf("can't easily detect unavailability - forcing off for now")

	u := (&url.URL{
		Scheme: "http",
		Host:   host,
		Path:   dataEndpoint,
	}).String()

	client := &http.Client{Timeout: time.Second}
	_, err := http.Head(u)
	if err != nil {
		return nil, err
	}
	return &locationProvider{endpoint: u, httpClient: client}, nil
}
