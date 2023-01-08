package fuzzy

import (
	"math"
	"testing"

	"github.com/sbiemont/fugologic/crisp"
	"github.com/sbiemont/fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExpression(t *testing.T) {
	var timesTwo Set = func(x float64) float64 { return x * 2 }

	fvA, _ := NewIDVal("a", crisp.Set{}, map[id.ID]Set{"a1": timesTwo})
	fsA1 := fvA.Get("a1")

	fvB, _ := NewIDVal("b", crisp.Set{}, map[id.ID]Set{"b1": timesTwo})
	fsB1 := fvB.Get("b1")

	fvC, _ := NewIDVal("c", crisp.Set{}, map[id.ID]Set{"c1": timesTwo})
	fsC1 := fvC.Get("c1")

	fvD, _ := NewIDVal("d", crisp.Set{}, map[id.ID]Set{"d1": timesTwo})
	fsD1 := fvD.Get("d1")

	fvE, _ := NewIDVal("e", crisp.Set{}, map[id.ID]Set{"e1": timesTwo})
	fsE1 := fvE.Get("e1")

	Convey("new expression", t, func() {
		Convey("when empty", func() {
			exp := NewExpression([]Premise{}, OperatorZadeh{}.And)
			result, err := exp.Evaluate(DataInput{})
			So(err, ShouldBeError, "expression: at least 1 premise expected")
			So(result, ShouldBeZeroValue)
		})

		Convey("when one premise", func() {
			dataIn := DataInput{
				fvA: 1,
			}
			result, err := fsA1.Evaluate(dataIn)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 1*2) // iso(1)*2
		})

		Convey("when complement", func() {
			exp := NewExpression([]Premise{fsA1}, nil).Not()
			result, err := exp.Evaluate(DataInput{
				fvA: 42,
			})
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 1-42*2)

			Convey("when complement not-not", func() {
				exp := NewExpression([]Premise{fsA1}, nil).Not().Not()
				result, err := exp.Evaluate(DataInput{
					fvA: 42,
				})
				So(err, ShouldBeNil)
				So(result, ShouldEqual, 42*2)
			})
		})

		Convey("when several premises", func() {
			dataIn := DataInput{
				fvA: 1,
				fvB: 2,
				fvC: 3,
			}

			Convey("when connector AND", func() {
				exp := NewExpression([]Premise{fsA1, fsB1, fsC1}, OperatorZadeh{}.And)
				result, err := exp.Evaluate(dataIn)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, 1*2) // min(1, 2, 3)*2

				Convey("when connector NOT-AND", func() {
					result, err := exp.Not().Evaluate(dataIn)
					So(err, ShouldBeNil)
					So(result, ShouldEqual, 1-1*2) // 1-min(1, 2, 3)*2
				})
			})

			Convey("when connector OR", func() {
				exp := NewExpression([]Premise{fsA1, fsB1, fsC1}, OperatorZadeh{}.Or)
				result, err := exp.Evaluate(dataIn)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, 3*2) // max(1, 2, 3)*2

				Convey("when connector NOT-OR", func() {
					result, err := exp.Not().Evaluate(dataIn)
					So(err, ShouldBeNil)
					So(result, ShouldEqual, 1-3*2) // 1-max(1, 2, 3)*2
				})
			})
		})

		Convey("when complex expression (A and B and C) or (D and E)", func() {
			dataIn := DataInput{
				fvA: 1,
				fvB: 2,
				fvC: 3,
				fvD: 4,
				fvE: 5,
			}

			expABC := NewExpression([]Premise{fsA1, fsB1, fsC1}, OperatorZadeh{}.And)
			expDE := NewExpression([]Premise{fsD1, fsE1}, OperatorZadeh{}.And)
			exp := NewExpression([]Premise{expABC, expDE}, OperatorZadeh{}.Or)

			result, err := exp.Evaluate(dataIn)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 8) // max(min(1, 2, 3)*2, min(4, 5)*2)
		})

		Convey("when complex expression complemented : (A and B) not-and C", func() {
			dataIn := DataInput{
				fvA: 1,
				fvB: 2,
				fvC: 3,
			}

			expAB := NewExpression([]Premise{fsA1, fsB1}, OperatorZadeh{}.And)
			exp := NewExpression([]Premise{expAB, fsC1}, OperatorZadeh{}.And).Not()

			result, err := exp.Evaluate(dataIn)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 1-2) // 1 - min(min(1, 2)*2, 3*2)
		})

		Convey("when complex expression complemented : A and (B not-and C)", func() {
			dataIn := DataInput{
				fvA: 1,
				fvB: 2,
				fvC: 3,
			}

			expBC := NewExpression([]Premise{fsB1, fsC1}, OperatorZadeh{}.And).Not()
			exp := NewExpression([]Premise{fsA1, expBC}, OperatorZadeh{}.And)

			result, err := exp.Evaluate(dataIn)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 1-4) // min(1*2, 1 - min(2, 3)*2)
		})

		Convey("when id-set fails", func() {
			fsKo := IDSet{uuid: "ko"}
			exp := NewExpression([]Premise{fsKo}, nil)
			res, err := exp.Evaluate(DataInput{})
			So(err, ShouldBeError, "input: cannot find parent for id set `ko`")
			So(res, ShouldBeZeroValue)
		})

	})

	Convey("connect", t, func() {
		Convey("when 2 premises", func() {
			exp := NewExpression([]Premise{fsA1}, nil).Connect(fsB1, math.Max)
			result, err := exp.Evaluate(DataInput{
				fvA: 1,
				fvB: 2,
			})
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 4) // 2*max(1, 2)
		})

		Convey("when 1 expression", func() {
			exp := NewExpression([]Premise{fsA1, fsB1}, nil).Connect(fsC1, math.Max)
			result, err := exp.Evaluate(DataInput{
				fvA: 1,
				fvB: 2,
				fvC: 3,
			})
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 6) // 2*max(1, 2, 3)
		})

		Convey("when several connect", func() {
			exp := NewExpression([]Premise{fsA1}, nil).
				Connect(fsB1, math.Max).
				Connect(fsC1, math.Min).
				Connect(fsD1, math.Max).
				Connect(fsE1, math.Min)
			result, err := exp.Evaluate(DataInput{
				fvA: 1,
				fvB: 2,
				fvC: 3,
				fvD: 4,
				fvE: 5,
			})
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 8) // 2*min(max(min(max(1, 2), 3), 4), 5)
		})
	})
}
