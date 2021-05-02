package locator

import (
	"encoding/json"
	"fmt"
	. "github.com/jphastings/jan-poka/pkg/common"
	"time"
)

type targetJSON struct {
	PollSeconds   int               `json:"poll"`
	LocationSpecs []json.RawMessage `json:"target"`
}

type deciderLocationSpec struct {
	Type string `json:"type"`
}

type TargetInstructions struct {
	pollTicker *time.Ticker
	sequence   []func() (TargetDetails, bool, error)
	Requester  OnTracked
}

func DecodeJSON(givenJSON []byte) (*TargetInstructions, error) {
	var target targetJSON
	err := json.Unmarshal(givenJSON, &target)
	if err != nil {
		return nil, err
	}

	ti := &TargetInstructions{
		sequence: []func() (TargetDetails, bool, error){},
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

		init, ok := Providers[decider.Type]
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

		ti.sequence = append(ti.sequence, prov.Location)
	}

	return ti, nil
}

func (ti *TargetInstructions) Poll() <-chan TargetDetails {
	locationsChan := make(chan TargetDetails)
	go func() {
		for {
			i := 0
			for _, locationRetriever := range ti.sequence {
				target, retry, err := locationRetriever()
				if retry {
					ti.sequence[i] = locationRetriever
					i++
				}

				if err == nil {
					locationsChan <- target
					break
				}
			}

			// Unset all erased items
			for j := i; j < len(ti.sequence); j++ {
				ti.sequence[j] = nil
			}
			ti.sequence = ti.sequence[:i]

			// Do the next tick, if there's something to poll
			if ti.pollTicker == nil || len(ti.sequence) == 0 {
				break
			}
			<-ti.pollTicker.C
		}
	}()

	return locationsChan
}
