package fuzzy

import (
	"math"
	"testing"

	"fugologic/crisp"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExpression(t *testing.T) {
	var timesTwo Set = func(x float64) float64 { return x * 2 }

	fvA := NewIDValCustom("a", crisp.Set{})
	fsA1 := NewIDSetCustom("a1", timesTwo, &fvA)

	fvB := NewIDValCustom("b", crisp.Set{})
	fsB1 := NewIDSetCustom("b1", timesTwo, &fvB)

	fvC := NewIDValCustom("c", crisp.Set{})
	fsC1 := NewIDSetCustom("c1", timesTwo, &fvC)

	fvD := NewIDValCustom("d", crisp.Set{})
	fsD1 := NewIDSetCustom("d1", timesTwo, &fvD)

	fvE := NewIDValCustom("e", crisp.Set{})
	fsE1 := NewIDSetCustom("e1", timesTwo, &fvE)

	Convey("new expression", t, func() {
		Convey("when empty", func() {
			exp := NewExpression([]Premise{}, ConnectorZadehAnd)
			result, err := exp.Evaluate(DataInput{})
			So(err, ShouldBeError, "expression: at least 1 premise expected")
			So(result, ShouldBeZeroValue)
		})

		Convey("when one premise", func() {
			dataIn := DataInput{
				fvA.uuid: 1,
			}
			result, err := fsA1.Evaluate(dataIn)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 1*2) // iso(1)*2
		})

		Convey("when several premises", func() {
			dataIn := DataInput{
				fvA.uuid: 1,
				fvB.uuid: 2,
				fvC.uuid: 3,
			}

			Convey("when connector AND", func() {
				exp := NewExpression([]Premise{fsA1, fsB1, fsC1}, ConnectorZadehAnd)
				result, err := exp.Evaluate(dataIn)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, 1*2) // min(1, 2, 3)*2
			})

			Convey("when connector OR", func() {
				exp := NewExpression([]Premise{fsA1, fsB1, fsC1}, ConnectorZadehOr)
				result, err := exp.Evaluate(dataIn)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, 3*2) // max(1, 2, 3)*2
			})
		})

		Convey("when complex expression (A and B and C) or (D and E)", func() {
			dataIn := DataInput{
				fvA.uuid: 1,
				fvB.uuid: 2,
				fvC.uuid: 3,
				fvD.uuid: 4,
				fvE.uuid: 5,
			}

			expABC := NewExpression([]Premise{fsA1, fsB1, fsC1}, ConnectorZadehAnd)
			expDE := NewExpression([]Premise{fsD1, fsE1}, ConnectorZadehAnd)
			exp := NewExpression([]Premise{expABC, expDE}, ConnectorZadehOr)

			result, err := exp.Evaluate(dataIn)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 8) // max(min(1, 2, 3)*2, min(4, 5)*2)
		})
	})

	Convey("connect", t, func() {
		Convey("when 2 premises", func() {
			exp := NewExpression([]Premise{fsA1}, nil).Connect(fsB1, math.Max)
			result, err := exp.Evaluate(DataInput{
				"a": 1,
				"b": 2,
			})
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 4) // 2*max(1, 2)
		})

		Convey("when 1 expression", func() {
			exp := NewExpression([]Premise{fsA1, fsB1}, nil).Connect(fsC1, math.Max)
			result, err := exp.Evaluate(DataInput{
				"a": 1,
				"b": 2,
				"c": 3,
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
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
				"e": 5,
			})
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 8) // 2*min(max(min(max(1, 2), 3), 4), 5)
		})
	})
}
