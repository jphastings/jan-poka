package env

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	. "github.com/jphastings/jan-poka/pkg/math"
)

type Config struct {
	Heading Degrees `default:"0"`
	Port    uint16  `default:"2678"`

	HomeLatitude  Degrees `required:"true"`
	HomeLongitude Degrees `required:"true"`
	HomeAltitude  Meters  `required:"true"`

	UseLog      bool `default:"true"`
	UseSteppers bool `default:"false"`
	UseTTS      bool `default:"false"`

	MotorSteps           int           `default:"200"`
	SphereDiameter       Meters        `default:"0.2"`
	OmniwheelDiameter    Meters        `default:"0.048"`
	MinStepInterval      time.Duration `default:"400us"`
	MotorAutoSleepLeeway time.Duration `default:"500ms"`

	Home LLACoords `ignored:"true"`
}

func ParseEnv() (Config, error) {
	var env Config
	err := envconfig.Process("jan-poka", &env)
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
