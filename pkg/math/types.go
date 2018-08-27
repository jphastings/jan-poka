package math

import (
	"math"
)

const Π = math.Pi

type Degrees float64
type Radians float64
type Meters float64

type Direction struct {
	Distance  Meters
	Heading   Degrees
	Elevation Degrees
}

type LLACoords struct {
	Φ Degrees
	Λ Degrees
	A Meters
}

type SphericalCoords struct {
	R Meters
	Θ Radians
	Φ Radians
}

type ECEFCoords struct {
	X Meters
	Y Meters
	Z Meters
}

func (deg Degrees) Radians() Radians {
	return Radians(deg * math.Pi / 180.0)
}

func (rad Radians) Degrees() Degrees {
	return Degrees(rad * 180 / math.Pi)
}

func Sin(rad Radians) Meters  { return Meters(math.Sin(float64(rad))) }
func Sinº(deg Degrees) Meters { return Sin(deg.Radians()) }
func Cos(rad Radians) Meters  { return Meters(math.Cos(float64(rad))) }
func Cosº(deg Degrees) Meters { return Cos(deg.Radians()) }

func Acos(a Meters) Radians      { return Radians(math.Acos(float64(a))) }
func Acosº(a Meters) Degrees     { return Acos(a).Degrees() }
func Atan2(y, x Meters) Radians  { return Radians(math.Atan2(float64(y), float64(x))) }
func Atan2º(y, x Meters) Degrees { return Atan2(y, x).Degrees() }

func Sqrt(a Meters) Meters    { return Meters(math.Sqrt(float64(a))) }
func Mod(v ECEFCoords) Meters { return Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z) }

// ModRad wraps the given number of radians to be within 0 and 2π
func ModRad(rad Radians) Radians {
	if rad < 0 {
		rad += 2 * Π
	}
	if rad >= 2*Π {
		rad -= 2 * Π
	}
	return rad
}

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
