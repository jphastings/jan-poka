package deliveroo

import (
	"encoding/json"
	"fmt"
	. "github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/math"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	TYPE              = "deliveroo"
	orderStatusAPIURL = "https://api.%s.deliveroo.com/consumer/v2-6/consumer_order_statuses/%s?sharing_token=%s"
)

var orderStatusPathRe = regexp.MustCompile(`^/orders/(\d+)/status$`)

type config struct {
	http             *http.Client
	orderStatusURL   string
	previousShareURL string
}

type request struct {
	ShareURL string `json:"share_url"`
}

func init() {
	Providers[TYPE] = func() LocationProvider {
		return &config{
			http: &http.Client{
				// Prevent following redirects
				CheckRedirect: func(*http.Request, []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
		}
	}
	log.Println("✅ Provider: Deliveroo order positions")
}

func (c *config) SetParams(decodeInto func(interface{}) error) error {
	var req request
	if err := decodeInto(&req); err != nil {
		return err
	}
	if req.ShareURL == c.previousShareURL {
		return nil
	}

	res, err := c.http.Head(req.ShareURL)
	if err != nil {
		return fmt.Errorf("given share url wasn't accessible: %w", err)
	}

	orderStatusURL, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		return fmt.Errorf("given share url wasn't usable: %w", err)
	}
	params, err := url.ParseQuery(orderStatusURL.RawQuery)
	if err != nil {
		return fmt.Errorf("given share url wasn't usable: %w", err)
	}
	matches := orderStatusPathRe.FindStringSubmatch(orderStatusURL.Path)
	if matches == nil {
		return fmt.Errorf("given share url wasn't usable")
	}

	start := len(orderStatusURL.Host) - 2
	end := len(orderStatusURL.Host)
	market := orderStatusURL.Host[start:end]

	c.orderStatusURL = fmt.Sprintf(orderStatusAPIURL, market, matches[1], params.Get("sharing_token"))
	c.previousShareURL = req.ShareURL
	return nil
}

type status struct {
	Data struct {
		Attributes struct {
			Message     string `json:"message"`
			IsCompleted bool   `json:"is_completed"`
			IsFailed    bool   `json:"is_failed"`
		} `json:"attributes"`
	} `json:"data"`
	Included []struct {
		Type       string `json:"type"`
		Attributes struct {
			Type      string       `json:"type"`
			Latitude  math.Degrees `json:"latitude"`
			Longitude math.Degrees `json:"longitude"`

			RestaurantName string `json:"restaurant_name"`
			Items          string `json:"items"`
		} `json:"attributes"`
	} `json:"included"`
}

func (c *config) Location() (TargetDetails, bool, error) {
	req, err := http.NewRequest("GET", c.orderStatusURL, nil)
	if err != nil {
		return TargetDetails{}, false, err
	}
	req.Header.Set("Accept", "application/json, application/vnd.api+json")

	res, err := c.http.Do(req)
	if err != nil {
		return TargetDetails{}, true, err
	}

	var s status
	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&s); err != nil {
		return TargetDetails{}, false, err
	}

	if s.Data.Attributes.IsCompleted || s.Data.Attributes.IsFailed {
		return TargetDetails{}, false, fmt.Errorf("order is no longer active")
	}

	for _, i := range s.Included {
		if i.Type != "location" || i.Attributes.Type != "RIDER" {
			continue
		}

		details := TargetDetails{
			Name: s.Data.Attributes.Message,
			Coords: math.LLACoords{
				Latitude:  i.Attributes.Latitude,
				Longitude: i.Attributes.Longitude,
			},
			AccurateAt: time.Now(),
		}

		return details, true, nil
	}

	return TargetDetails{}, true, fmt.Errorf("rider has not picked up order yet")
}
