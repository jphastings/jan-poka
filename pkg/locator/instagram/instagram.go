package instagram

import (
	"github.com/ahmdrz/goinsta"
	"github.com/jphastings/jan-poka/pkg/locator/common"
	. "github.com/jphastings/jan-poka/pkg/math"
	"log"
)

const TYPE = "instagram"

type locationProvider struct {
	Name      string   `json:"name"`
	Username  string   `json:"username"`
}

var insta *goinsta.Instagram

func init() {
	insta = goinsta.New("pilenticular", "e3Qp3yE4WgCpKvmeuF")
	err := insta.Login()
	if err != nil {
		log.Println("❌ Provider: Instagram could not log in.")
		return
	}

	common.Providers[TYPE] = func() common.LocationProvider { return &locationProvider{} }
	log.Println("✅ Provider: Instagram post positions available.")
}

func (lp *locationProvider) SetParams(decodeInto func(interface{}) error) error {
	if err := decodeInto(lp); err != nil {
		return err
	}

	user, err := insta.Profiles.ByName(lp.Username)
	if err != nil {
		return err
	}

	if lp.Name == "" {
		lp.Name = user.FullName
	}

	if lp.Name == "" {
		lp.Name = lp.Username
	}

	//fmt.Println(feed)
	//for _, item := range feed.Items {
	//	fmt.Println(item.Lat, item.Lng)
	//}

	return nil
}

func (lp *locationProvider) Location() (LLACoords, string, bool) {
	return LLACoords{
		Latitude:  0,
		Longitude: 0,
		Altitude:  0,
	}, lp.Username, true
}
