package math_test

import (
	"testing"

	. "github.com/jphastings/jan-poka/pkg/math"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestCoordinate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Coordinate Suite")
}

var _ = Describe("LLACoords.ECEF()", func() {
	accuracy := Meters(1)

	DescribeTable("known coordinate transforms",
		func(input LLACoords, expected ECEFCoords) {
			actual := input.ECEF(WGS84)
			Expect(actual.X).To(BeNumerically("~", expected.X, accuracy))
			Expect(actual.Y).To(BeNumerically("~", expected.Y, accuracy))
			Expect(actual.Z).To(BeNumerically("~", expected.Z, accuracy))
		},
		Entry("surface X",
			LLACoords{Latitude: 0, Longitude: 0, Altitude: 0},
			ECEFCoords{X: 6378137, Y: 0, Z: 0}),
		Entry("surface Y",
			LLACoords{Latitude: 0, Longitude: 90, Altitude: 0},
			ECEFCoords{X: 0, Y: 6378137, Z: 0}),
		Entry("surface Z",
			LLACoords{Latitude: 90, Longitude: 0, Altitude: 0},
			ECEFCoords{X: 0, Y: 0, Z: 6356752}),
		Entry("surface -X",
			LLACoords{Latitude: 0, Longitude: 180, Altitude: 0},
			ECEFCoords{X: -6378137, Y: 0, Z: 0}),
		Entry("surface -Y",
			LLACoords{Latitude: 0, Longitude: -90, Altitude: 0},
			ECEFCoords{X: 0, Y: -6378137, Z: 0}),
		Entry("surface -Z",
			LLACoords{Latitude: -90, Longitude: 0, Altitude: 0},
			ECEFCoords{X: 0, Y: 0, Z: -6356752}),
		Entry("Greenwich observatory",
			LLACoords{Latitude: 51.4769, Longitude: 0.0005, Altitude: 48},
			ECEFCoords{X: 3980689, Y: 35, Z: 4966800}),
		Entry("Vernadsky Station Bar",
			LLACoords{Latitude: -65.245724, Longitude: -64.257668, Altitude: 4},
			ECEFCoords{X: 1163168, Y: -2412321, Z: -5769239}),
		Entry("Marist Brothers Primary School in Suva",
			LLACoords{Latitude: -18.140535, Longitude: 178.428644, Altitude: 33},
			ECEFCoords{X: -6060835, Y: 166262, Z: -1973182}),
	)
})
