package l10n

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/math"
)

func TrackerCallback(name string, _ math.LLACoords, bearing math.AERCoords, distance math.Meters, _ bool) future.Future {
	fmt.Println(Phrase(name, bearing, distance, false))
	f := future.New()
	f.Succeed()
	return f
}
