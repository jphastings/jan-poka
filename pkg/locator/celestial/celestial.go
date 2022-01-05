//go:build libnova

package celestial

import (
	"fmt"
	"log"
	"strings"
	"time"

	. "github.com/jphastings/jan-poka/pkg/common"
)

const TYPE = "celestial"

var _ LocationProvider = (*locationProvider)(nil)

type locationProvider struct {
	target Body
}

type config struct {
	Body string `json:"body"`
}

func init() {
	Providers[TYPE] = func() LocationProvider { return &locationProvider{} }
	log.Println("✅ Provider: Celestial body positions")
}

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	c := &config{}
	if err := decodeInto(c); err != nil {
		return err
	}

	body, ok := Bodies[c.Body]
	if !ok {
		return fmt.Errorf("unknown celestial body '%s'", c.Body)
	}
	lp.target = body
	return nil
}

func (lp *locationProvider) Location() TargetDetails {
	name := strings.Title(string(lp.target))
	if lp.target == Moon || lp.target == Sun {
		name = "The " + name
	}

	at := time.Now()
	loc, err := GeocentricCoordinates(lp.target, at)

	return TargetDetails{
		Name:       name,
		Coords:     loc,
		AccurateAt: at,
		Final:      false,
		Err:        err,
	}
}
