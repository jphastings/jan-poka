package math

type geometry struct {
	a Meters
	f Meters
}

// WGS84 is the earth geometry used by GPS.
var WGS84 = geometry{
	a: 6378137,
	f: 1.0 / 298.257223518,
}

// ECEF transforms into earth centered, earth fixed coordinates using the given earth geometry.
func (lla LLACoords) ECEF(geom geometry) ECEFCoords {
	sinφ := Sinº(lla.Φ)
	cosφ := Cosº(lla.Φ)
	sinλ := Sinº(lla.Λ)
	cosλ := Cosº(lla.Λ)

	a := geom.a
	b := geom.a * (1 - geom.f)

	a2 := a * a
	b2 := b * b

	e := Sqrt((a2 - b2) / a2)

	e2 := e * e
	sin2φ := sinφ * sinφ

	n := a / Sqrt(1-e2*sin2φ)
	h := lla.A

	return ECEFCoords{
		X: (n + h) * cosφ * cosλ,
		Y: (n + h) * cosφ * sinλ,
		Z: ((b2/a2)*n + h) * sinφ,
	}
}

// ENU transforms into topocentric meters North, East and Up at the given location.
func (ecef ECEFCoords) ENU(here LLACoords) ENUCoords {
	sinφ := Sinº(here.Φ)
	cosφ := Cosº(here.Φ)
	sinλ := Sinº(here.Λ)
	cosλ := Cosº(here.Λ)

	return ENUCoords{
		East:  -ecef.X*sinλ + ecef.Y*cosλ,
		North: -ecef.X*cosλ*sinφ - ecef.Y*sinλ*sinφ + ecef.Z*cosφ,
		Up:    ecef.X*cosλ*cosφ + ecef.Y*sinλ*cosφ + ecef.Z*sinφ,
	}
}

// AER transforms into a bearing; in degrees Elevated from the horizon, in degrees of Azimuth measured from the given heading (measured from North) as well as the Range.
func (enu ENUCoords) AER(facing Degrees) AERCoords {
	distance := Sqrt(enu.North*enu.North + enu.East*enu.East + enu.Up*enu.Up)
	return AERCoords{
		Azimuth:   ModDeg(90 - Atan2º(enu.North, enu.East) - facing),
		Elevation: 90 - Acosº(enu.Up/distance),
		Range:     distance,
	}
}
