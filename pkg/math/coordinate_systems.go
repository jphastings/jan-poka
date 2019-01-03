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
	sinTheta := SinDeg(lla.Latitude)
	cosTheta := CosDeg(lla.Latitude)
	sinLambda := SinDeg(lla.Longitude)
	cosLambda := CosDeg(lla.Longitude)

	a := geom.a
	b := geom.a * (1 - geom.f)

	a2 := a * a
	b2 := b * b

	e := Sqrt((a2 - b2) / a2)

	e2 := e * e
	sin2Theta := sinTheta * sinTheta

	n := a / Sqrt(1-e2*sin2Theta)
	h := lla.Altitude

	return ECEFCoords{
		X: (n + h) * cosTheta * cosLambda,
		Y: (n + h) * cosTheta * sinLambda,
		Z: ((b2/a2)*n + h) * sinTheta,
	}
}

// ENU transforms into topocentric meters North, East and Up at the given location.
func (ecef ECEFCoords) ENU(here LLACoords) ENUCoords {
	sinTheta := SinDeg(here.Latitude)
	cosTheta := CosDeg(here.Latitude)
	sinLambda := SinDeg(here.Longitude)
	cosLambda := CosDeg(here.Longitude)

	return ENUCoords{
		East:  -ecef.X*sinLambda + ecef.Y*cosLambda,
		North: -ecef.X*cosLambda*sinTheta - ecef.Y*sinLambda*sinTheta + ecef.Z*cosTheta,
		Up:    ecef.X*cosLambda*cosTheta + ecef.Y*sinLambda*cosTheta + ecef.Z*sinTheta,
	}
}

// AER transforms into a bearing; in degrees Elevated from the horizon, in degrees of Azimuth measured from the given heading (measured from North) as well as the Range.
func (enu ENUCoords) AER(facing Degrees) AERCoords {
	distance := Sqrt(enu.North*enu.North + enu.East*enu.East + enu.Up*enu.Up)
	return AERCoords{
		Azimuth:   ModDeg(90 - Atan2Deg(enu.North, enu.East) - facing),
		Elevation: 90 - AcosDeg(enu.Up/distance),
		Range:     distance,
	}
}
