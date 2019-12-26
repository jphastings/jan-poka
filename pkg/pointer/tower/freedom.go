package tower

import "github.com/jphastings/jan-poka/pkg/math"

func Pointer(theta, phi math.Degrees) (math.Degrees, math.Degrees) {
	clockwise := theta < 180

	var base math.Degrees
	if clockwise {
		base = theta - 90
	} else {
		base = theta - 270
	}

	arm := 180 - phi
	if clockwise {
		arm *= -1
	}

	return base, arm
}
