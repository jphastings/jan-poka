package locations

import (
	"encoding/json"
	"time"

	"github.com/jphastings/corviator/pkg/locations/iss"
	"github.com/jphastings/corviator/pkg/transforms"
)

type locationProvider interface {
	Location(interface{}) (transforms.LLACoords, bool)
	EmptyParams() interface{}
}

var locationProviders map[string]locationProvider

type targetJSON struct {
	PollSeconds   int               `json:"poll"`
	LocationSpecs []json.RawMessage `json:"target"`
}

type deciderLocationSpec struct {
	Type string `json:"type"`
}

func init() {
	locationProviders = map[string]locationProvider{
		"iss": iss.NewLocationProvider(),
	}
}

func DecodeJSON(givenJSON []byte) (<-chan transforms.LLACoords, error) {
	var target targetJSON
	err := json.Unmarshal(givenJSON, &target)
	if err != nil {
		return nil, err
	}

	targetSequence := []func() (transforms.LLACoords, bool){}
	for _, locationSpec := range target.LocationSpecs {
		var decider deciderLocationSpec
		err = json.Unmarshal(locationSpec, &decider)
		if err != nil {
			return nil, err
		}

		provider, ok := locationProviders[decider.Type]
		if !ok {
			// TODO: Log issue
			continue
		}

		params := provider.EmptyParams()
		err = json.Unmarshal(locationSpec, &params)
		if err != nil {
			return nil, err
		}

		targetSequence = append(targetSequence, func() (transforms.LLACoords, bool) {
			return provider.Location(params)
		})
	}

	wait := time.Duration(target.PollSeconds) * time.Second
	locationsChan := make(chan transforms.LLACoords)
	go func() {
		for {
			for _, locationRetreiver := range targetSequence {
				location, ok := locationRetreiver()
				if ok {
					locationsChan <- location
					break
				}
			}
			time.Sleep(wait)
		}
	}()

	return locationsChan, nil
}
