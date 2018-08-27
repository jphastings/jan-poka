package sphere

import (
	"github.com/jphastings/corviator/pkg/hardware/wheel"
	. "github.com/jphastings/corviator/pkg/math"
	"math"
	"time"
)

type Config struct {
	motors []*wheel.Motor

	sphereRotationSteps float64
	minStepInterval     time.Duration

	Facing float64

	isSetUp        bool
	currentHeading Degrees
	currentΘ       Degrees
}

func New(
	motors []*wheel.Motor,
	wheelRotationSteps int,
	wheelRatio float64,
	minStepInterval time.Duration,
	facing float64,
) *Config {
	return &Config{
		motors:              motors,
		sphereRotationSteps: wheelRatio * float64(wheelRotationSteps),
		minStepInterval:     minStepInterval,
		Facing:              facing,

		currentHeading: 0,
		currentΘ:       0,
	}
}

func (s *Config) StepToElevation(heading, elevation Degrees) time.Duration {
	Θ := 90 - elevation
	finalΘ := Θ
	completesIn := time.Duration(0)

	if heading == s.currentHeading {
		Θ = s.currentΘ - Θ
	} else {
		completesIn = s.stepHome(completesIn)
	}

	if Θ != 0 {
		completesIn = s.stepToΘ(heading, Θ, completesIn)
	}

	finished := time.NewTimer(completesIn)
	go func() {
		<-finished.C
		s.currentΘ = finalΘ
		s.currentHeading = heading
	}()

	return completesIn
}

// Home is at Θ = 0 (strait up)
func (s *Config) stepHome(wait time.Duration) time.Duration {
	oppositeHeading := 180 + s.currentHeading
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

	for _, motor := range s.motors {
		motorSteps := -int(math.Ceil(float64(Cosº(heading-motor.Angle)) * maxSteps))
		go travelMotor(wait, travelTime, motor, motorSteps)
	}

	// We assume the execution time of this function is negligible
	return wait + travelTime
}

func travelMotor(w, t time.Duration, m *wheel.Motor, s int) {
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
