package env

import (
	_ "embed"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	janpoka "github.com/jphastings/jan-poka"
	. "github.com/jphastings/jan-poka/pkg/math"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Facing Degrees `default:"0"`
	Port   uint16  `default:"2678"`

	HomeLatitude  Degrees `required:"true"`
	HomeLongitude Degrees `required:"true"`
	HomeAltitude  Meters  `default:"0"`

	UseLog    bool `default:"true"`
	UseAudio  bool `default:"false"`
	UseMapper bool `default:"false"`

	IPAddress string `default:""`
	MQTTPort  int    `default:"1883"`

	Persistence string `default:"~/.jan-poka"`

	TCPTimeout time.Duration `default:"1s"`

	InstagramUsername string
	InstagramPassword string
	GoogleMapsAPIKey  string
	Dump1090Host      string

	Home LLACoords `ignored:"true"`

	Version string `ignored:"true"`
}

func ParseEnv() (Config, error) {
	var env Config
	err := envconfig.Process("jp", &env)
	if err != nil {
		return env, err
	}

	env.Version = janpoka.Version

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
		if s, sErr := os.Stat(env.Persistence); os.IsNotExist(sErr) || !s.IsDir() {
			if mkErr := os.Mkdir(env.Persistence, 0755); mkErr != nil {
				return env, mkErr
			}
		}
	}

	return env, nil
}
