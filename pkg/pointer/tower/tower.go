package tower

import (
	"github.com/jphastings/jan-poka/pkg/math"
	"log"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/pca9685"
)

type Config struct {
	facing math.Degrees

	// The tower is expected to be facing the face of the device when turned on
	thetaServo *pca9685.Servo
	// The arrow is expected to be pointing directly downwards when turned on
	phiServo *pca9685.Servo
	// The needle is expected to be pointing to zero when turned on
	distanceNumeral *pca9685.Servo
	// The multiplier is intended to be pointing at "by foot" when turned on
	distanceScale *pca9685.Servo
}

// Expects the Adafruit Servo controller board: https://learn.adafruit.com/16-channel-pwm-servo-driver/overview
func New(bus i2c.Bus, facing math.Degrees) (*Config, error) {
	dev, err := pca9685.NewI2C(bus, pca9685.I2CAddr)
	if err != nil {
		return nil, err
	}

	if err := dev.SetPwmFreq(50 * physic.Hertz); err != nil {
		log.Fatal(err)
	}
	if err := dev.SetAllPwm(0, 0); err != nil {
		log.Fatal(err)
	}

	servos := pca9685.NewServoGroup(dev, 50, 650, -360, +360)

	config := &Config{
		thetaServo:      servos.GetServo(0),
		phiServo:        servos.GetServo(1),
		distanceNumeral: servos.GetServo(2),
		distanceScale:   servos.GetServo(3),

		facing: facing,
	}

	config.thetaServo.SetMinMaxAngle(-180, +180)
	config.distanceNumeral.SetMinMaxAngle(-90, +90)

	return config, nil
}

func (s *Config) setDirection(bearing math.AERCoords) error {
	degreesFromAhead := math.ModDeg(bearing.Azimuth - s.facing)
	var arrowRight bool

	var err error
	if degreesFromAhead <= 90 {
		arrowRight = true
		err = s.thetaServo.SetAngle((90 - degreesFromAhead).Angle())
	} else if degreesFromAhead <= 180 {
		arrowRight = true
		err = s.thetaServo.SetAngle((degreesFromAhead - 90).Angle())
	} else if degreesFromAhead <= 270 {
		arrowRight = false
		err = s.thetaServo.SetAngle((270 - degreesFromAhead).Angle())
	} else {
		arrowRight = false
		err = s.thetaServo.SetAngle((degreesFromAhead - 270).Angle())
	}
	if err != nil {
		return err
	}

	phi := bearing.Elevation + 90
	if arrowRight {
		phi *= -1
	}

	if err := s.phiServo.SetAngle(phi.Angle()); err != nil {
		return err
	}

	return nil
}

const maxNumeral = 15

func (s *Config) setDistance(distance math.Meters) error {
	numeral, scale := math.Scaled(distance)

	var err error
	switch scale {
	case math.ByFoot:
		err = s.distanceScale.SetAngle(0 * physic.Degree)
	case math.ByCar:
		err = s.distanceScale.SetAngle(120 * physic.Degree)
	case math.ByPlane:
		err = s.distanceScale.SetAngle(240 * physic.Degree)
	}
	if err != nil {
		return err
	}

	angle := (numeral / maxNumeral) * 180
	if angle > 180 {
		angle = 180
	}

	err = s.distanceScale.SetAngle(physic.Angle(angle) * physic.Degree)
	if err != nil {
		return err
	}

	return nil
}
