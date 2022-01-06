package temporal

import (
	"sort"
	"time"

	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/twilight"
	"github.com/zsefvlol/timezonemapper"
)

func LocalTimeAndSkiesAt(coords math.LLACoords) (time.Time, []common.SkyChange, error) {
	now := time.Now()

	tz := timezonemapper.LatLngToTimezoneString(float64(coords.Latitude), float64(coords.Longitude))
	location, err := time.LoadLocation(tz)
	if err != nil {
		return now, nil, err
	}

	localNow := now.In(location)
	return localNow, skyChanges(localNow, coords), nil
}

func skyChanges(localTime time.Time, coords math.LLACoords) []common.SkyChange {
	today := skyChangesForDay(localTime, coords)
	dayAhead := localTime.Add(24 * time.Hour)
	tomorrow := skyChangesForDay(dayAhead, coords)

	var skyChanges []common.SkyChange
	for _, change := range today {
		if !change.Time.Before(localTime) {
			skyChanges = append(skyChanges, change)
		}
	}
	for _, change := range tomorrow {
		if !change.Time.After(dayAhead) {
			skyChanges = append(skyChanges, change)
		}
	}

	sort.Slice(skyChanges, func(i, j int) bool {
		return skyChanges[i].Time.Before(skyChanges[j].Time)
	})

	skyNow := common.SkyChange{Sky: skyChanges[len(skyChanges)-1].Sky, Time: localTime}
	skyChanges = append([]common.SkyChange{skyNow}, skyChanges...)

	return skyChanges
}

func skyChangesForDay(localTime time.Time, coords math.LLACoords) []common.SkyChange {
	lat, lng := float64(coords.Latitude), float64(coords.Longitude)

	sunrise, sunset, status := twilight.SunRiseSet(localTime, lat, lng)
	switch status {
	case twilight.SunriseStatusAboveHorizon:
		// Day all the time
		return []common.SkyChange{{Sky: common.SkyDay, Time: localTime}}
	case twilight.SunriseStatusBelowHorizon:
		// Night all the time
		return []common.SkyChange{{Sky: common.SkyNight, Time: localTime}}
	}

	// TODO: Better names
	civilRise, civilSet, _ := twilight.CivilTwilight(localTime, lat, lng)
	astroRise, astroSet, _ := twilight.AstronomicalTwilight(localTime, lat, lng)

	return []common.SkyChange{
		{Sky: common.SkyDay, Time: sunrise},
		{Sky: common.SkyCivil, Time: sunset},
		{Sky: common.SkyAstro, Time: civilSet},
		{Sky: common.SkyNight, Time: astroSet},
		{Sky: common.SkyAstro, Time: astroRise},
		{Sky: common.SkyCivil, Time: civilRise},
	}
}
