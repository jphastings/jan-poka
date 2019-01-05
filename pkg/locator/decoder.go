package locator

import (
	"encoding/json"
	"fmt"
	"github.com/jphastings/corviator/pkg/locator/celestial"
	"time"

	"github.com/jphastings/corviator/pkg/locator/iss"
	"github.com/jphastings/corviator/pkg/locator/lla"
	"github.com/jphastings/corviator/pkg/math"
)

type locationProvider interface {
	// SetParams provides you with a function that will populate the given (json annotated) struct pointer with the JSON params. May return error if JSON does not match. Should validate given parameters and return error if unusable.
	SetParams(func(decodeInto interface{}) error) error
	// Location returns the location according to the params set earlier at the current time. Second argument can be false if provider is currently offline.
	Location() (target math.LLACoords, suggestedName string, isUsable bool)
}

type targetJSON struct {
	PollSeconds   int               `json:"poll"`
	LocationSpecs []json.RawMessage `json:"target"`
}

type deciderLocationSpec struct {
	Type string `json:"type"`
}

type TargetDetails struct {
	Name   string
	Coords math.LLACoords
}

type TargetInstructions struct {
	pollTicker *time.Ticker
	sequence   []func() (TargetDetails, bool)
}

func provider(decider string) (locationProvider, error) {
	switch decider {
	case lla.TYPE:
		return lla.NewLocationProvider(), nil
	case iss.TYPE:
		return iss.NewLocationProvider(), nil
	case celestial.TYPE:
		return celestial.NewLocationProvider(), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", decider)
	}
}

func DecodeJSON(givenJSON []byte) (*TargetInstructions, error) {
	var target targetJSON
	err := json.Unmarshal(givenJSON, &target)
	if err != nil {
		return nil, err
	}

	ti := &TargetInstructions{
		sequence: []func() (TargetDetails, bool){},
	}

	if target.PollSeconds > 0 {
		ti.pollTicker = time.NewTicker(time.Duration(target.PollSeconds) * time.Second)
	}

	for _, locationSpec := range target.LocationSpecs {
		var decider deciderLocationSpec
		err = json.Unmarshal(locationSpec, &decider)
		if err != nil {
			return nil, err
		}

		prov, err := provider(decider.Type)
		if err != nil {
			return nil, err
		}

		err = prov.SetParams(func(params interface{}) error {
			err = json.Unmarshal(locationSpec, &params)
			return err
		})
		if err != nil {
			return nil, err
		}

		ti.sequence = append(ti.sequence, func() (TargetDetails, bool) {
			coords, name, isUsable := prov.Location()
			return TargetDetails{Name: name, Coords: coords}, isUsable
		})
	}

	return ti, nil
}

func (ti *TargetInstructions) Poll() <-chan TargetDetails {
	locationsChan := make(chan TargetDetails)
	go func() {
		for {
			for _, locationRetriever := range ti.sequence {
				target, ok := locationRetriever()
				if ok {
					locationsChan <- target
					break
				}
			}
			if ti.pollTicker == nil {
				break
			}
			<-ti.pollTicker.C
		}
	}()

	return locationsChan
}
