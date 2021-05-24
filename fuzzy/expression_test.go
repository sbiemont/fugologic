package fuzzy

import (
	"testing"

	"fugologic.git/crisp"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPremise(t *testing.T) {
	Convey("if", t, func() {
		Convey("when fuzzy set", func() {
			fv := NewIDVal(crisp.Set{})
			fs := NewIDSet(NewSetTriangular(0, 1, 2), &fv)

			premise := If(fs).premise
			So(premise, ShouldHaveSameTypeAs, IDSet{})
			So(premise.(IDSet).uuid, ShouldEqual, fs.uuid)
		})

		Convey("when expression", func() {
			fv := NewIDVal(crisp.Set{})
			fs := NewIDSet(NewSetTriangular(0, 1, 2), &fv)

			expression := NewExpression([]Premise{fs}, nil)
			premise := If(expression).premise
			So(premise, ShouldHaveSameTypeAs, Expression{})
		})
	})
}

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
}
