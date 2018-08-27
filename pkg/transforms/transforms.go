package transforms

import . "github.com/jphastings/corviator/pkg/math"

func RelativeDirection(source LLACoords, target LLACoords, facing Degrees) Direction {
	abs := AbsoluteDirection(source, target)

	direction := Direction{
		Distance: abs.R,
	}

	// Absolute direction in cartesian
	x := direction.Distance * Sin(abs.Θ) * Cos(abs.Φ)
	y := direction.Distance * Sin(abs.Θ) * Sin(abs.Φ)
	z := direction.Distance * Cos(abs.Θ)

	sinφ := Sinº(source.Φ)
	cosφ := Cosº(source.Φ)
	sinλ := Sinº(source.Λ)
	cosλ := Cosº(source.Λ)

	// Relative direction in cartesian (rotate so forward is North, up is away from earth core)
	// ie. Rotate 270 - Φ about y axis, then Λ - 180 about the new x axis
	xr := x*sinφ - y*cosφ*sinλ - z*cosφ*cosλ
	yr := z*sinλ - y*cosλ
	zr := x*cosφ + y*sinφ*sinλ + z*sinφ*cosλ

	// Not sure why this 180 needs to be there… I thought I'd need to invert as spherical is anticlockwise but
	// headings are clockwise. Tests pass though!
	direction.Heading = ModDeg(180 + Atan2º(yr, xr) - facing)
	direction.Elevation = 90 - Acosº(zr/direction.Distance)

	return direction
}

func AbsoluteDirection(sourceLLA LLACoords, targetLLA LLACoords) SphericalCoords {
	source := LLAToECEF(sourceLLA)
	target := LLAToECEF(targetLLA)

	d := ECEFCoords{
		X: target.X - source.X,
		Y: target.Y - source.Y,
		Z: target.Z - source.Z,
	}

	r := Mod(d)

	return SphericalCoords{
		R: r,
		Θ: Acos(d.Z / r),
		Φ: ModRad(Atan2(d.Y, d.X)),
	}
}
