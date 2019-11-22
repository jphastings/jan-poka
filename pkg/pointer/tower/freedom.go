package tower

import "github.com/jphastings/jan-poka/pkg/math"

func Pointer(theta, phi math.Degrees) (base math.Degrees, arm math.Degrees) {
	clockwise := theta < 180

	if clockwise {
		base = theta - 90
	} else {
		base = theta - 270
	}

	if clockwise {
		phi *= -1
	}

	return base, phi
}
