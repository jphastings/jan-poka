package sphere

import (
	"github.com/jphastings/corviator/pkg/future"
	"github.com/jphastings/corviator/pkg/math"
)

func (s *Config) TrackerCallback(_ string, bearing math.AERCoords, _ bool) future.Future {
	return s.StepToDirection(bearing)
}
