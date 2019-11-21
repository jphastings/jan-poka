package math_test

import (
	"testing"

	. "github.com/jphastings/jan-poka/pkg/math"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestTrigonometry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trigonometry Suite")
}

var _ = Describe("ModDeg", func() {

	DescribeTable("known angles",
		func(deg Degrees, expected Degrees) {
			Expect(ModDeg(deg)).To(Equal(expected))

		},
		Entry("0º", Degrees(0), Degrees(0)),
		Entry("15º", Degrees(15), Degrees(15)),
		Entry("360º", Degrees(360), Degrees(0)),
		Entry("540º", Degrees(540), Degrees(180)),
		Entry("720º", Degrees(720), Degrees(0)),
		Entry("-45º", Degrees(-45), Degrees(315)),
		Entry("-360º", Degrees(-360), Degrees(0)),
		Entry("-540º", Degrees(-540), Degrees(180)),
		Entry("-720º", Degrees(-720), Degrees(0)),
	)
})
