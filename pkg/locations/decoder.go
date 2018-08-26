package locations

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jphastings/corviator/pkg/locations/iss"
	"github.com/jphastings/corviator/pkg/locations/lla"
	"github.com/jphastings/corviator/pkg/transforms"
)

type locationProvider interface {
	// SetParams provides you with a function that will populate the given (json annotated) struct pointer with the JSON params. May return error if JSON does not match. Should validate given parameters and return error if unusable.
	SetParams(func(decodeInto interface{}) error) error
	// Location returns the location according to the params set earlier at the current time. Second argument can be false if provider is currently offline.
	Location() (target transforms.LLACoords, isUsable bool)
}

var locationProviders map[string]locationProvider

type targetJSON struct {
	PollSeconds   int               `json:"poll"`
	LocationSpecs []json.RawMessage `json:"target"`
}

type deciderLocationSpec struct {
	Type string `json:"type"`
}

type TargetInstructions struct {
	pollTicker *time.Ticker
	sequence   []func() (transforms.LLACoords, bool)
}

func init() {
	locationProviders = map[string]locationProvider{
		iss.TYPE: iss.NewLocationProvider(),
		lla.TYPE: lla.NewLocationProvider(),
	}
}

func DecodeJSON(givenJSON []byte) (*TargetInstructions, error) {
	var target targetJSON
	err := json.Unmarshal(givenJSON, &target)
	if err != nil {
		return nil, err
	}

	ti := &TargetInstructions{
		sequence:   []func() (transforms.LLACoords, bool){},
		pollTicker: time.NewTicker(time.Duration(target.PollSeconds) * time.Second),
	}

	for _, locationSpec := range target.LocationSpecs {
		var decider deciderLocationSpec
		err = json.Unmarshal(locationSpec, &decider)
		if err != nil {
			return nil, err
		}

		provider, ok := locationProviders[decider.Type]
		if !ok {
			return nil, fmt.Errorf("unknown provider: %s", decider.Type)
		}

		err := provider.SetParams(func(params interface{}) error {
			err = json.Unmarshal(locationSpec, &params)
			return err
		})
		if err != nil {
			return nil, err
		}

		ti.sequence = append(ti.sequence, func() (transforms.LLACoords, bool) {
			return provider.Location()
		})
	}

	return ti, nil
}

func (ti *TargetInstructions) Poll() <-chan transforms.LLACoords {
	locationsChan := make(chan transforms.LLACoords)
	go func() {
		for {
			for _, locationRetriever := range ti.sequence {
				if location, ok := locationRetriever(); ok {
					locationsChan <- location
					break
				}
			}
			<-ti.pollTicker.C
		}
	}()

	return locationsChan
}
