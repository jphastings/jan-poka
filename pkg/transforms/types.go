package transforms

import "math"

type degrees float64
type radians float64
type meters float64

type LLACoords struct {
	Φ degrees
	Λ degrees
	A meters
}

type SphericalCoords struct {
	R meters
	Θ radians
	Φ radians
}

type ECEFCoords struct {
	X meters
	Y meters
	Z meters
}

func Rad(deg degrees) radians {
	return radians(deg * math.Pi / 180.0)
}

func deg(rad radians) degrees {
	return degrees(rad * 180 / math.Pi)
}

func sin(rad radians) meters  { return meters(math.Sin(float64(rad))) }
func sinº(deg degrees) meters { return sin(Rad(deg)) }
func cos(rad radians) meters  { return meters(math.Cos(float64(rad))) }
func cosº(deg degrees) meters { return cos(Rad(deg)) }

func acos(a meters) radians      { return radians(math.Acos(float64(a))) }
func acosD(a meters) degrees     { return deg(acos(a)) }
func atan2(y, x meters) radians  { return radians(math.Atan2(float64(y), float64(x))) }
func atan2D(y, x meters) degrees { return deg(atan2(y, x)) }

func mod(v ECEFCoords) meters { return meters(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z))) }
