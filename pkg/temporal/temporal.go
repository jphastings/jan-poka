package temporal

import (
	"time"

	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/zsefvlol/timezonemapper"
)

func LocalTime(coords math.LLACoords) (time.Time, error) {
	now := time.Now()

	tz := timezonemapper.LatLngToTimezoneString(float64(coords.Latitude), float64(coords.Longitude))
	location, err := time.LoadLocation(tz)
	if err != nil {
		return now, err
	}

	localNow := now.In(location)
	return localNow, nil
}
