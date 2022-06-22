package projections

import (
	. "github.com/jphastings/jan-poka/pkg/math"
	"math"
)

var Winkel = Projection{
	Anchors: []AnchorPoint{
		{Name: "Top top, middle", Coords: LLACoords{Latitude: 90, Longitude: 0}},
		{Name: "Rightmost", Coords: LLACoords{Latitude: 0, Longitude: 180}},
		{Name: "Bottom middle", Coords: LLACoords{Latitude: -90, Longitude: 0}},
		{Name: "Leftmost", Coords: LLACoords{Latitude: 0, Longitude: -180}},
	},
	Normalize: func(coords LLACoords) Pos {
		lam := float64(coords.Longitude * (Pi / 180))
		phi := float64(coords.Latitude * (Pi / 180))
		alf := math.Acos(math.Cos(phi) * math.Cos(lam/2))

		sincAlf := float64(1)
		if alf != 0 {
			sincAlf = math.Sin(alf) / alf
		}

		// NB. all constant scaling factors removed (per axis), cos will be transformed anyway
		return Pos{
			X: lam*2/Pi + ((2 * math.Cos(phi) * math.Sin(lam/2)) / sincAlf),
			Y: phi + (math.Sin(phi) / sincAlf),
		}
	},
}
