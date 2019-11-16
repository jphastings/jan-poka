package instagram

import (
	"fmt"
	"github.com/ahmdrz/goinsta"
	"time"

	"github.com/jphastings/jan-poka/pkg/locator/common"
	. "github.com/jphastings/jan-poka/pkg/math"
)

const TYPE = "instagram"

type locationProvider struct {
	Name      string   `json:"name"`
	Username  string   `json:"username"`
	Freshness Duration `json:"freshness"`
}

var insta *goinsta.Instagram

// TODO: Could this be replaced by init()?
func Load() {
	insta = goinsta.New("pilenticular", "e3Qp3yE4WgCpKvmeuF")
	insta.Login()

	common.Providers[TYPE] = func() common.LocationProvider { return &locationProvider{} }
}

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	lp.Freshness = Duration{time.Duration(24) * time.Hour}
	if err := decodeInto(lp); err != nil {
		return err
	}

	if lp.Name == "" {
		lp.Name = lp.Username
	}

	user, err := insta.Profiles.ByName(lp.Username)
	if err != nil {
		return err
	}

	fmt.Println(user)

	feed := user.Feed() // time.Now().Add(lp.Freshness.Duration * time.Duration(-1)).String()

	fmt.Println(feed)
	for _, item := range feed.Items {
		fmt.Println(item.Lat, item.Lng)
	}

	return nil
}

func (lp *locationProvider) Location() (LLACoords, string, bool) {
	return LLACoords{
		Latitude:  0,
		Longitude: 0,
		Altitude:  0,
	}, lp.Username, true
}
