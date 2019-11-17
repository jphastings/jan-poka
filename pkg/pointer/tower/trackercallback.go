package tower

import (
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/math"
)

func (s *Config) TrackerCallback(_ string, bearing math.AERCoords, distance math.Meters, _ bool) future.Future {
	fn := future.New()
	go func() {
		err := s.setDirection(bearing)
		if err != nil {
			fn.Fail(err)
		} else {
			fn.Succeed()
		}
	}()

	fe := future.New()
	go func() {
		err := s.setDistance(distance)
		if err != nil {
			fe.Fail(err)
		} else {
			fe.Succeed()
		}
	}()

	return future.All(fn, fe)
}
