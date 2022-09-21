package fuzzy

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDataInput(t *testing.T) {
	fvA, fsA1 := newTestVal("a", "a1")
	_, fsMissingParent := newTestVal("b", "b1")

	Convey("find", t, func() {
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
			provider := DataInput{fvA: 1}
			y, err := provider.find(fsA1)
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

	empty := map[*IDVal]float64{}
	filled1 := map[*IDVal]float64{
		fvA: 1,
		fvB: 2,
	}
	filled2 := map[*IDVal]float64{
		fvC: 3,
		fvD: 4,
	}

	Convey("merge", t, func() {
		Convey("when empty", func() {
			So(mergeData(empty, map[*IDVal]float64{}), ShouldBeEmpty)
		})

		Convey("when merge with empty", func() {
			So(mergeData(empty, filled1), ShouldResemble, map[*IDVal]float64(filled1))
			So(mergeData(filled1, empty), ShouldResemble, map[*IDVal]float64(filled1))
		})

		Convey("when ok", func() {
			So(mergeData(filled1, filled2), ShouldResemble, map[*IDVal]float64{
				fvA: 1,
				fvB: 2,
				fvC: 3,
				fvD: 4,
			})
			So(filled1, ShouldResemble, map[*IDVal]float64{ // filled1: unchanged
				fvA: 1,
				fvB: 2,
			})
			So(filled2, ShouldResemble, map[*IDVal]float64{ // filled2: unchanged
				fvC: 3,
				fvD: 4,
			})
		})
	})
}
