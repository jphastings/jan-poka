package instagram

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/env"
	"github.com/jphastings/jan-poka/pkg/locator/common"
	. "github.com/jphastings/jan-poka/pkg/math"
	"log"

	"github.com/ahmdrz/goinsta/v2"
)

const TYPE = "instagram"

type config struct {
	client *goinsta.Instagram
	target request
}

type request struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func Login(environment env.Config) {
	c := &config{
		client: goinsta.New(environment.InstagramUsername, environment.InstagramPassword),
	}
	if err := c.client.Login(); err != nil {
		log.Println("❌ Provider: Instagram could not log in.")
		return
	}

	common.Providers[TYPE] = func() common.LocationProvider { return c }
	log.Println("✅ Provider: Instagram post positions available.")
}

func (c *config) SetParams(decodeInto func(interface{}) error) error {
	if err := decodeInto(&c.target); err != nil {
		return err
	}

	user, err := c.client.Profiles.ByName(c.target.Username)
	if err != nil {
		return err
	}

	if c.target.Name == "" {
		c.target.Name = user.FullName
	}

	if c.target.Name == "" {
		c.target.Name = c.target.Username
	}

	//fmt.Println(feed)
	//for _, item := range feed.Items {
	//	fmt.Println(item.Lat, item.Lng)
	//}

	// Haven't currently figured out how to extract lat/long
	return fmt.Errorf("not implemented")
}

func (c *config) Location() (LLACoords, string, bool) {
	return LLACoords{
		Latitude:  0,
		Longitude: 0,
		Altitude:  0,
	}, c.target.Name, true
}
