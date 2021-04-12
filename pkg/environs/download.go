package environs

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/jphastings/jan-poka/pkg/math"
	geojson "github.com/paulmach/go.geojson"
	"io"
	"strconv"
	"strings"

	"github.com/tidwall/buntdb"
)

type database struct {
	client *buntdb.DB
}

func New(dbFile string) (*database, error) {
	db, err := buntdb.Open(dbFile)
	if err != nil {
		return nil, err
	}
	if err := db.CreateSpatialIndex("geo", "*", buntdb.IndexRect); err != nil {
		return nil, err
	}
	return &database{client: db}, nil
}

func (d *database) At(lla math.LLACoords) error {
	return d.client.View(func(tx *buntdb.Tx) error {
		rect := fmt.Sprintf("[%f %f]", float64(lla.Latitude), float64(lla.Longitude))
		return tx.Nearby("geo", rect, func(key, val string, dist float64) bool {
			fmt.Printf("%s, %+v\n", key, val)
			return false
		})
	})
}

func (d *database) BuildDB(gzCloser io.ReadCloser) error {
	closer, err := gzip.NewReader(gzCloser)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(closer)
	for scanner.Scan() {
		rec := strings.Split(scanner.Text(), "\t")
		if err := importRecord(d.client, rec); err != nil {
			return err
		}
	}

	return scanner.Err()
}

const (
	countFields    = 19
	fieldID        = 0
	fieldLat       = 4
	fieldLng       = 5
	fieldShortCode = 6
	fieldCode      = 7
)

var requiredFields = []int{fieldID, fieldLat, fieldLng, fieldCode, fieldShortCode}

func importRecord(db *buntdb.DB, line []string) error {
	if len(line) != countFields {
		return nil
	}
	for _, req := range requiredFields {
		if len(line) <= req || line[req] == "" {
			return nil
		}
	}

	id := line[fieldID]

	lat, err := strconv.ParseFloat(line[fieldLat], 64)
	if err != nil {
		return fmt.Errorf("latitude (%s) not parseable: %w", line[fieldLat], err)
	}
	lng, err := strconv.ParseFloat(line[fieldLng], 64)
	if err != nil {
		return fmt.Errorf("longitude (%s) not parseable: %w", line[fieldLng], err)
	}

	data, err := toGeoJSON(lat, lng, line[fieldShortCode], line[fieldCode])
	if err != nil {
		return fmt.Errorf("couldn't marshal geojson: %w", err)
	}

	return db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(id, data, nil)
		return err
	})
}

func toGeoJSON(lat, lng float64, shortCode, code string) (string, error) {
	f := geojson.NewFeature(geojson.NewPointGeometry([]float64{lat, lng}))
	f.Properties["gn.s"] = shortCode
	f.Properties["gn.c"] = code

	rawJSON, err := f.MarshalJSON()
	return string(rawJSON), err
}
