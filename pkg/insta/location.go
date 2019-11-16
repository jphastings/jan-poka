package insta

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	. "github.com/jphastings/jan-poka/pkg/math"
	"log"
	"time"
)

type Post struct {
	FullName   string
	TakenAt    time.Time
	locationID string
	Coords     LLACoords
}

type sharedDataJSON struct {
	Rhxgis    string `json:"rhx_gis"`
	EntryData struct {
		ProfilePage []struct {
			GraphQL struct {
				User struct {
					FullName                 string `json:"full_name"`
					EdgeOwnerToTimelineMedia struct {
						Edges []struct {
							Node struct {
								TakenAtTimestamp int64 `json:"taken_at_timestamp"`
								Location         struct {
									ID string `json:"id"`
								} `json:"location"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"edge_owner_to_timeline_media"`
				} `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}

func GetLatest(username string, freshness time.Duration) (Post, error) {
	post, err := getLatestPostDetails(username)
	if err != nil {
		return Post{}, err
	}

	return addCoords(post)
}

func addCoords(post Post) (Post, error) {
	url := "https://www.instagram.com/explore/locations/" + post.locationID

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", url)
		if r.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", r.Ctx.Get("gis"), r.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			r.Headers.Set("X-Instagram-GIS", gisHash)
		}

	})

	c.OnHTML(`meta`, func(e *colly.HTMLElement) {
		fmt.Println(e, e.Attr("property"), e.Text)
	})

	err := c.Visit(url)
	return post, err
}

func getLatestPostDetails(username string) (Post, error) {
	url := "https://www.instagram.com/" + username
	foundPost := make(chan Post, 1)
	forceStop := make(chan error, 1)

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	err := c.Post("http://example.com/login", map[string]string{"username": "admin", "password": "admin"})
	if err != nil {
		log.Fatal(err)
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", url)
		if r.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", r.Ctx.Get("gis"), r.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			r.Headers.Set("X-Instagram-GIS", gisHash)
		}

	})

	c.OnHTML(`body`, func(e *colly.HTMLElement) {
		js := e.ChildText("script:first-of-type")
		jsonData := js[21 : len(js)-1]

		var sharedData sharedDataJSON
		if err := json.Unmarshal([]byte(jsonData), &sharedData); err != nil {
			panic(err)
		}

		profiles := sharedData.EntryData.ProfilePage
		if len(profiles) == 0 {
			panic("no profiles")
		}
		user := profiles[0].GraphQL.User

		for _, edge := range user.EdgeOwnerToTimelineMedia.Edges {
			p := Post{
				FullName:   user.FullName,
				TakenAt:    time.Unix(edge.Node.TakenAtTimestamp, 0),
				locationID: edge.Node.Location.ID,
			}

			if p.locationID != "" {
				foundPost <- p
			}
		}
	})

	go func() {
		err := c.Visit(url)
		if err == nil {
			err = fmt.Errorf("no posts with locations found")
		}
		forceStop <- err
	}()

	for {
		select {
		case post := <-foundPost:
			return post, nil
		case err := <-forceStop:
			return Post{}, err
		default:
		}
	}
}
