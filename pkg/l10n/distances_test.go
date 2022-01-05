package l10n_test

import (
	"testing"

	"github.com/jphastings/jan-poka/pkg/l10n"

	. "github.com/jphastings/jan-poka/pkg/math"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestDistances(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Distances Suite")
}

var _ = Describe("LLACoords.Distance", func() {
	DescribeTable("important distances",
		func(distance Meters, description string) {
			actual := l10n.Distance(distance)
			Expect(actual).To(Equal(description))
		},

		// Hours walk
		Entry("Very close", Meters(1), "0 hours' walk away"),
		Entry("A quarter hour walk", Meters(1260), "¼ hours' walk away"),
		Entry("Half an hour walk", Meters(2520), "½ hours' walk away"),
		Entry("An hour walk", Meters(5040), "1 hours' walk away"),
		Entry("Ten hour walk", Meters(50400), "10 hours' walk away"),

		// // Hours Drive
		Entry("An hour's drive", Meters(72420), "1 hours' drive away"),
		Entry("Two hours' drive", Meters(144840), "2 hours' drive away"),

		// // Hours Flight
		Entry("An hour's flight", Meters(885139), "1 hours' flight away"),
	)
})
