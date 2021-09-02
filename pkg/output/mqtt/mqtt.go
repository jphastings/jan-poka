package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/future"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Message will be interpreted by microcontrollers. Keep JSON keys at 8 or fewer characters.
type Message struct {
	TargetLatitude  float32 `json:"lat"`
	TargetLongitude float32 `json:"lng"`
	TargetAltitude  float32 `json:"alt"`

	CalculatedAzimuth         float32 `json:"azi"`
	CalculatedElevation       float32 `json:"ele"`
	CalculatedRange           float32 `json:"rng"`
	CalculatedSurfaceDistance float32 `json:"dst"`

	CalculatedMapperLengths map[int]Lengths `json:"map"`
}

type Lengths struct {
	Left  float32 `json:"l"`
	Right float32 `json:"r"`
}

type Config struct {
	client  mqtt.Client
	topic   string
	timeout time.Duration
}

func New(broker, username, password, topic string, timeout time.Duration) (*Config, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", broker))
	opts.SetClientID(clientID())
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.OnConnect = func(mqtt.Client) { log.Println("Connected to MQTT") }
	opts.OnConnectionLost = func(_ mqtt.Client, err error) { log.Printf("MQTT connection lost: %v\n", err) }

	client := mqtt.NewClient(opts)
	s := &Config{
		client:  client,
		topic:   topic,
		timeout: timeout,
	}
	return s, s.tokenOk(client.Connect())
}

// clientID generates a unique ID for MQTT to use (as only one client of each name can exist on a server at once)
func clientID() string {
	uid, err := machineid.ID()
	if err != nil {
		uid = "unknown-machine-id"
	}
	return fmt.Sprintf("jan-poka:publisher:%s", uid)
}

func (c *Config) tokenOk(token mqtt.Token) error {
	if token.WaitTimeout(c.timeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *Config) TrackerCallback(details common.TrackedDetails) future.Future {
	return future.Exec(func() error {
		msg := Message{
			TargetLatitude:            float32(details.Target.Latitude),
			TargetLongitude:           float32(details.Target.Longitude),
			TargetAltitude:            float32(details.Target.Altitude),
			CalculatedAzimuth:         float32(details.Bearing.Azimuth),
			CalculatedElevation:       float32(details.Bearing.Elevation),
			CalculatedRange:           float32(details.Bearing.Range),
			CalculatedSurfaceDistance: float32(details.UnobstructedDistance),
			CalculatedMapperLengths:   mapperLengths(details.MapperLengths),
		}

		enc, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		return c.tokenOk(c.client.Publish(c.topic, 0, false, enc))
	})
}

func mapperLengths(mls map[int]common.WallPos) map[int]Lengths {
	ls := make(map[int]Lengths)
	for i, ml := range mls {
		ls[i] = Lengths{
			Left:  float32(ml.LengthLeft),
			Right: float32(ml.LengthRight),
		}
	}
	return ls
}
