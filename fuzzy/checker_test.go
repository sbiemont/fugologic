package fuzzy

import (
	"testing"

	"fugologic.git/crisp"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewChecker(t *testing.T) {
	fvA := NewIDValCustom("a", crisp.Set{})
	fsA1 := NewIDSetCustom("a1", nil, &fvA)
	fsA2 := NewIDSetCustom("a2", nil, &fvA)

	fvB := NewIDValCustom("b", crisp.Set{})
	fsB1 := NewIDSetCustom("b1", nil, &fvB)
	fsB2 := NewIDSetCustom("b2", nil, &fvB)

	Convey("new checker", t, func() {
		Convey("when empty", func() {
			checker := newChecker([]IDSet{})
			So(checker.check(), ShouldBeNil)
		})

		Convey("when ok", func() {
			checker := newChecker([]IDSet{fsA1, fsA2, fsB1, fsB2})
			So(checker.check(), ShouldBeNil)
		})

		Convey("when ok with repeated sets", func() {
			checker := newChecker([]IDSet{fsA1, fsA1, fsA1, fsA1})
			So(checker.check(), ShouldBeNil)
		})

		Convey("when id set error", func() {
			fvC := NewIDValCustom("c", crisp.Set{})
			fsC1 := NewIDSetCustom("a1", nil, &fvC)

			checker := newChecker([]IDSet{fsA1, fsA2, fsB1, fsB2, fsC1})
			err := checker.check()
			// Values are not sorted, check the beginning of the error
			So(err, ShouldNotBeNil)
			So(checker.check().Error(), ShouldContainSubstring, "sets: id `a1` already present (for val id `")
		})

		Convey("when id val error", func() {
			fvC := NewIDValCustom("a", crisp.Set{}) // value already defined
			fsC1 := NewIDSetCustom("c1", nil, &fvC)

			checker := newChecker([]IDSet{fsA1, fsC1})
			So(checker.check(), ShouldBeError, "values: id `a` already present")
		})

		Convey("when id val missing", func() {
			fvC := NewIDValCustom("", crisp.Set{})
			fsC1 := NewIDSetCustom("c1", nil, &fvC)

			checker := newChecker([]IDSet{fsC1})
			So(checker.check(), ShouldBeError, "values: id required")
		})

		Convey("when id set missing", func() {
			fvC := NewIDValCustom("c", crisp.Set{})
			fsC1 := NewIDSetCustom("", nil, &fvC)

			checker := newChecker([]IDSet{fsC1})
			So(checker.check(), ShouldBeError, "sets: id required (for val id `c`)")
		})

		Convey("when id set parent error", func() {
			fvParent := NewIDValCustom("c", crisp.Set{})
			fsMissingParent := NewIDSetCustom("missing#1", nil, &fvParent)
			fsMissingParent.parent = nil
			checker := newChecker([]IDSet{fsMissingParent})
			So(checker.check(), ShouldBeError, "sets: no parent found for id set `missing#1`")
		})
	})
}
