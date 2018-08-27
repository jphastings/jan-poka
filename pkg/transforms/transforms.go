package transforms

func RelativeDirection(source LLACoords, target LLACoords, facing degrees) (distance meters, heading, elevation degrees) {
	abs := AbsoluteDirection(source, target)

	distance = abs.R

	// Absolute direction in cartesian
	x := distance * sin(abs.Θ) * cos(abs.Φ)
	y := distance * sin(abs.Θ) * sin(abs.Φ)
	z := distance * cos(abs.Θ)

	sinφ := sinº(source.Φ)
	cosφ := cosº(source.Φ)
	sinλ := sinº(source.Λ)
	cosλ := cosº(source.Λ)

	// Relative direction in cartesian (rotate so forward is North, up is away from earth core)
	// ie. Rotate 270 - Φ about y axis, then Λ - 180 about the new x axis
	xr := x*sinφ - y*cosφ*sinλ - z*cosφ*cosλ
	yr := z*sinλ - y*cosλ
	zr := x*cosφ + y*sinφ*sinλ + z*sinφ*cosλ

	// Because heading is positive clockwise, but spherical is positive anticlockwise
	heading = 180 - atan2D(yr, xr) - facing
	if heading >= 360 {
		heading -= 360
	}
	elevation = 90 - acosD(zr/distance)

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

	r := mod(d)

	return SphericalCoords{
		R: r,
		Θ: acos(d.Z / r),
		Φ: atan2(d.Y, d.X),
	}
}
