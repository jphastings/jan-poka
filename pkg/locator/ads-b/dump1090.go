// +build !sbs

package ads_b

import (
	"encoding/json"
	"fmt"
	. "github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/math"
	"io/ioutil"
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

func Connect(dump1090Host string) error {
	positionQuery, err := newClient(dump1090Host)
	if err != nil {
		return err
	}
	Providers[TYPE] = func() LocationProvider { return positionQuery }
	return nil
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
	u := (&url.URL{
		Scheme: "http",
		Host:   host,
		Path:   dataEndpoint,
	}).String()

	client := &http.Client{Timeout: time.Second}

	_, err := client.Head(u)
	if err != nil {
		return nil, fmt.Errorf("the dump1090 server is not running")
	}

	return &locationProvider{endpoint: u, httpClient: client}, nil
}
