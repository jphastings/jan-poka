package tower

import (
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/math"
)

func (s *Config) TrackerCallback(_ string, _ math.LLACoords, bearing math.AERCoords, distance math.Meters, _ bool) future.Future {
	fn := future.Exec(func() error { return s.setDirection(bearing) })
	fe := future.Exec(func() error { return s.setDistance(distance) })
	return future.All(fn, fe)
}
