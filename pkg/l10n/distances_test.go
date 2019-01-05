package l10n_test

import (
	"github.com/jphastings/corviator/pkg/l10n"
	"testing"

	. "github.com/jphastings/corviator/pkg/math"
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
		// Meters
		Entry("Zero", Meters(0), "0m"),
		Entry("Exact", Meters(1), "1m"),
		Entry("Round down", Meters(1.1), "1m"),
		Entry("Quarter round up", Meters(1.2), "1¼m"),
		Entry("Quarter exact", Meters(1.25), "1¼m"),
		Entry("Quarter round down", Meters(1.3), "1¼m"),
		Entry("Half round up", Meters(1.4), "1½m"),
		Entry("Half exact", Meters(1.5), "1½m"),
		Entry("Half round down", Meters(1.6), "1½m"),
		Entry("Three quarters round up", Meters(1.7), "1½m"),
		Entry("Three quarters exact", Meters(1.75), "2m"),
		Entry("Three quarters round down", Meters(1.8), "2m"),
		Entry("Round up", Meters(1.9), "2m"),
		Entry("Three significant figures", Meters(900), "990m"),

		// Kilometers
		Entry("Three significant figures", Meters(999), "1km"),
		Entry("Exact", Meters(1000), "1km"),
	)
})
