package math

import m "math"

func (deg Degrees) Radians() Radians {
	return Radians(deg * Π / 180.0)
}

func (rad Radians) Degrees() Degrees {
	return Degrees(rad * 180 / Π)
}

func Sin(rad Radians) Meters  { return Meters(m.Sin(float64(rad))) }
func Sinº(deg Degrees) Meters { return Sin(deg.Radians()) }
func Cos(rad Radians) Meters  { return Meters(m.Cos(float64(rad))) }
func Cosº(deg Degrees) Meters { return Cos(deg.Radians()) }

func Acos(a Meters) Radians      { return Radians(m.Acos(float64(a))) }
func Acosº(a Meters) Degrees     { return Acos(a).Degrees() }
func Atan2(y, x Meters) Radians  { return Radians(m.Atan2(float64(y), float64(x))) }
func Atan2º(y, x Meters) Degrees { return Atan2(y, x).Degrees() }

func Sqrt(a Meters) Meters { return Meters(m.Sqrt(float64(a))) }

// ModDeg wraps the given number of degrees to be within 0 and 360
func ModDeg(deg Degrees) Degrees {
	if deg < 0 {
		deg += 360
	}
	if deg >= 360 {
		deg -= 360
	}
	return deg
}
