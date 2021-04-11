package env

import (
	"github.com/kelseyhightower/envconfig"
	"time"

	. "github.com/jphastings/jan-poka/pkg/math"
)

type Config struct {
	Facing Degrees `default:"0"`
	Port   uint16  `default:"2678"`

	HomeLatitude  Degrees `required:"true"`
	HomeLongitude Degrees `required:"true"`
	HomeAltitude  Meters  `required:"true"`

	UseLog   bool `default:"true"`
	UseTower bool `default:"false"`
	UseAudio bool `default:"false"`

	TowerStatePath string `default:"/run/jan-poka/"`

	MQTTBroker   string `default:"mqtt.local:1883"`
	MQTTUsername string `default:"jan-poka"`
	MQTTPassword string `default:""`
	MQTTTopic    string `default:"home/geo/target"`

	TCPTimeout time.Duration `default:"1s"`

	InstagramUsername string
	InstagramPassword string

	GoogleMapsAPIKey string

	Home LLACoords `ignored:"true"`
}

func ParseEnv() (Config, error) {
	var env Config
	err := envconfig.Process("jp", &env)
	if err != nil {
		return env, err
	}

	env.Home = LLACoords{
		Latitude:  env.HomeLatitude,
		Longitude: env.HomeLongitude,
		Altitude:  env.HomeAltitude,
	}

	return env, nil
}
