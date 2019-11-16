package l10n

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/math"
)

func TrackerCallback(name string, bearing math.AERCoords, distance math.Meters, _ bool) future.Future {
	f := future.New()
	go func() {
		fmt.Println(Phrase(name, bearing, distance, false))
		f.Succeed()
	}()
	return f
}
