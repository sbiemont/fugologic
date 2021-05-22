package fuzzy

import (
	"testing"

	"fugologic.git/crisp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProviderGet(t *testing.T) {
	newSet := func() Set { return func(x float64) float64 { return x } }

	fvA := NewIDValCustom("a", crisp.Set{})
	fsA1 := NewIDSetCustom("a1", newSet(), &fvA)

	fvMissingParent := NewIDValCustom("b", crisp.Set{})
	fsMissingParent := NewIDSetCustom("b1", newSet(), &fvMissingParent)

	Convey("data", t, func() {
		Convey("when parent missing", func() {
			provider := DataInput{}
			fsMissingParent.parent = nil
			y, err := provider.find(fsMissingParent)
			So(err, ShouldBeError, "input: cannot find parent for id set `b1`")
			So(y, ShouldBeZeroValue)
		})

		Convey("when data missing", func() {
			provider := DataInput{}
			y, err := provider.find(fsA1)
			So(err, ShouldBeError, "input: cannot find data for id val `a` (id set `a1`)")
			So(y, ShouldBeZeroValue)
		})

		Convey("when ok", func() {
			provider := DataInput{"a": 1}
			y, err := provider.find(fsA1)
			So(err, ShouldBeNil)
			So(y, ShouldEqual, 1)
		})
	})
}
