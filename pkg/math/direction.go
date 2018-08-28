package math

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
