package mapper

import (
	. "github.com/jphastings/jan-poka/pkg/math"
	"math"
)

func (c *Config) Calculate(coords LLACoords) []*WallPos {
	pos := make([]*WallPos, len(c.Mappers))
	for i, mc := range c.Mappers {
		pos[i] = calcWallPosition(mc, coords)
	}

	return pos
}

func calcWallPosition(s State, coords LLACoords) *WallPos {
	for _, m := range s.Maps {
		if !coordsWithinBounds(m.BottomLeft, m.TopRight, coords) {
			continue
		}

		pos, err := projectCoords(m, coords, s.WallConfig)
		if err != nil {
			panic(err) // TODO: don't panic
			continue
		}
		return &pos
	}
	return nil
}

func coordsWithinBounds(BottomRight, TopLeft Correlation, coords LLACoords) bool {
	// Longitude is positive towards the right
	// Latitude is positive towards the top
	return coords.Longitude <= BottomRight.Longitude &&
		coords.Longitude >= TopLeft.Longitude &&
		coords.Latitude <= TopLeft.Latitude &&
		coords.Latitude >= BottomRight.Latitude
}

func projectCoords(ms MapSpec, coords LLACoords, w WallConfig) (WallPos, error) {
	x, y, err := ms.ToCartesian(coords)
	if err != nil {
		return WallPos{}, err
	}

	transforms, err := ms.Transforms(w)
	if err != nil {
		return WallPos{}, err
	}

	X := transforms.Scale*x + transforms.Tx
	Y := -transforms.Scale*y + transforms.Ty

	X2 := math.Pow(X, 2)
	Y2 := math.Pow(Y, 2)
	dxmX2 := math.Pow(float64(w.Width)-X, 2)

	wheelRadiusSquared := math.Pow(float64(w.WheelRadius), 2)

	return WallPos{
		LengthLeft:  Meters(math.Sqrt(X2 + Y2 - wheelRadiusSquared)),
		LengthRight: Meters(math.Sqrt(dxmX2 + Y2 - wheelRadiusSquared)),
	}, nil
}
