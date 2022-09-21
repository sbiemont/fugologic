package fuzzy

import (
	"testing"

	"fugologic/crisp"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewIDSet(t *testing.T) {
	iso := func(x float64) float64 { return x }

	Convey("id set", t, func() {
		Convey("when not", func() {
			fvA := NewIDVal(crisp.Set{})
			fsA1 := NewIDSet(iso, &fvA)

			not := fsA1.Not()
			So(not.parent, ShouldEqual, &fvA)
			So(not.uuid, ShouldEqual, fsA1.uuid)
			So(not.set(0.1), ShouldEqual, 0.9)
		})
	})
}

func TestNewIDVal(t *testing.T) {
	Convey("id val", t, func() {
		fvA := NewIDVal(crisp.Set{})
		fsA1 := NewIDSet(nil, &fvA)

		So(fsA1.parent, ShouldEqual, &fvA)
	})
}

func TestCheckIDs(t *testing.T) {
	fvA := NewIDValCustom("a", crisp.Set{})
	fsA1 := NewIDSetCustom("a1", nil, &fvA)
	fsA2 := NewIDSetCustom("a2", nil, &fvA)

	fvB := NewIDValCustom("", crisp.Set{})
	fsB1 := NewIDSetCustom("b1", nil, &fvB)

	fvC := NewIDValCustom("c", crisp.Set{})
	fsC1 := NewIDSetCustom("", nil, &fvC)

	fvABis := NewIDValCustom("a", crisp.Set{})
	fsABis1 := NewIDSetCustom("a-bis-1", nil, &fvABis)

	fvATer := NewIDValCustom("a-ter", crisp.Set{})
	fsATer1 := NewIDSetCustom("a1", nil, &fvATer)

	Convey("id val", t, func() {
		Convey("when ok", func() {
			So(checkIDs([]IDSet{fsA1, fsA2}), ShouldBeNil)
		})

		Convey("when same id val uuid", func() {
			So(checkIDs([]IDSet{fsA1, fsA2, fsABis1}), ShouldBeError, "values: id `a` already defined")
		})

		Convey("when same id set uuid", func() {
			So(checkIDs([]IDSet{fsA1, fsA2, fsATer1}), ShouldBeError, "sets: id `a1` already defined")
		})

		Convey("when no id val uuid", func() {
			So(checkIDs([]IDSet{fsA1, fsA2, fsB1}), ShouldBeError, "values: id required")
		})

		Convey("when no id set uuid", func() {
			So(checkIDs([]IDSet{fsA1, fsA2, fsC1}), ShouldBeError, "sets: id required")
		})
	})
}
