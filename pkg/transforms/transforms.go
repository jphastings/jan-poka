package transforms

import (
	"math"
)

type LLACoords struct {
	Φ float64
	Λ float64
	A float64
}

type SphericalCoords struct {
	R float64
	Θ float64
	Φ float64
}

type ECEFCoords struct {
	X float64
	Y float64
	Z float64
}

type geometry struct {
	a float64
	f float64
}

var WGS84 = geometry{
	a: 6378137,
	f: 1.0 / 298.257223518,
}

func rad(deg float64) float64 {
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

	sinφ := math.Sin(rad(source.Φ))
	cosφ := math.Cos(rad(source.Φ))
	sinλ := math.Sin(rad(source.Λ))
	cosλ := math.Cos(rad(source.Λ))

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

func LLAToECEF(coords LLACoords) ECEFCoords {
	geo := WGS84

	sinφ := math.Sin(rad(coords.Φ))
	cosφ := math.Cos(rad(coords.Φ))
	sinλ := math.Sin(rad(coords.Λ))
	cosλ := math.Cos(rad(coords.Λ))

	a := geo.a
	b := geo.a * (1 - geo.f)

	a2 := a * a
	b2 := b * b

	e := math.Sqrt((a2 - b2) / a2)

	e2 := e * e
	sin2φ := sinφ * sinφ

	n := a / math.Sqrt(1-e2*sin2φ)
	h := coords.A

	return ECEFCoords{
		X: (n + h) * cosφ * cosλ,
		Y: (n + h) * cosφ * sinλ,
		Z: ((b2/a2)*n + h) * sinφ,
	}
}
