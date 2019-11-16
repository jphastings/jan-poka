package sphere

import (
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/math"
)

func (s *Config) TrackerCallback(_ string, bearing math.AERCoords, _ bool) future.Future {
	return s.StepToDirection(bearing)
}
