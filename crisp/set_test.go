package crisp

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSet(t *testing.T) {
	Convey("new set with dx", t, func() {
		Convey("when ok", func() {
			set, err := NewSet(0, 1, 0.25)
			So(err, ShouldBeNil)
			So(set, ShouldResemble, Set{
				xmin: 0.0,
				xmax: 1.0,
				dx:   0.25,
			})
		})

		Convey("when values", func() {
			set, _ := NewSet(0, 1, 0.25)
			So(set.Values(), ShouldResemble, []float64{0.0, 0.25, 0.5, 0.75, 1.0})
		})

		Convey("when values with error introduction", func() {
			set, _ := NewSet(0, 0.5, 0.1)
			So(set.Values(), ShouldResemble, []float64{0.0, 0.1, 0.2, 0.30000000000000004, 0.4, 0.5})
		})

		Convey("when empty", func() {
			So(Set{}.Values(), ShouldBeEmpty)
		})
	})

	Convey("new set with error", t, func() {
		Convey("when dx=0", func() {
			set, err := NewSet(0, 1, 0)
			So(err, ShouldBeError, "crisp set: dx shall be > 0")
			So(set, ShouldResemble, Set{})
		})

		Convey("when xmin>xmax", func() {
			set, err := NewSet(1, 0, 0.1)
			So(err, ShouldBeError, "crisp set: xmin shall be < xmax")
			So(set, ShouldResemble, Set{})
		})
	})
}
