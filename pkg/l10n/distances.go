package l10n

import (
	"fmt"
	. "github.com/jphastings/corviator/pkg/math"
	"math"
)

type rangedDistance struct {
	unit   string
	amount float64
}

const (
	sigfigs    = 3
	jumpFactor = 10000
)

var distanceUnits = []rangedDistance{
	{"km", 1000},
	{"AU", AstronomicalUnitInMeters},
}

type numberRepresentation struct {
	amount   float64
	accuracy float64
}

func Distance(m Meters) string {
	oldRep := roundSigFigs(float64(m))
	r := rangedDistance{"m", oldRep.amount}

	for _, conversion := range distanceUnits {
		newRep := roundSigFigs(float64(m) / conversion.amount)

		if isBetter(newRep, oldRep) {
			r.unit = conversion.unit
			r.amount = newRep.amount
			oldRep.accuracy = newRep.accuracy
		}
	}

	frac := ""
	rem := r.amount - math.Floor(r.amount)
	if rem > 0 {
		r.amount = math.Floor(r.amount)
		switch math.Round(rem * 4) {
		case 0:
			break
		case 1:
			frac = "¼"
		case 2:
			frac = "½"
		case 3:
			if math.Round(rem*2) == 1 {
				frac = "½"
			} else {
				r.amount += 1
			}
		default:
			r.amount += 1
		}
	}

	return fmt.Sprintf("%.0f%s%s", r.amount, frac, r.unit)
}

func isBetter(newRep numberRepresentation, oldRep numberRepresentation) bool {
	return newRep.amount >= 0.5 && (oldRep.amount > jumpFactor || newRep.accuracy >= oldRep.accuracy)
}

func roundSigFigs(num float64) numberRepresentation {
	oldNum := num

	scale := int(math.Ceil(math.Log10(num)))
	if scale < 0 {
		// Less than a meter distances are irrelevant
		return numberRepresentation{}
	}
	scaling := math.Pow10(scale - sigfigs)
	num = math.Round(num/scaling) * scaling

	return numberRepresentation{amount: num, accuracy: num / oldNum}
}
