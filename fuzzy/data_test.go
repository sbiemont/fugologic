package fuzzy

import (
	"testing"

	"fugologic.git/crisp"
	"fugologic.git/id"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDataInput(t *testing.T) {
	newSet := func() Set { return func(x float64) float64 { return x } }

	fvA := NewIDValCustom("a", crisp.Set{})
	fsA1 := NewIDSetCustom("a1", newSet(), &fvA)

	fvMissingParent := NewIDValCustom("b", crisp.Set{})
	fsMissingParent := NewIDSetCustom("b1", newSet(), &fvMissingParent)

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
			provider := DataInput{"a": 1}
			y, err := provider.find(fsA1)
			So(err, ShouldBeNil)
			So(y, ShouldEqual, 1)
		})
	})
}

func TestMergeData(t *testing.T) {
	empty := Data{}
	filled1 := Data{
		"a": 1,
		"b": 2,
	}
	filled2 := Data{
		"c": 3,
		"d": 4,
	}

	Convey("merge", t, func() {
		Convey("when empty", func() {
			So(mergeData(empty, Data{}), ShouldBeEmpty)
		})

		Convey("when merge with empty", func() {
			So(mergeData(empty, filled1), ShouldResemble, map[id.ID]float64(filled1))
			So(mergeData(filled1, empty), ShouldResemble, map[id.ID]float64(filled1))
		})

		Convey("when ok", func() {
			So(mergeData(filled1, filled2), ShouldResemble, map[id.ID]float64{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			})
			So(filled1, ShouldResemble, Data{ // filled1: unchanged
				"a": 1,
				"b": 2,
			})
			So(filled2, ShouldResemble, Data{ // filled2: unchanged
				"c": 3,
				"d": 4,
			})
		})
	})
}
