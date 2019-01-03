package env

import (
	"time"

	"github.com/kelseyhightower/envconfig"

	. "github.com/jphastings/corviator/pkg/math"
)

type Config struct {
	Heading   Degrees `default:"0"`
	MicroStep int     `default:"1"`
	Port      uint16  `default:"2678"`

	HomeLatitude  Degrees `required:"true"`
	HomeLongitude Degrees `required:"true"`
	HomeAltitude  Meters  `required:"true"`

	UseLog      bool `default:"true"`
	UseSteppers bool `default:"false"`
	UseTTS      bool `default:"false"`

	MotorSteps           int           `default:"200"`
	SphereDiameter       Meters        `default:"0.2"`
	OmniwheelDiameter    Meters        `default:"0.048"`
	MinStepInterval      time.Duration `default:"400ms"`
	MotorAutoSleepLeeway time.Duration `default:"5ms"`

	Home LLACoords `ignored:"true"`
}

func ParseEnv() (Config, error) {
	var env Config
	err := envconfig.Process("corviator", &env)
	if err != nil {
		return env, err
	}

	env.Home = LLACoords{
		Φ: env.HomeLatitude,
		Λ: env.HomeLongitude,
		A: env.HomeAltitude,
	}

	return env, nil
}
