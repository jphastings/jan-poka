package common

import (
	"time"

	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/math"
)

// OnTracked is a function which can be called when a tracking is complete
type OnTracked func(TrackedDetails) future.Future

type TrackedDetails struct {
	Name                 string
	AccurateAt           time.Time
	Target               math.LLACoords
	Bearing              math.AERCoords
	UnobstructedDistance math.Meters
	LocalTime            time.Time

	IsFirstTrack bool
}
