package sphere

import (
	"fmt"
	"github.com/jphastings/corviator/pkg/future"
	"github.com/jphastings/corviator/pkg/hardware/motor"
	. "github.com/jphastings/corviator/pkg/math"
	"math"
	"periph.io/x/periph/conn/gpio"
	"sync"
	"time"
)

type Config struct {
	motors     []*motor.Motor
	powerSaver *motor.PowerSaver

	sphereRotationSteps int
	minStepInterval     time.Duration

	Facing Degrees

	isSetUp        bool
	currentAzimuth Degrees
	currentTheta   Degrees
}

func New(
	motors []*motor.Motor,
	motorsActivePin gpio.PinOut,
	motorsActiveLeeway time.Duration,
	wheelRotationSteps int,
	wheelRatio float64,
	minStepInterval time.Duration,
	facing Degrees,
) *Config {
	config := &Config{
		motors:     motors,
		powerSaver: motor.NewPowerSaver(motorsActivePin, motorsActiveLeeway),

		sphereRotationSteps: int(wheelRatio) * wheelRotationSteps,
		minStepInterval:     minStepInterval,
		Facing:              facing,

		currentAzimuth: 0,
		currentTheta:   0,
	}

	return config
}

func (s *Config) StepToDirection(bearing AERCoords) future.Future {
	f := future.New()

	go func() {
		if err := s.powerSaver.PowerOn(); err != nil {
			f.Fail(err)
			return
		}

		theta := 90 - bearing.Elevation
		finalTheta := theta

		if bearing.Azimuth == s.currentAzimuth {
			theta = s.currentTheta - theta
		} else {
			if result := <-s.stepHome(); !result.IsOK() {
				s.powerSaver.PowerOff()
				f.Bubble(result)
				return
			}
		}

		if theta != 0 {
			if result := <-s.stepToTheta(bearing.Azimuth, theta); !result.IsOK() {
				s.powerSaver.PowerOff()
				f.Bubble(result)
				return
			}
		}

		s.powerSaver.PowerOff()
		s.currentTheta = finalTheta
		s.currentAzimuth = bearing.Azimuth

		f.Succeed()
	}()

	return f
}

// Home is at Θ = 0 (straight up)
func (s *Config) stepHome() future.Future {
	oppositeHeading := 180 + s.currentAzimuth
	if oppositeHeading >= 360 {
		oppositeHeading -= 360
	}

	return s.stepToTheta(oppositeHeading, s.currentTheta)
}

func (s *Config) stepToTheta(heading, theta Degrees) future.Future {
	f := future.New()

	if theta == 0 {
		f.Succeed()
		return f
	}

	maxSteps := float64(theta) * float64(s.sphereRotationSteps) / 360
	fmt.Println("Maximum steps is", maxSteps)
	travelTime := time.Duration(maxSteps) * s.minStepInterval
	fmt.Println("Travel time is", travelTime.String())
	if travelTime < 0 {
		travelTime = -travelTime
	}

	var wg sync.WaitGroup

	for _, mtr := range s.motors {
		motorSteps := -int(math.Ceil(float64(CosDeg(heading-mtr.Angle)) * maxSteps))
		wg.Add(1)
		go func() {
			travelMotor(travelTime, mtr, motorSteps)
			wg.Done()
		}()
	}

	wg.Wait()

	f.Succeed()
	return f
}

func travelMotor(t time.Duration, m *motor.Motor, s int) {
	var f bool
	if s >= 0 {
		f = true
	} else {
		s = -s
		f = false
	}

	pulseWidth := t / time.Duration(s)
	ticker := time.NewTicker(pulseWidth)
	for range ticker.C {
		m.StepChannel <- f

		s--
		if s < 0 {
			break
		}
	}
}
