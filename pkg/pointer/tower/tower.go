package tower

import (
	"github.com/jphastings/jan-poka/pkg/hardware/stepper"
	"github.com/jphastings/jan-poka/pkg/math"
	"sync"
)

type Config struct {
	facing math.Degrees

	// The tower is expected to be facing the face of the device when turned on
	thetaServo *stepper.Stepper
	// The arrow is expected to be pointing directly downwards when turned on
	phiServo *stepper.Stepper
	// The needle is expected to be pointing to zero when turned on
	distanceNumeral *stepper.Stepper
	// The multiplier is intended to be pointing at "by foot" when turned on
	distanceScale *stepper.Stepper
}

func New(facing math.Degrees) (*Config, error) {
	steppers := stepper.Pi2Quad(stepper.Motors["28BYJ-48"])

	steppers[0].SetSpeed(20)
	steppers[1].SetSpeed(20)
	steppers[2].SetSpeed(20)
	steppers[3].SetSpeed(20)

	config := &Config{
		thetaServo:      steppers[0],
		phiServo:        steppers[1],
		distanceNumeral: steppers[2],
		distanceScale:   steppers[3],

		facing: facing,
	}

	return config, nil
}

func (s *Config) Shutdown() {
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		s.thetaServo.SetAngle(0)
		s.thetaServo.Off()
		s.phiServo.SetAngle(0)
		s.phiServo.Off()
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		s.distanceNumeral.SetAngle(0)
		s.distanceNumeral.Off()
		s.distanceScale.SetAngle(0)
		s.distanceScale.Off()
		wg.Done()
	}()

	wg.Wait()
}

func (s *Config) setDirection(bearing math.AERCoords) error {
	base, arm := Pointer(math.ModDeg(bearing.Azimuth-s.facing), bearing.Elevation+90)

	if err := s.thetaServo.SetAngle(base); err != nil {
		return err
	}
	if err := s.thetaServo.Off(); err != nil {
		return err
	}

	if err := s.phiServo.SetAngle(arm); err != nil {
		return err
	}
	if err := s.phiServo.Off(); err != nil {
		return err
	}

	return nil
}

const maxNumeral = 15

func (s *Config) setDistance(distance math.Meters) error {
	numeral, scale := math.Scaled(distance)
	var err error

	angle := (numeral / maxNumeral) * 180
	if angle > 180 {
		angle = 180
	}

	err = s.distanceNumeral.SetAngle(math.Degrees(angle))
	if err != nil {
		return err
	}
	if err := s.distanceNumeral.Off(); err != nil {
		return err
	}

	switch scale {
	case math.ByFoot:
		err = s.distanceScale.SetAngle(0)
	case math.ByCar:
		err = s.distanceScale.SetAngle(120)
	case math.ByPlane:
		err = s.distanceScale.SetAngle(240)
	}
	if err != nil {
		return err
	}
	if err := s.distanceScale.Off(); err != nil {
		return err
	}

	return nil
}
