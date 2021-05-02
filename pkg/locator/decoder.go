package locator

import (
	"encoding/json"
	"fmt"
	"github.com/jphastings/jan-poka/pkg/common"
	"time"

	_ "github.com/jphastings/jan-poka/pkg/locator/ads-b"
	_ "github.com/jphastings/jan-poka/pkg/locator/celestial"
	_ "github.com/jphastings/jan-poka/pkg/locator/deliveroo"
	_ "github.com/jphastings/jan-poka/pkg/locator/instagram"
	_ "github.com/jphastings/jan-poka/pkg/locator/iss"
	_ "github.com/jphastings/jan-poka/pkg/locator/lla"
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
	sequence   []func() (common.TargetDetails, bool)
	Requester  common.OnTracked
}

func DecodeJSON(givenJSON []byte) (*TargetInstructions, error) {
	var target targetJSON
	err := json.Unmarshal(givenJSON, &target)
	if err != nil {
		return nil, err
	}

	ti := &TargetInstructions{
		sequence: []func() (common.TargetDetails, bool){},
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

		ti.sequence = append(ti.sequence, func() (common.TargetDetails, bool) {
			details, retry, err := prov.Location()
			_ = retry // TODO: Stop trying on subsequents if false && err
			return details, err == nil
		})
	}

	return ti, nil
}

func (ti *TargetInstructions) Poll() <-chan common.TargetDetails {
	locationsChan := make(chan common.TargetDetails)
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
