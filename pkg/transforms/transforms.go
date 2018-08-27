package transforms

import (
	"math"
)

func Rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func deg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func RelativeDirection(source LLACoords, target LLACoords, facing float64) (distance, heading, elevation float64) {
	abs := AbsoluteDirection(source, target)

	distance = abs.R

	// Absolute direction in cartesian
	x := distance * math.Sin(abs.Θ) * math.Cos(abs.Φ)
	y := distance * math.Sin(abs.Θ) * math.Sin(abs.Φ)
	z := distance * math.Cos(abs.Θ)

	sinφ := math.Sin(Rad(source.Φ))
	cosφ := math.Cos(Rad(source.Φ))
	sinλ := math.Sin(Rad(source.Λ))
	cosλ := math.Cos(Rad(source.Λ))

	// Relative direction in cartesian (rotate so forward is North, up is away from earth core)
	// ie. Rotate 270 - Φ about y axis, then Λ - 180 about the new x axis
	xr := x*sinφ - y*cosφ*sinλ - z*cosφ*cosλ
	yr := z*sinλ - y*cosλ
	zr := x*cosφ + y*sinφ*sinλ + z*sinφ*cosλ

	// Because heading is positive clockwise, but spherical is positive anticlockwise
	heading = 180 - deg(math.Atan2(yr, xr)) - facing
	if heading >= 360 {
		heading -= 360
	}
	elevation = 90 - deg(math.Acos(zr/distance))

	return distance, heading, elevation
}

func AbsoluteDirection(sourceLLA LLACoords, targetLLA LLACoords) SphericalCoords {
	source := LLAToECEF(sourceLLA)
	target := LLAToECEF(targetLLA)

	d := ECEFCoords{
		X: target.X - source.X,
		Y: target.Y - source.Y,
		Z: target.Z - source.Z,
	}

	r := math.Sqrt(d.X*d.X + d.Y*d.Y + d.Z*d.Z)

	return SphericalCoords{
		R: r,
		Θ: math.Acos(d.Z / r),
		Φ: math.Atan2(d.Y, d.X),
	}
}
