package math

import geo "github.com/kellydunn/golang-geo"

func (source LLACoords) DirectionTo(target LLACoords, facing Degrees) AERCoords {
	sourceECEF := source.ECEF(WGS84)
	targetECEF := target.ECEF(WGS84)

	geocentric := ECEFCoords{
		X: targetECEF.X - sourceECEF.X,
		Y: targetECEF.Y - sourceECEF.Y,
		Z: targetECEF.Z - sourceECEF.Z,
	}

	topocentric := geocentric.ENU(source)
	bearing := topocentric.AER(facing)

	return bearing
}

func (source LLACoords) GreatCircleDistance(target LLACoords) Meters {
	src := geo.NewPoint(float64(source.Latitude), float64(source.Longitude))
	trg := geo.NewPoint(float64(target.Latitude), float64(target.Longitude))

	// find the great circle distance between them
	flatDist := Meters(src.GreatCircleDistance(trg) * 1000)

	heightDiff := target.Altitude - source.Altitude
	dist := Sqrt(flatDist*flatDist + heightDiff*heightDiff)

	return dist
}
