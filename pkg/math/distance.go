package math

type DistanceScale float64

const (
	ByFoot  DistanceScale = 5040
	ByCar   DistanceScale = 72420
	ByPlane DistanceScale = 885139
)

func Scaled(distance Meters) (float64, DistanceScale) {
	hoursByFoot := float64(distance) / float64(ByFoot)
	hoursByCar := float64(distance) / float64(ByCar)
	hoursByPlane := float64(distance) / float64(ByPlane)

	if hoursByCar < 1 {
		return hoursByFoot, ByFoot
	} else if hoursByPlane < 1 {
		return hoursByCar, ByCar
	} else {
		return hoursByPlane, ByPlane
	}
}
