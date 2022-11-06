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

	Convey("new set with n", t, func() {
		Convey("when ok", func() {
			set, err := NewSetN(0, 1, 5)
			So(err, ShouldBeNil)
			So(set, ShouldResemble, Set{
				xmin: 0.0,
				xmax: 1.0,
				dx:   0.25,
			})
		})

		Convey("when values", func() {
			set, _ := NewSetN(0, 1, 5)
			So(set.Values(), ShouldResemble, []float64{0.0, 0.25, 0.5, 0.75, 1.0})

			set, _ = NewSetN(0, 1, 2)
			So(set.Values(), ShouldResemble, []float64{0.0, 1.0})

			set, _ = NewSetN(0, 0.5, 4)
			So(set.Values(), ShouldResemble, []float64{0, 0.16666666666666666, 0.3333333333333333, 0.5})
		})
	})

	Convey("new set with error", t, func() {
		Convey("when dx", func() {
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

		Convey("when n", func() {
			Convey("when n<=1", func() {
				set, err := NewSetN(0, 1, 1)
				So(err, ShouldBeError, "crisp set: n shall be >= 2")
				So(set, ShouldResemble, Set{})
			})

			Convey("when xmin>xmax", func() {
				set, err := NewSetN(1, 0, 2)
				So(err, ShouldBeError, "crisp set: xmin shall be < xmax")
				So(set, ShouldResemble, Set{})
			})
		})
	})

	Convey("new set by changing n", t, func() {
		Convey("when ok", func() {
			set1, err1 := NewSetN(0, 1, 5)
			So(err1, ShouldBeNil)
			So(set1.Values(), ShouldResemble, []float64{0.0, 0.25, 0.5, 0.75, 1.0})

			set2, err2 := set1.N(3)
			So(err2, ShouldBeNil)
			So(set2.Values(), ShouldResemble, []float64{0.0, 0.5, 1.0})
		})

		Convey("when ok and change n<=1", func() {
			set1, err1 := NewSetN(0, 1, 5)
			So(err1, ShouldBeNil)
			So(set1.Values(), ShouldResemble, []float64{0.0, 0.25, 0.5, 0.75, 1.0})

			set2, err2 := set1.N(1)
			So(err2, ShouldBeError, "crisp set: n shall be >= 2")
			So(set2, ShouldResemble, Set{})
		})
	})

	Convey("new set by changing dx", t, func() {
		Convey("when ok", func() {
			set1, err1 := NewSet(0, 1, 0.25)
			So(err1, ShouldBeNil)
			So(set1.Values(), ShouldResemble, []float64{0.0, 0.25, 0.5, 0.75, 1.0})

			set2, err2 := set1.Dx(0.5)
			So(err2, ShouldBeNil)
			So(set2.Values(), ShouldResemble, []float64{0.0, 0.5, 1.0})
		})

		Convey("when ok and change n<=1", func() {
			set1, err1 := NewSet(0, 1, 0.25)
			So(err1, ShouldBeNil)
			So(set1.Values(), ShouldResemble, []float64{0.0, 0.25, 0.5, 0.75, 1.0})

			set2, err2 := set1.Dx(0)
			So(err2, ShouldBeError, "crisp set: dx shall be > 0")
			So(set2, ShouldResemble, Set{})
		})
	})
}
