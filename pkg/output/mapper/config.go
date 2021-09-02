package mapper

import (
	"encoding/json"
	"fmt"
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/pebbe/proj/v5"
	"io/ioutil"
	"math"
	"sync"

	. "github.com/jphastings/jan-poka/pkg/math"
)

type Config struct {
	Path string
	mu   sync.Mutex

	Mappers []State
}

type State struct {
	common.WallPos `json:"currentPosition"`
	WallConfig     WallConfig `json:"wallConfig"`
	// MapSpecs earlier in this slice will take precedence
	Maps []MapSpec `json:"maps"`
}

type WallConfig struct {
	Width       Meters
	WheelRadius Meters
}

type MapSpec struct {
	TopRight   Correlation
	BottomLeft Correlation

	// proj.PJ.Info().Description
	ProjectionDescription string `json:"projection"`
	projection            *proj.PJ

	// memoization of the derived values
	transforms *Transforms
}

type Correlation struct {
	common.WallPos `json:"wallPosition"`
	LLACoords      `json:"latLong"`
}

type Transforms struct {
	Scale float64
	Tx    float64
	Ty    float64
}

func (ms MapSpec) ToCartesian(coords LLACoords) (float64, float64, error) {
	x, y, _, _, err := ms.projection.Trans(proj.Fwd, float64(coords.Latitude), float64(coords.Longitude), 0, 0)
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

func (ms *MapSpec) Transforms(w WallConfig) (*Transforms, error) {
	if ms.transforms != nil {
		return ms.transforms, nil
	}

	xa, ya, err := ms.ToCartesian(ms.BottomLeft.LLACoords)
	if err != nil {
		return nil, err
	}

	xb, yb, err := ms.ToCartesian(ms.TopRight.LLACoords)
	if err != nil {
		return nil, err
	}

	Xa, Ya := calcXY(ms.BottomLeft.WallPos, w.Width, w.WheelRadius)
	Xb, Yb := calcXY(ms.TopRight.WallPos, w.Width, w.WheelRadius)

	ms.transforms.Tx = (xa*Xb - xb*Xa) / (xa - xb)
	ms.transforms.Ty = (yb*Ya - ya*Yb) / (yb - ya)
	ms.transforms.Scale = (Xa - ms.transforms.Tx) / xa

	return ms.transforms, nil
}

func calcXY(wall common.WallPos, width, wheelRadius Meters) (float64, float64) {
	left2 := math.Pow(float64(wall.LengthLeft), 2)
	right2 := math.Pow(float64(wall.LengthRight), 2)
	widthSquared := math.Pow(float64(width), 2)
	wheelRadiusSquared := math.Pow(float64(wheelRadius), 2)

	Y := math.Sqrt(wheelRadiusSquared +
		right2 +
		(left2-widthSquared-right2)/2*float64(width))
	X := math.Sqrt(wheelRadiusSquared +
		left2 -
		math.Pow(Y, 2))

	return X, Y
}

func New(configPath string) (*Config, error) {
	s := &Config{Path: configPath}

	data, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &s.Mappers); err != nil {
		return nil, err
	}

	projCtx := proj.NewContext()
	// Never closed, as the application's purpose is to map

	// Check that all the projections are well-understood
	for i, mc := range s.Mappers {
		for j, m := range mc.Maps {
			p, err := projCtx.Create(m.ProjectionDescription)
			if err != nil {
				return nil, err
			}
			s.Mappers[i].Maps[j].projection = p
		}
	}

	return s, nil
}

func (c *Config) TrackerCallback(details common.TrackedDetails) future.Future {
	return future.Exec(func() error {
		details.MapperLengths = c.Calculate(details.Target)
		// TODO: Remove this fake logging
		fmt.Println(details.MapperLengths)
		return nil
	})
}

func (c *Config) writeConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(c.Mappers)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.Path, data, 0644)
}
