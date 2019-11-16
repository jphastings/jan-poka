package locator

import (
	"encoding/json"
	"fmt"
	"github.com/jphastings/jan-poka/pkg/locator/common"
	"github.com/jphastings/jan-poka/pkg/locator/instagram"
	"github.com/jphastings/jan-poka/pkg/locator/iss"
	"github.com/jphastings/jan-poka/pkg/locator/lla"
	"time"

	"github.com/jphastings/jan-poka/pkg/math"
)

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

func init() {
	// Load locators that have no external dependencies
	lla.Load()
	iss.Load()
	instagram.Load()
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

		init, ok := common.Providers[decider.Type]
		if !ok {
			return nil, fmt.Errorf("unknown provider: %s", decider)
		}
		prov := init()

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
