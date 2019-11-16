package celestial

/*
#cgo LDFLAGS: -lnova
#include <libnova/julian_day.h>
#include <libnova/sidereal_time.h>
#include <libnova/lunar.h>
#include <libnova/solar.h>
#include <libnova/mercury.h>
#include <libnova/venus.h>
#include <libnova/mars.h>
#include <libnova/jupiter.h>
#include <libnova/saturn.h>
#include <libnova/neptune.h>
#include <libnova/uranus.h>
#include <libnova/pluto.h>
*/
import "C"
import (
	"time"

	. "github.com/jphastings/jan-poka/pkg/math"
)

type Julian C.double

func JulianDay(t time.Time) Julian {
	t1 := t.UTC()
	lnd := C.struct_ln_date{
		years:   C.int(t1.Year()),
		months:  C.int(t1.Month()),
		days:    C.int(t1.Day()),
		hours:   C.int(t1.Hour()),
		minutes: C.int(t1.Minute()),
		seconds: C.double(t1.Second()) + C.double(t1.Nanosecond())*1e-9,
	}

	return Julian(C.ln_get_julian_day(&lnd))
}

func GeocentricCoordinates(cb Body, j Julian) LLACoords {
	cj := C.double(j)
	var equ C.struct_ln_equ_posn
	var dist Meters

	switch cb {
	case Moon:
		C.ln_get_lunar_equ_coords(cj, &equ)
		dist = Meters(C.ln_get_lunar_earth_dist(cj) * 1000)
	case Sun:
		C.ln_get_solar_equ_coords(cj, &equ)
		dist = AstronomicalUnits(1).Meters()
	case Mercury:
		C.ln_get_mercury_equ_coords(cj, &equ)
		dist = AstronomicalUnits(C.ln_get_mercury_earth_dist(cj)).Meters()
	case Venus:
		C.ln_get_venus_equ_coords(cj, &equ)
		dist = AstronomicalUnits(C.ln_get_venus_earth_dist(cj)).Meters()
	case Mars:
		C.ln_get_mars_equ_coords(cj, &equ)
		dist = AstronomicalUnits(C.ln_get_mars_earth_dist(cj)).Meters()
	case Jupiter:
		C.ln_get_jupiter_equ_coords(cj, &equ)
		dist = AstronomicalUnits(C.ln_get_jupiter_earth_dist(cj)).Meters()
	case Saturn:
		C.ln_get_saturn_equ_coords(cj, &equ)
		dist = AstronomicalUnits(C.ln_get_saturn_earth_dist(cj)).Meters()
	case Neptune:
		C.ln_get_neptune_equ_coords(cj, &equ)
		dist = AstronomicalUnits(C.ln_get_neptune_earth_dist(cj)).Meters()
	case Uranus:
		C.ln_get_uranus_equ_coords(cj, &equ)
		dist = AstronomicalUnits(C.ln_get_uranus_earth_dist(cj)).Meters()
	case Pluto:
		C.ln_get_pluto_equ_coords(cj, &equ)
		dist = AstronomicalUnits(C.ln_get_pluto_earth_dist(cj)).Meters()
	default:
		return LLACoords{}
	}

	return LLACoords{
		Latitude:  Degrees(equ.dec),
		Longitude: Degrees(equ.ra - C.ln_get_apparent_sidereal_time(cj)*15),
		Altitude:  dist,
	}
}
