package transforms_test

import (
	"testing"

	"github.com/jphastings/corviator/pkg/transforms"

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

var _ = Describe("AbsoluteDirection", func() {
	quarterTurn := Radians(Π / 2)
	halfTurn := Radians(Π)

	DescribeTable("known directions",
		func(home LLACoords, target LLACoords, expected SphericalCoords) {
			actual := transforms.AbsoluteDirection(home, target)

			Expect(actual.R).To(BeNumerically("~", expected.R, Meters(1)))
			// Angle down from ECEF Z
			Expect(actual.Θ).To(BeNumerically("~", expected.Θ, arcSecond.Radians()))
			// Pointing straight up or down the angle around is irrelevant
			if expected.Θ != 0 && expected.Θ != halfTurn {
				// Angle anticlockwise round from ECEF X looking from ECEF +Z
				Expect(actual.Φ).To(BeNumerically("~", expected.Φ, arcSecond.Radians()))
			}
		},
		// NB. 1 arcSecond =~ 31m at the equator
		Entry("North at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: arcSecond, Λ: 0, A: 0},
			SphericalCoords{R: 31, Θ: 0}),
		Entry("East at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: arcSecond, A: 0},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: quarterTurn}),
		Entry("South at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: -arcSecond, Λ: 0, A: 0},
			SphericalCoords{R: 31, Θ: halfTurn}),
		Entry("West at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: -arcSecond, A: 0},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: 3 * quarterTurn}),
		Entry("Up at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: 0, A: 31},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: 0}),
		Entry("Down at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: 0, A: -31},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: halfTurn}),

		Entry("North at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: arcSecond, Λ: 90, A: 0},
			SphericalCoords{R: 31, Θ: 0}),
		Entry("East at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90 + arcSecond, A: 0},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: halfTurn}),
		Entry("South at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: -arcSecond, Λ: 90, A: 0},
			SphericalCoords{R: 31, Θ: halfTurn}),
		Entry("West at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90 - arcSecond, A: 0},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: 2 * Π}), // This only isn't zero because it's _just_ under 2 pi
		Entry("Up at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90, A: 31},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: quarterTurn}),
		Entry("Down at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90, A: -31},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: 3 * quarterTurn}),

		Entry("Dateline from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90 - arcSecond, Λ: 180, A: 0},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: halfTurn}),
		Entry("Towards Russia from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90 - arcSecond, Λ: 90, A: 0},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: quarterTurn}),
		Entry("Meridian from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90 - arcSecond, Λ: 0, A: 0},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: 0}),
		Entry("Towards Canada from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90 - arcSecond, Λ: -90, A: 0},
			SphericalCoords{R: 31, Θ: quarterTurn, Φ: 3 * quarterTurn}),
		Entry("Up from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90, Λ: 0, A: 31},
			SphericalCoords{R: 31, Θ: 0}),
		Entry("Down from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90, Λ: 0, A: -31},
			SphericalCoords{R: 31, Θ: halfTurn}),

		Entry("Across the dateline, directionally looking at America from NZ",
			LLACoords{Φ: 0, Λ: 180 - arcSecond, A: 0},
			LLACoords{Φ: 0, Λ: arcSecond - 180, A: 0},
			SphericalCoords{R: 62, Θ: quarterTurn, Φ: 3 * quarterTurn}),
		Entry("Across the dateline, directionally looking at NZ from America",
			LLACoords{Φ: 0, Λ: arcSecond - 180, A: 0},
			LLACoords{Φ: 0, Λ: 180 - arcSecond, A: 0},
			SphericalCoords{R: 62, Θ: quarterTurn, Φ: quarterTurn}),
	)
})

var _ = Describe("RelativeDirection", func() {

	DescribeTable("known directions",
		func(home, target LLACoords, facing Degrees, expected Direction) {
			actual := transforms.RelativeDirection(home, target, facing)

			Expect(actual.Distance).To(BeNumerically("~", expected.Distance, Meters(1)))
			// Angle up (+) or down (-) from the horizon
			Expect(actual.Elevation).To(BeNumerically("~", expected.Elevation, 2*arcSecond))
			// Pointing straight up or down the angle around is irrelevant
			if expected.Elevation != 90 && expected.Elevation != -90 {
				// Angle clockwise round from the direction specified by 'facing'
				Expect(actual.Heading).To(BeNumerically("~", expected.Heading, 2*arcSecond))
			}
		},
		Entry("North at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: arcSecond, Λ: 0, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 0, Elevation: 0}),
		Entry("East at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: arcSecond, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 90, Elevation: 0}),
		Entry("South at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: -arcSecond, Λ: 0, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 180, Elevation: 0}),
		Entry("West at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: -arcSecond, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 270, Elevation: 0}),
		Entry("Up at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: 0, A: 31},
			Degrees(0),
			Direction{Distance: 31, Elevation: 90}),
		Entry("Down at equator/meridian",
			LLACoords{Φ: 0, Λ: 0, A: 0},
			LLACoords{Φ: 0, Λ: 0, A: -31},
			Degrees(0),
			Direction{Distance: 31, Elevation: -90}),

		Entry("North at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: arcSecond, Λ: 90, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 0, Elevation: 0}),
		Entry("East at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90 + arcSecond, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 90, Elevation: 0}),
		Entry("South at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: -arcSecond, Λ: 90, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 180, Elevation: 0}),
		Entry("West at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90 - arcSecond, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 270, Elevation: 0}),
		Entry("Up at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90, A: 31},
			Degrees(0),
			Direction{Distance: 31, Elevation: 90}),
		Entry("Down at equator/90 long",
			LLACoords{Φ: 0, Λ: 90, A: 0},
			LLACoords{Φ: 0, Λ: 90, A: -31},
			Degrees(0),
			Direction{Distance: 31, Elevation: -90}),

		Entry("Dateline from North Pole",
			LLACoords{Φ: 90 - arcSecond, Λ: 180, A: 0},
			LLACoords{Φ: 90 - 2*arcSecond, Λ: 180, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 180, Elevation: 0}),
		Entry("Towards Russia from North Pole",
			LLACoords{Φ: 90 - arcSecond, Λ: 90, A: 0},
			LLACoords{Φ: 90 - 2*arcSecond, Λ: 90, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 180, Elevation: 0}),
		Entry("Meridian from North Pole",
			LLACoords{Φ: 90 - arcSecond, Λ: 0, A: 0},
			LLACoords{Φ: 90 - 2*arcSecond, Λ: 0, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 180, Elevation: 0}),
		Entry("Towards Canada from North Pole",
			LLACoords{Φ: 90 - arcSecond, Λ: -90, A: 0},
			LLACoords{Φ: 90 - 2*arcSecond, Λ: -90, A: 0},
			Degrees(0),
			Direction{Distance: 31, Heading: 180, Elevation: 0}),
		Entry("Up from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90, Λ: 0, A: 31},
			Degrees(0),
			Direction{Distance: 31, Elevation: 90}),
		Entry("Down from North Pole",
			LLACoords{Φ: 90, Λ: 0, A: 0},
			LLACoords{Φ: 90, Λ: 0, A: -31},
			Degrees(0),
			Direction{Distance: 31, Elevation: -90}),

		Entry("Across the dateline, directionally looking at America from NZ",
			LLACoords{Φ: 0, Λ: 180 - arcSecond, A: 0},
			LLACoords{Φ: 0, Λ: arcSecond - 180, A: 0},
			Degrees(0),
			Direction{Distance: 62, Heading: 90, Elevation: 0}),
		Entry("Across the dateline, directionally looking at NZ from America",
			LLACoords{Φ: 0, Λ: arcSecond - 180, A: 0},
			LLACoords{Φ: 0, Λ: 180 - arcSecond, A: 0},
			Degrees(0),
			Direction{Distance: 62, Heading: 270, Elevation: 0}),
	)
})
