package common

import (
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/math"
	"time"
)

// OnTracked is a function which can be called when a tracking is complete
type OnTracked func(TrackedDetails) future.Future

type TrackedDetails struct {
	Name                 string
	AccurateAt           time.Time
	Target               math.LLACoords
	Bearing              math.AERCoords
	UnobstructedDistance math.Meters
	IsFirstTrack         bool
}
