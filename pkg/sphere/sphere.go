package sphere

import (
	"fmt"
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
	currentΘ       Degrees
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
		currentΘ:       0,
	}

	return config
}

func (s *Config) TrackerCallback(_ string, bearing AERCoords, _ bool) chan error {
	return s.StepToDirection(bearing)
}

func (s *Config) StepToDirection(bearing AERCoords) chan error {
	promise := make(chan error)

	go func() {
		if err := s.powerSaver.PowerOn(); err != nil {
			promise <- err
			return
		}

		theta := 90 - bearing.Elevation
		finalTheta := theta

		if bearing.Azimuth == s.currentAzimuth {
			theta = s.currentΘ - theta
		} else {
			if err := <-s.stepHome(); err != nil {
				s.powerSaver.PowerOff()
				promise <- err
				return
			}
		}

		if theta != 0 {
			if err := <-s.stepToΘ(bearing.Azimuth, theta); err != nil {
				s.powerSaver.PowerOff()
				promise <- err
				return
			}
		}

		s.powerSaver.PowerOff()
		s.currentΘ = finalTheta
		s.currentAzimuth = bearing.Azimuth

		promise <- nil
	}()

	return promise
}

// Home is at Θ = 0 (straight up)
func (s *Config) stepHome() chan error {
	oppositeHeading := 180 + s.currentAzimuth
	if oppositeHeading >= 360 {
		oppositeHeading -= 360
	}

	return s.stepToΘ(oppositeHeading, s.currentΘ)
}

func (s *Config) stepToΘ(heading, theta Degrees) chan error {
	promise := make(chan error)

	if theta == 0 {
		promise <- nil
		return promise
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
		motorSteps := -int(math.Ceil(float64(Cosº(heading-mtr.Angle)) * maxSteps))
		wg.Add(1)
		go func() {
			travelMotor(travelTime, mtr, motorSteps)
			wg.Done()
		}()
	}

	wg.Wait()

	promise <- nil
	return promise
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
