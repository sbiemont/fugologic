package crisp

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSet(t *testing.T) {
	Convey("new set with dx", t, func() {
		set := NewSetDx(0, 1, 0.25)

		Convey("when ok", func() {
			So(set, ShouldResemble, Set{
				xmin: 0.0,
				xmax: 1.0,
				dx:   0.25,
			})
		})

		Convey("when values", func() {
			So(set.Values(), ShouldResemble, []float64{0.0, 0.25, 0.5, 0.75, 1.0})
		})
	})
}
