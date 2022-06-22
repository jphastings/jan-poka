package projections

import (
	. "github.com/jphastings/jan-poka/pkg/math"
)

var Equirectangular = Projection{
	ID: "equirectangular-globe",
	Anchors: []AnchorPoint{
		{Name: "Top left corner", Coords: LLACoords{Latitude: 90, Longitude: -180}},
		{Name: "Top right corner", Coords: LLACoords{Latitude: 90, Longitude: 180}},
		{Name: "Bottom right corner", Coords: LLACoords{Latitude: -90, Longitude: 180}},
		{Name: "Bottom left corner", Coords: LLACoords{Latitude: -90, Longitude: -180}},
	},
	Normalize: func(coords LLACoords) Pos {
		return Pos{
			X: float64(coords.Longitude),
			Y: float64(coords.Latitude),
		}
	},
	Reverse: func(pos Pos) LLACoords {
		return LLACoords{
			Latitude:  Degrees(pos.Y),
			Longitude: Degrees(pos.X),
		}
	},
}
