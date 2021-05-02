package deliveroo

import (
	"bytes"
	"github.com/jphastings/jan-poka/pkg/math"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

var (
	orderDelivered = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":{"attributes":{"message":"Your order took just 16 minutes. Enjoy!"}},"included":[{"id":"12345","type":"order","attributes":{"items":"Build your own + one more item","restaurant_name":"🏄Honi Poke - Angel🏄"}},{"id":"12345","type":"order_banner","attributes":{}}]}`)),
		Header:     make(http.Header),
	}
	orderArriving = &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"data":{"attributes":{"message":"Thales is nearby"}},"included":[{"id":"12345","type":"order","attributes":{"items":"Build your own + one more item","restaurant_name":"🏄Honi Poke - Angel🏄"}},{"id":"12345","type":"order_banner","attributes":{}},{"id":"customer-67890","type":"location","attributes":{"latitude":51.6,"longitude":-0.1,"type":"CUSTOMER"}},{"id":"rider-54321","type":"location","attributes":{"latitude":51.5,"longitude":-0.1,"type":"RIDER"}},{"id":"54321","type":"rider","attributes":{}}]}`)),
		Header:     make(http.Header),
	}
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req), nil }

func MockClient(t *testing.T, testURL string, response *http.Response) *http.Client {
	return &http.Client{Transport: RoundTripFunc(func(req *http.Request) *http.Response {
		if req.URL.String() != testURL {
			t.Errorf("query made to incorrect URL: %s", req.URL.String())
		}
		return response
	})}
}

func Test_config_Location(t *testing.T) {
	tests := []struct {
		name      string
		response  *http.Response
		wantLLA   math.LLACoords
		wantName  string
		wantRetry bool
	}{
		{"Arriving order", orderArriving, math.LLACoords{Latitude: 51.5, Longitude: -0.1}, "Thales is nearby", true},
		{"Delivered order", orderDelivered, math.LLACoords{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryURL := "https://example.com/order/status"
			c := &config{
				http:           MockClient(t, queryURL, tt.response),
				orderStatusURL: queryURL,
			}

			details, retry, _ := c.Location()
			if !reflect.DeepEqual(details.Coords, tt.wantLLA) {
				t.Errorf("Location() got = %v, want %v", details.Coords, tt.wantLLA)
			}
			if details.Name != tt.wantName {
				t.Errorf("Location() got = %v, want %v", details.Name, tt.wantName)
			}
			if retry != tt.wantRetry {
				t.Errorf("Location() got = %v, want %v", retry, tt.wantRetry)
			}
		})
	}
}
