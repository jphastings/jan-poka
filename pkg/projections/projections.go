package projections

import "github.com/jphastings/jan-poka/pkg/math"

type Projection struct {
	Anchors []AnchorPoint
	// Normalize uses the projection to turn a Lat/Long into an arbitrarily scaled cartesian pair
	Normalize func(math.LLACoords) Pos
	// Reverse may be nil if there is no trivial reversibility of this function
	Reverse func(Pos) math.LLACoords
}

type AnchorPoint struct {
	Name   string
	Coords math.LLACoords
}
