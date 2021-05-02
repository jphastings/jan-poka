package math

import (
	"github.com/pebbe/proj/v5"
)

var projCtx = proj.NewContext()

type Projection struct {
	proj *proj.PJ
}

func NewProjection(projStr string) (*Projection, error) {
	pj, err := projCtx.Create(projStr)
	if err != nil {
		return nil, err
	}

	return &Projection{proj: pj}, nil
}

func (lla LLACoords) Map(projection *Projection) (MapCoords, error) {
	x, y, _, _, err := projection.proj.Trans(proj.Fwd,
		float64(lla.Latitude), float64(lla.Longitude), 0, 0)
	if err != nil {
		return MapCoords{}, err
	}

	return MapCoords{Horizontal: x, Vertical: y}, nil
}
