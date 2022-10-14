package fuzzy

import (
	"fugologic/crisp"
	"fugologic/id"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Create a fuzzy value, a fuzzy set and link both
func newTestVal(val, set id.ID) (*IDVal, IDSet) {
	fuzzySet := func(x float64) float64 { return x }
	fv, _ := NewIDVal(val, crisp.Set{}, map[id.ID]Set{
		set: fuzzySet,
	})
	return fv, fv.Get(set)
}

func TestDataInput(t *testing.T) {
	fvA, fsA1 := newTestVal("a", "a1")
	_, fsMissingParent := newTestVal("b", "b1")

	Convey("value", t, func() {
		Convey("when parent missing", func() {
			provider := DataInput{}
			fsMissingParent.parent = nil
			y, err := provider.value(fsMissingParent)
			So(err, ShouldBeError, "input: cannot find parent for id set `b1`")
			So(y, ShouldBeZeroValue)
		})

		Convey("when data missing", func() {
			provider := DataInput{}
			y, err := provider.value(fsA1)
			So(err, ShouldBeError, "input: cannot find data for id val `a` (id set `a1`)")
			So(y, ShouldBeZeroValue)
		})

		Convey("when ok", func() {
			provider := DataInput{fvA: 1}
			y, err := provider.value(fsA1)
			So(err, ShouldBeNil)
			So(y, ShouldEqual, 1)
		})
	})
}

func TestMergeData(t *testing.T) {
	fvA, _ := newTestVal("a", "a1")
	fvB, _ := newTestVal("b", "b1")
	fvC, _ := newTestVal("c", "c1")
	fvD, _ := newTestVal("d", "d1")

	empty := dataIO{}
	filled1 := dataIO{
		fvA: 1,
		fvB: 2,
	}
	filled2 := dataIO{
		fvC: 3,
		fvD: 4,
	}

	Convey("merge", t, func() {
		Convey("when empty", func() {
			So(empty.merge(dataIO{}), ShouldBeEmpty)
		})

		Convey("when merge with empty", func() {
			So(empty.merge(filled1), ShouldResemble, filled1)
			So(filled1.merge(empty), ShouldResemble, filled1)
		})

		Convey("when ok", func() {
			So(filled1.merge(filled2), ShouldResemble, dataIO{
				fvA: 1,
				fvB: 2,
				fvC: 3,
				fvD: 4,
			})
			So(filled1, ShouldResemble, dataIO{ // filled1: unchanged
				fvA: 1,
				fvB: 2,
			})
			So(filled2, ShouldResemble, dataIO{ // filled2: unchanged
				fvC: 3,
				fvD: 4,
			})
		})
	})
}
