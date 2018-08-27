package transforms

import "math"

type geometry struct {
	a meters
	f meters
}

var WGS84 = geometry{
	a: 6378137,
	f: 1.0 / 298.257223518,
}

func LLAToECEF(coords LLACoords) ECEFCoords {
	sinφ := sinº(coords.Φ)
	cosφ := cosº(coords.Φ)
	sinλ := sinº(coords.Λ)
	cosλ := cosº(coords.Λ)

	a := WGS84.a
	b := WGS84.a * (1 - WGS84.f)

	a2 := a * a
	b2 := b * b

	e := meters(math.Sqrt(float64((a2 - b2) / a2)))

	e2 := e * e
	sin2φ := sinφ * sinφ

	n := a / meters(math.Sqrt(float64(1-e2*sin2φ)))
	h := coords.A

	return ECEFCoords{
		X: (n + h) * cosφ * cosλ,
		Y: (n + h) * cosφ * sinλ,
		Z: ((b2/a2)*n + h) * sinφ,
	}
}
