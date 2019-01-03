package l10n

import (
	"fmt"
	"github.com/jphastings/corviator/pkg/future"
	"github.com/jphastings/corviator/pkg/math"
)

func TrackerCallback(name string, bearing math.AERCoords, _ bool) future.Future {
	f := future.New()
	go func() {
		fmt.Println(Phrase(name, bearing, false))
		f.Succeed()
	}()
	return f
}
