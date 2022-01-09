package mapper

import (
	"encoding/json"
	"fmt"
	. "github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/wroge/wgs84"
	"io/ioutil"
	"math"
	"path/filepath"
	"sync"

	. "github.com/jphastings/jan-poka/pkg/math"
)

type Config struct {
	Path string
	mu   sync.Mutex

	Mappers []State
}

type State struct {
	WallPos    `json:"currentPosition"`
	WallConfig WallConfig `json:"wallConfig"`
	// MapSpecs earlier in this slice will take precedence
	Maps []MapSpec `json:"maps"`
}

type WallConfig struct {
	Width       Meters
	WheelRadius Meters
}

func New(configRoot string) (*Config, error) {
	configPath := filepath.Join(configRoot, "mapper.json")
	s := &Config{Path: configPath}

	data, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &s.Mappers); err != nil {
		return nil, err
	}

	availableProjections := wgs84.EPSG()

	for i, mc := range s.Mappers {
		for j, m := range mc.Maps {
			projection := availableProjections.Code(m.EPSGCode)
			// Bizarrely, this library just returns geocentric if the code isn't valid
			if projection == wgs84.XYZ() {
				return nil, fmt.Errorf("projection code EPSG:%d not supported", m.EPSGCode)
			}

			s.Mappers[i].Maps[j].transform = wgs84.LonLat().To(projection)
		}
	}

	return s, nil
}

type MapSpec struct {
	Name       string
	TopRight   Correlation
	BottomLeft Correlation

	EPSGCode  int `json:"epsgCode"`
	transform wgs84.Func

	// memoization of the derived values
	transforms *Transforms
}

type Correlation struct {
	WallPos   `json:"wallPosition"`
	LLACoords `json:"latLong"`
}

type Transforms struct {
	Scale float64
	Tx    float64
	Ty    float64
}

func (ms MapSpec) ToCartesian(coords LLACoords) (float64, float64, error) {
	x, y, _ := ms.transform(float64(coords.Latitude), float64(coords.Longitude), float64(coords.Altitude))
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

	ms.transforms = &Transforms{
		Tx: (xa*Xb - xb*Xa) / (xa - xb),
		Ty: (yb*Ya - ya*Yb) / (yb - ya),
	}
	ms.transforms.Scale = (Xa - ms.transforms.Tx) / xa

	return ms.transforms, nil
}

func calcXY(wall WallPos, width, wheelRadius Meters) (float64, float64) {
	left2 := math.Pow(float64(wall.Left), 2)
	right2 := math.Pow(float64(wall.Right), 2)
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

func (c *Config) TrackerCallback(details TrackedDetails) future.Future {
	return future.Exec(func() error {
		details.MapperLengths = c.Calculate(details.Target)
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
