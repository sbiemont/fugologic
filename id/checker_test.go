package id

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type mockID ID

func (mid mockID) ID() ID {
	return ID(mid)
}

func TestNewChecker(t *testing.T) {
	a := mockID("a")
	b := mockID("b")
	c := mockID("c")
	d := mockID("d")
	empty := mockID("")

	Convey("new checker", t, func() {
		Convey("when empty", func() {
			checker := NewChecker([]Identifiable{})
			So(checker.Check(), ShouldBeNil)
		})

		Convey("when ok", func() {
			checker := NewChecker([]Identifiable{a, b, c, d})
			So(checker.Check(), ShouldBeNil)
		})

		Convey("when ok with repeated sets", func() {
			checker := NewChecker([]Identifiable{a, b, c, a})
			So(checker.Check(), ShouldBeError, "id `a` already defined")
		})

		Convey("when id set error", func() {
			checker := NewChecker([]Identifiable{a, b, c, empty})

			// Values are not sorted, check the beginning of the error
			err := checker.Check()
			So(err, ShouldNotBeNil)
			So(checker.Check().Error(), ShouldEqual, "id required")
		})
	})
}
