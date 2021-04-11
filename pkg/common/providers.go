package common

import (
	"github.com/jphastings/jan-poka/pkg/math"
	"time"
)

type LocationProvider interface {
	// SetParams provides you with a function that will populate the given (json annotated) struct pointer with the JSON params. May return error if JSON does not match. Should validate given parameters and return error if unusable.
	SetParams(func(decodeInto interface{}) error) error
	// Location returns the location according to the params set earlier. Will aim to get the most recent location, and will specify the time in accurateAt. Second argument can be false if provider can't provide a reasonable fix.
	Location() (target math.LLACoords, accurateAt time.Time, suggestedName string, isUsable bool)
}

var Providers = make(map[string]func() LocationProvider)
