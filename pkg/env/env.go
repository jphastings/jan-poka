package env

import (
	"os/user"
	"path/filepath"
	"strings"
	"time"

	. "github.com/jphastings/jan-poka/pkg/math"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Facing Degrees `default:"0"`
	Port   uint16  `default:"2678"`

	HomeLatitude  Degrees `required:"true"`
	HomeLongitude Degrees `required:"true"`
	HomeAltitude  Meters  `required:"true"`

	UseLog    bool `default:"true"`
	UseAudio  bool `default:"false"`
	UseMapper bool `default:"false"`

	MQTTBroker   string `default:"mqtt.local:1883"`
	MQTTUsername string `default:"jan-poka"`
	MQTTPassword string `default:""`
	MQTTTopic    string `default:"home/geo/target"`

	Persistence string `default:"~/.jan-poka"`

	TCPTimeout time.Duration `default:"1s"`

	InstagramUsername string
	InstagramPassword string
	GoogleMapsAPIKey  string
	Dump1090Host      string

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

	if strings.HasPrefix(env.Persistence, "~/") {
		usr, err := user.Current()
		if err != nil {
			return env, err
		}
		env.Persistence = filepath.Join(usr.HomeDir, env.Persistence[2:])
	}

	return env, nil
}
