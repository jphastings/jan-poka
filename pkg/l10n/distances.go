package l10n

import (
	"fmt"
	"math"

	. "github.com/jphastings/jan-poka/pkg/math"
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
	{"hours' walk away", 5040},     // meters per hour
	{"hours' drive away", 72420},   // meters per hour
	{"hours' flight away", 885139}, // meters per hour
}

type numberRepresentation struct {
	amount   float64
	accuracy float64
}

func Distance(m Meters) string {
	var oldRep numberRepresentation
	var r rangedDistance

	for i, conversion := range distanceUnits {
		newRep := roundSigFigs(float64(m) / conversion.amount)

		if i == 0 || isBetter(newRep, oldRep) {
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

	if r.amount == 0 && frac != "" {
		return fmt.Sprintf("%s %s", frac, r.unit)
	} else {
		return fmt.Sprintf("%.0f%s %s", r.amount, frac, r.unit)
	}
}

func isBetter(newRep numberRepresentation, oldRep numberRepresentation) bool {
	return newRep.amount >= 1
}

// Uses a heuristic for SI units for jumps
//func isBetter(newRep numberRepresentation, oldRep numberRepresentation) bool {
//	return newRep.amount >= 0.5 && (oldRep.amount > jumpFactor || newRep.accuracy >= oldRep.accuracy)
//}

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
