package common

import "github.com/jphastings/jan-poka/pkg/math"

type LocationProvider interface {
	// SetParams provides you with a function that will populate the given (json annotated) struct pointer with the JSON params. May return error if JSON does not match. Should validate given parameters and return error if unusable.
	SetParams(func(decodeInto interface{}) error) error
	// Location returns the location according to the params set earlier at the current time. Second argument can be false if provider is currently offline.
	Location() (target math.LLACoords, suggestedName string, isUsable bool)
}

var Providers = make(map[string]func() LocationProvider)
