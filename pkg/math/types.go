package math

import (
	m "math"
)

const Π = m.Pi

type Degrees float64
type Radians float64
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

// LLACoords represents a position on earth by latitude (Φ), longitude (Λ) and altitude above average sea level (Azimuth).
type LLACoords struct {
	Φ Degrees
	Λ Degrees
	A Meters
}

// ECEFCoords represents an Earth Centered, Earth Fixed cartesian location in meters from the center towards lat/long 0, 0 (X), towards lat/long 0,90 (Y) and towards lat/long 90,0 (Z).
type ECEFCoords struct {
	X Meters
	Y Meters
	Z Meters
}
