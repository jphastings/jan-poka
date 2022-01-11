package mqtt

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/mdns"
	"github.com/jphastings/jan-poka/pkg/shutdown"
	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
	"github.com/mochi-co/mqtt/server/persistence/bolt"
	"go.etcd.io/bbolt"
)

const Topic = "home/geo/target"

// Message will be interpreted by microcontrollers. Keep JSON keys at 8 or fewer characters.
type Message struct {
	TargetLatitude  float32 `json:"lat"`
	TargetLongitude float32 `json:"lng"`
	TargetAltitude  float32 `json:"alt"`

	CalculatedAzimuth         float32 `json:"azi"`
	CalculatedElevation       float32 `json:"ele"`
	CalculatedRange           float32 `json:"rng"`
	CalculatedSurfaceDistance float32 `json:"dst"`

	LocalTime       string `json:"time"`
	UTCOffsetInMins int16  `json:"tutc"`
	DSTOffsetInMins int16  `json:"tdst"`

	CalculatedMapperLeft    float32         `json:"r1"`
	CalculatedMapperRight   float32         `json:"r2"`
	CalculatedMapperLengths map[int]Lengths `json:"map"`

	Reset bool `json:"reset"`
}

type Lengths struct {
	Left  float32 `json:"l"`
	Right float32 `json:"r"`
}

type Config struct {
	server *mqtt.Server
}

func New(port int, persistence string) (*Config, error) {
	if port == 0 {
		// This is because mqtt doesn't expose the Port in its internal tcp listener
		return nil, fmt.Errorf("MQTT subsystem cannot select a random free port")
	}

	server := mqtt.New()
	tcp := listeners.NewTCP("tcp", fmt.Sprintf(":%d", port))
	err := server.AddListener(tcp, &listeners.Config{Auth: ReadOnlyAuth})
	if err != nil {
		return nil, err
	}

	if persistence != "" {
		db := bolt.New(filepath.Join(persistence, "mqtt.db"), &bbolt.Options{
			Timeout: 500 * time.Millisecond,
		})
		if err := server.AddStore(db); err != nil {
			return nil, err
		}
	}

	shutdown.Ensure("MQTT server", server.Close)
	if err := server.Serve(); err != nil {
		return nil, err
	}
	if _, err := mdns.Register("MQTT", "_jan_poka_mqtt._tcp", port); err != nil {
		return nil, err
	}

	return &Config{server: server}, nil
}

func (c *Config) Close() {
	c.server.Close()
}

func (c *Config) TrackerCallback(details common.TrackedDetails) future.Future {
	_, offsetSecs := details.LocalTime.Zone()
	ml := mapperLengths(details.MapperLengths)

	return c.Publish(Message{
		TargetLatitude:            float32(details.Target.Latitude),
		TargetLongitude:           float32(details.Target.Longitude),
		TargetAltitude:            float32(details.Target.Altitude),
		CalculatedAzimuth:         float32(details.Bearing.Azimuth),
		CalculatedElevation:       float32(details.Bearing.Elevation),
		CalculatedRange:           float32(details.Bearing.Range),
		CalculatedSurfaceDistance: float32(details.UnobstructedDistance),

		LocalTime:       details.LocalTime.Format("15:04:05"),
		UTCOffsetInMins: int16(offsetSecs / 60),
		DSTOffsetInMins: 0, // TODO

		CalculatedMapperLeft:    ml[0].Left,
		CalculatedMapperRight:   ml[0].Right,
		CalculatedMapperLengths: ml,
	})
}

func (c *Config) Publish(msg Message) future.Future {
	return future.Exec(func() error {
		enc, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		return c.server.Publish(Topic, enc, true)
	})
}

func mapperLengths(mls map[int]common.WallPos) map[int]Lengths {
	ls := make(map[int]Lengths)
	for i, ml := range mls {
		ls[i] = Lengths{
			Left:  float32(ml.Left),
			Right: float32(ml.Right),
		}
	}
	return ls
}
