package mapper

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/common"
	. "github.com/jphastings/jan-poka/pkg/math"
	"math"
)

func (c *Config) Calculate(coords LLACoords) map[int]common.WallPos {
	pos := make(map[int]common.WallPos)
	for i, mc := range c.Mappers {
		if wp, err := calcWallPosition(mc, coords); err != nil {
			// TODO: Use logger
			fmt.Printf("Error while mapping: %v\n", err)
		} else {
			pos[i] = wp
		}
	}

	return pos
}

func calcWallPosition(s State, coords LLACoords) (common.WallPos, error) {
	for _, m := range s.Maps {
		if !coordsWithinBounds(m.BottomLeft, m.TopRight, coords) {
			continue
		}

		pos, err := projectCoords(m, coords, s.WallConfig)
		if err != nil {
			return common.WallPos{}, err
		}
		return pos, nil
	}
	return common.WallPos{}, fmt.Errorf("map specs don't cover the specified lat/long")
}

func coordsWithinBounds(BottomLeft, TopRight Correlation, coords LLACoords) bool {
	// Longitude is positive towards the right
	// Latitude is positive towards the top
	return coords.Longitude <= TopRight.Longitude &&
		coords.Longitude >= BottomLeft.Longitude &&
		coords.Latitude <= TopRight.Latitude &&
		coords.Latitude >= BottomLeft.Latitude
}

func projectCoords(ms MapSpec, coords LLACoords, w WallConfig) (common.WallPos, error) {
	x, y, err := ms.ToCartesian(coords)
	if err != nil {
		return common.WallPos{}, err
	}

	transforms, err := ms.Transforms(w)
	if err != nil {
		return common.WallPos{}, err
	}

	X := transforms.Scale*x + transforms.Tx
	Y := -transforms.Scale*y + transforms.Ty

	X2 := math.Pow(X, 2)
	Y2 := math.Pow(Y, 2)
	dxmX2 := math.Pow(float64(w.Width)-X, 2)

	wheelRadiusSquared := math.Pow(float64(w.WheelRadius), 2)

	return common.WallPos{
		LengthLeft:  Meters(math.Sqrt(X2 + Y2 - wheelRadiusSquared)),
		LengthRight: Meters(math.Sqrt(dxmX2 + Y2 - wheelRadiusSquared)),
	}, nil
}
