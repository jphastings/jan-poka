package mqtt

import (
	"encoding/json"
	"fmt"
	"github.com/jphastings/jan-poka/pkg/future"
	"github.com/jphastings/jan-poka/pkg/math"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// This will be interpreted by microcontrollers. Keep JSON keys at 8 or fewer characters.
type Message struct {
	TargetLatitude  float32 `json:"lat"`
	TargetLongitude float32 `json:"lng"`
	TargetAltitude  float32 `json:"alt"`

	CalculatedAzimuth         float32 `json:"azi"`
	CalculatedElevation       float32 `json:"ele"`
	CalculatedRange           float32 `json:"rng"`
	CalculatedSurfaceDistance float32 `json:"dst"`
}

type Config struct {
	client  mqtt.Client
	topic   string
	timeout time.Duration
}

func New(broker, username, password, topic string, timeout time.Duration) (*Config, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", broker))
	opts.SetClientID("jan-poka:publisher")
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
func (s *Config) tokenOk(token mqtt.Token) error {
	if token.WaitTimeout(s.timeout) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *Config) TrackerCallback(_ string, target math.LLACoords, bearing math.AERCoords, surfaceDistance math.Meters, _ bool) future.Future {
	return future.Exec(func() error {
		msg := Message{
			TargetLatitude:            float32(target.Latitude),
			TargetLongitude:           float32(target.Longitude),
			TargetAltitude:            float32(target.Altitude),
			CalculatedAzimuth:         float32(bearing.Azimuth),
			CalculatedElevation:       float32(bearing.Elevation),
			CalculatedRange:           float32(bearing.Range),
			CalculatedSurfaceDistance: float32(surfaceDistance),
		}
		enc, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		return s.tokenOk(s.client.Publish(s.topic, 0, false, enc))
	})
}
