package instagram

import (
	"github.com/jphastings/jan-poka/pkg/locator/common"
	. "github.com/jphastings/jan-poka/pkg/math"
	"time"

	"github.com/ahmdrz/goinsta/v2"
)

const TYPE = "instagram"

type config struct {
	client *goinsta.Instagram
	target request
	user   *goinsta.User
}

type request struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func Login(username, password string) error {
	c := &config{client: goinsta.New(username, password)}
	if err := c.client.Login(); err != nil {
		return err
	}

	common.Providers[TYPE] = func() common.LocationProvider { return c }
	return nil
}

func (c *config) SetParams(decodeInto func(interface{}) error) error {
	if err := decodeInto(&c.target); err != nil {
		return err
	}

	user, err := c.client.Profiles.ByName(c.target.Username)
	if err != nil {
		return err
	}
	c.user = user

	if c.target.Name == "" {
		c.target.Name = user.FullName
	}

	if c.target.Name == "" {
		c.target.Name = c.target.Username
	}

	return nil
}

func (c *config) Location() (LLACoords, time.Time, string, bool) {
	feed := c.user.Feed()
	// If a location isn't in the first page, assume it'd be too old
	feed.Next(false)

	for _, item := range feed.Items {
		// Imported photos aren't recent, so ignore them
		// goinsta doesn't specify if there's no lat/long, so assume that no-one will ever post at 0,0
		if item.ImportedTakenAt == 0 || item.Lat == 0 && item.Lng == 0 {
			continue
		}

		accurateAt := time.Unix(item.TakenAt, 0)

		return LLACoords{
			Latitude:  Degrees(item.Lat),
			Longitude: Degrees(item.Lng),
			Altitude:  0,
		}, accurateAt, c.target.Name, true
	}

	return LLACoords{}, time.Time{}, "", false
}
