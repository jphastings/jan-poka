package math

const AstronomicalUnitInMeters = 149597870000

type Degrees float64
type Radians float64
type ArcSeconds float64
type AstronomicalUnits float64
type Meters float64

// AERCoords represents an Azimuth (degrees clockwise from North), Elevation (above the horizon) and Range.
type AERCoords struct {
	Azimuth   Degrees
	Elevation Degrees
	Range     Meters
}

// ENUCoords represents a distance East, North and Up.
type ENUCoords struct {
	East  Meters
	North Meters
	Up    Meters
}

// LLACoords represents a position on earth with the (WGS84 referenced) latitude (Latitude), longitude (Longitude) and altitude above average sea level (Altitude).
type LLACoords struct {
	Latitude  Degrees
	Longitude Degrees
	Altitude  Meters
}

// ECEFCoords represents an Earth Centered, Earth Fixed cartesian location in meters from the center of Earth towards lat/long 0, 0 (X), towards lat/long 0,90 (Y) and towards lat/long 90,0 (Z).
type ECEFCoords struct {
	X Meters
	Y Meters
	Z Meters
}

// MapCoords represents a position on a 2d map projection. The horizontal & vertical directions are unspecified.
type MapCoords struct {
	Horizontal float64
	Vertical   float64
}

func (au AstronomicalUnits) Meters() Meters {
	return Meters(au * AstronomicalUnitInMeters)
}
