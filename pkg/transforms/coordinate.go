package transforms

import (
	. "github.com/jphastings/corviator/pkg/math"
)

type geometry struct {
	a Meters
	f Meters
}

var WGS84 = geometry{
	a: 6378137,
	f: 1.0 / 298.257223518,
}

func LLAToECEF(coords LLACoords) ECEFCoords {
	sinφ := Sinº(coords.Φ)
	cosφ := Cosº(coords.Φ)
	sinλ := Sinº(coords.Λ)
	cosλ := Cosº(coords.Λ)

	a := WGS84.a
	b := WGS84.a * (1 - WGS84.f)

	a2 := a * a
	b2 := b * b

	e := Sqrt((a2 - b2) / a2)

	e2 := e * e
	sin2φ := sinφ * sinφ

	n := a / Sqrt(1-e2*sin2φ)
	h := coords.A

	return ECEFCoords{
		X: (n + h) * cosφ * cosλ,
		Y: (n + h) * cosφ * sinλ,
		Z: ((b2/a2)*n + h) * sinφ,
	}
}
