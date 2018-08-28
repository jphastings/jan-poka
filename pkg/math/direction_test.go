package math_test

import (
	"testing"

	. "github.com/jphastings/corviator/pkg/math"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestDirections(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Directions Suite")
}

const arcSecond = Degrees(0.00027778)

var _ = FDescribe("LLACoords.DirectionTo", func() {
	cancer := Degrees(23.43687)

	DescribeTable("known directions",
		func(home, target LLACoords, facing Degrees, expected AERCoords) {
			actual := home.DirectionTo(target, facing)

			Expect(actual.Range).To(BeNumerically("~", expected.Range, Meters(1)))
			// Angle up (+) or down (-) from the horizon
			Expect(actual.Elevation).To(BeNumerically("~", expected.Elevation, 4*arcSecond))
			// Pointing straight up or down the angle around is irrelevant
			if expected.Elevation != 90 && expected.Elevation != -90 {
				// Angle clockwise round from the direction specified by 'facing'
				Expect(actual.Azimuth).To(BeNumerically("~", expected.Azimuth, 4*arcSecond))
			}
		},
		Entry("North at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: arcSecond, Λ: 0, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 0, Elevation: 0}),
		Entry("East at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: arcSecond, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 90, Elevation: 0}),
		Entry("South at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: -arcSecond, Λ: 0, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 180, Elevation: 0}),
		Entry("West at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: -arcSecond, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 270, Elevation: 0}),
		Entry("Up at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: 0, A: 31},
			Degrees(0),
			AERCoords{Range: 31, Elevation: 90}),
		Entry("Down at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: 0, A: -31},
			Degrees(0),
			AERCoords{Range: 31, Elevation: -90}),

		Entry("North at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: arcSecond, Λ: 90, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 0, Elevation: 0}),
		Entry("East at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90 + arcSecond, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 90, Elevation: 0}),
		Entry("South at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: -arcSecond, Λ: 90, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 180, Elevation: 0}),
		Entry("West at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90 - arcSecond, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 270, Elevation: 0}),
		Entry("Up at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90, A: 31},
			Degrees(0),
			AERCoords{Range: 31, Elevation: 90}),
		Entry("Down at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90, A: -31},
			Degrees(0),
			AERCoords{Range: 31, Elevation: -90}),

		Entry("North at tropic cancer/180 long",
			LLACoords{Φ: cancer, Λ: 180, A: 0},
			LLACoords{Φ: cancer + arcSecond, Λ: 180, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 0, Elevation: 0}),
		Entry("East at tropic cancer/180 long (across dateline)",
			LLACoords{Φ: cancer, Λ: 180, A: 0},
			LLACoords{Φ: cancer, Λ: arcSecond - 180, A: 0},
			Degrees(0),
			// The range is a little less than the height specified because of the WGS translation
			AERCoords{Range: 28, Azimuth: 90, Elevation: 0}),
		Entry("South at tropic cancer/180 long",
			LLACoords{Φ: cancer, Λ: 180, A: 0},
			LLACoords{Φ: cancer - arcSecond, Λ: 180, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 180, Elevation: 0}),
		Entry("West at tropic cancer/180 long (across dateline)",
			LLACoords{Φ: cancer, Λ: arcSecond - 180, A: 0},
			LLACoords{Φ: cancer, Λ: 180, A: 0},
			Degrees(0),
			// The range is a little less than the height specified because of the WGS translation
			AERCoords{Range: 28, Azimuth: 270, Elevation: 0}),
		Entry("Up at tropic cancer/180 long",
			LLACoords{Φ: cancer, Λ: 180, A: 0},
			LLACoords{Φ: cancer, Λ: 180, A: 31},
			Degrees(0),
			AERCoords{Range: 31, Elevation: 90}),
		Entry("Down at tropic cancer/180 long",
			LLACoords{Φ: cancer, Λ: 180, A: 0},
			LLACoords{Φ: cancer, Λ: 180, A: -31},
			Degrees(0),
			AERCoords{Range: 31, Elevation: -90}),

		Entry("Dateline from North Pole",
			LLACoords{Φ: 90 - arcSecond, Λ: 180, A: 0},
			LLACoords{Φ: 90 - 2*arcSecond, Λ: 180, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 180, Elevation: 0}),
		Entry("Towards Russia from North Pole",
			LLACoords{Φ: 90 - arcSecond, Λ: 90, A: 0},
			LLACoords{Φ: 90 - 2*arcSecond, Λ: 90, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 180, Elevation: 0}),
		Entry("Meridian from North Pole",
			LLACoords{Φ: 90 - arcSecond, Λ: 0, A: 0},
			LLACoords{Φ: 90 - 2*arcSecond, Λ: 0, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 180, Elevation: 0}),
		Entry("Towards Canada from North Pole",
			LLACoords{Φ: 90 - arcSecond, Λ: -90, A: 0},
			LLACoords{Φ: 90 - 2*arcSecond, Λ: -90, A: 0},
			Degrees(0),
			AERCoords{Range: 31, Azimuth: 180, Elevation: 0}),
		Entry("Up from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90, Λ: 0, A: 31},
			Degrees(0),
			AERCoords{Range: 31, Elevation: 90}),
		Entry("Down from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90, Λ: 0, A: -31},
			Degrees(0),
			AERCoords{Range: 31, Elevation: -90}),

		Entry("North at equator/meridian while facing East",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: arcSecond, Λ: 0, A: 0},
			Degrees(90),
			AERCoords{Range: 31, Azimuth: 270, Elevation: 0}),
		Entry("North at equator/meridian while facing West",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: arcSecond, Λ: 0, A: 0},
			Degrees(270),
			AERCoords{Range: 31, Azimuth: 90, Elevation: 0}),
	)
})
