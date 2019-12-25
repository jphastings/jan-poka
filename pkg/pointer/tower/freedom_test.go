package tower_test

import (
	"github.com/jphastings/jan-poka/pkg/math"
	"github.com/jphastings/jan-poka/pkg/pointer/tower"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

func TestFreedom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Freedom Suite")
}

var _ = Describe("Pointer", func() {

	DescribeTable("known angles",
		func(theta, phi math.Degrees, base, arm math.Degrees) {
			actualBase, actualArm := tower.Pointer(theta, phi)
			Expect(actualBase).To(Equal(base))
			Expect(actualArm).To(Equal(arm))
		},
		Entry("0º", math.Degrees(0), math.Degrees(90), math.Degrees(-90), math.Degrees(90)),
		Entry("1º", math.Degrees(1), math.Degrees(90), math.Degrees(-89), math.Degrees(90)),
		Entry("89º", math.Degrees(89), math.Degrees(90), math.Degrees(-1), math.Degrees(90)),
		Entry("90º", math.Degrees(90), math.Degrees(90), math.Degrees(0), math.Degrees(90)),
		Entry("91º", math.Degrees(91), math.Degrees(90), math.Degrees(1), math.Degrees(90)),
		Entry("179º", math.Degrees(179), math.Degrees(90), math.Degrees(89), math.Degrees(90)),
		Entry("180º", math.Degrees(180), math.Degrees(90), math.Degrees(-90), math.Degrees(-90)),
		Entry("181º", math.Degrees(181), math.Degrees(90), math.Degrees(-89), math.Degrees(-90)),
		Entry("269º", math.Degrees(269), math.Degrees(90), math.Degrees(-1), math.Degrees(-90)),
		Entry("270º", math.Degrees(270), math.Degrees(90), math.Degrees(0), math.Degrees(-90)),
		Entry("271º", math.Degrees(271), math.Degrees(90), math.Degrees(1), math.Degrees(-90)),
		Entry("359º", math.Degrees(359), math.Degrees(90), math.Degrees(89), math.Degrees(-90)),
	)
})
