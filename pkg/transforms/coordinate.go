package transforms

import "math"

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

func LLAToECEF(coords LLACoords) ECEFCoords {
	sinφ := math.Sin(Rad(coords.Φ))
	cosφ := math.Cos(Rad(coords.Φ))
	sinλ := math.Sin(Rad(coords.Λ))
	cosλ := math.Cos(Rad(coords.Λ))

	a := WGS84.a
	b := WGS84.a * (1 - WGS84.f)

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
