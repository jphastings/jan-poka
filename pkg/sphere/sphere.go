package sphere

import (
	"github.com/jphastings/corviator/pkg/hardware/motor"
	. "github.com/jphastings/corviator/pkg/math"
	"math"
	"periph.io/x/periph/conn/gpio"
	"time"
)

type Config struct {
	motors     []*motor.Motor
	powerSaver *motor.PowerSaver

	sphereRotationSteps float64
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

		sphereRotationSteps: wheelRatio * float64(wheelRotationSteps),
		minStepInterval:     minStepInterval,
		Facing:              facing,

		currentAzimuth: 0,
		currentΘ:       0,
	}

	return config
}

func (s *Config) TrackerCallback(_ string, bearing AERCoords, _ bool) error {
	s.StepToDirection(bearing)
	return nil
}

func (s *Config) StepToDirection(bearing AERCoords) time.Duration {
	Θ := 90 - bearing.Elevation
	finalΘ := Θ
	completesIn := time.Duration(0)

	if bearing.Azimuth == s.currentAzimuth {
		Θ = s.currentΘ - Θ
	} else {
		completesIn = s.stepHome(completesIn)
	}

	if Θ != 0 {
		completesIn = s.stepToΘ(bearing.Azimuth, Θ, completesIn)
	}

	s.powerSaver.PowerUntil(completesIn)
	finished := time.NewTimer(completesIn)
	go func() {
		<-finished.C
		s.currentΘ = finalΘ
		s.currentAzimuth = bearing.Azimuth
	}()

	return completesIn
}

// Home is at Θ = 0 (straight up)
func (s *Config) stepHome(wait time.Duration) time.Duration {
	oppositeHeading := 180 + s.currentAzimuth
	if oppositeHeading >= 360 {
		oppositeHeading -= 360
	}

	return s.stepToΘ(oppositeHeading, s.currentΘ, wait)
}

func (s *Config) stepToΘ(heading, Θ Degrees, wait time.Duration) time.Duration {
	if Θ == 0 {
		return wait
	}

	maxSteps := float64(Θ) * s.sphereRotationSteps / 360
	travelTime := time.Duration(maxSteps) * s.minStepInterval
	if travelTime < 0 {
		travelTime = -travelTime
	}

	for _, mtr := range s.motors {
		motorSteps := -int(math.Ceil(float64(Cosº(heading-mtr.Angle)) * maxSteps))
		go travelMotor(wait, travelTime, mtr, motorSteps)
	}

	// We assume the execution time of this function is negligible
	return wait + travelTime
}

func travelMotor(w, t time.Duration, m *motor.Motor, s int) {
	<-time.NewTimer(w).C

	var f bool
	if s >= 0 {
		f = true
	} else {
		s = -s
		f = false
	}

	ticker := time.NewTicker(t / time.Duration(s))
	for range ticker.C {
		m.StepChannel <- f

		s--
		if s < 0 {
			break
		}
	}
}
