package google

import (
	"context"
	"fmt"
	. "github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/math"
	"googlemaps.github.io/maps"
	"time"
)

const TYPE = "google"

type config struct {
	client *maps.Client
	target request
}

type request struct {
	Name       string `json:"name"`
	SearchTerm string `json:"q"`
}

func Login(apiKey string) error {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return err
	}
	c := &config{client: client}

	Providers[TYPE] = func() LocationProvider { return c }
	return nil
}

func (c *config) SetParams(decodeInto func(interface{}) error) error {
	return decodeInto(&c.target)
}

func (c *config) Location() TargetDetails {
	ctx := context.Background()

	fields := []maps.PlaceSearchFieldMask{maps.PlaceSearchFieldMaskGeometryLocation}
	if c.target.Name == "" {
		fields = append(fields, maps.PlaceSearchFieldMaskName)
	}

	places, err := c.client.FindPlaceFromText(ctx, &maps.FindPlaceFromTextRequest{
		Input:     c.target.SearchTerm,
		InputType: maps.FindPlaceFromTextInputTypeTextQuery,
		Fields:    fields,
	})
	if err != nil {
		return TargetDetails{Final: true, Err: err}
	}
	if len(places.Candidates) == 0 {
		return TargetDetails{Final: true, Err: fmt.Errorf("no search results")}
	}

	place := places.Candidates[0]
	name := c.target.Name
	if name == "" {
		name = place.Name
	}

	lla := math.LLACoords{
		Latitude:  math.Degrees(place.Geometry.Location.Lat),
		Longitude: math.Degrees(place.Geometry.Location.Lng),
	}
	_ = c.GuessElevation(&lla)
	return TargetDetails{
		Name:       name,
		Coords:     lla,
		AccurateAt: time.Now(),
		Final:      true,
	}
}

func (c *config) GuessElevation(lla *math.LLACoords) error {
	elevations, err := c.client.Elevation(context.Background(), &maps.ElevationRequest{
		Locations: []maps.LatLng{{
			Lat: float64(lla.Latitude),
			Lng: float64(lla.Longitude),
		}},
	})
	if err != nil {
		return err
	}
	if len(elevations) == 0 {
		return fmt.Errorf("no elevation data available")
	}

	lla.Altitude = math.Meters(elevations[0].Elevation)
	return nil
}
