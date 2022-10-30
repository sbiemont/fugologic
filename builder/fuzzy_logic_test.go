package builder

import (
	"fmt"
	"testing"

	"fugologic/crisp"
	"fugologic/fuzzy"
	"fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

// Create a fuzzy value, a fuzzy set and link both
func newTestVal(val, set id.ID) (*fuzzy.IDVal, fuzzy.IDSet) {
	fuzzySet := func(x float64) float64 { return x }
	fv, _ := fuzzy.NewIDVal(val, crisp.Set{}, map[id.ID]fuzzy.Set{
		set: fuzzySet,
	})
	return fv, fv.Get(set)
}

func TestIf(t *testing.T) {
	fvA, fsA1 := newTestVal("a", "a1")
	fvB, fsB1 := newTestVal("b", "b1")
	fvC, fsC1 := newTestVal("c", "c1")
	fvD, fsD1 := newTestVal("d", "d1")
	fvE, fsE1 := newTestVal("e", "e1")

	Convey("if", t, func() {
		Convey("when zadeh connectors", func() {
			bld := NewFuzzyLogic(
				fuzzy.OperatorZadeh,
				nil,
				nil,
				nil,
			)

			// (A and B and C) or (D and E)
			expABC := bld.If(fsA1).And(fsB1).And(fsC1)
			expCD := bld.If(fsD1).And(fsE1)
			rule := bld.If(expABC.Or(expCD))

			res, err := rule.Evaluate(fuzzy.DataInput{
				fvA: 1,
				fvB: 2,
				fvC: 3,
				fvD: 4,
				fvE: 5,
			})
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 4) // max(min(1,2,3), min(4,5))
		})

		Convey("when connectors hyberbolic", func() {
			// (A and B and C) or (D and E)
			bld := NewFuzzyLogic(
				fuzzy.Operator{
					And: fuzzy.OperatorHyperbolic.And,
					Or:  fuzzy.OperatorHyperbolic.Or,
				},
				nil,
				nil,
				nil,
			)

			expABC := bld.If(fsA1).And(fsB1).And(fsC1)
			expCD := bld.If(fsD1).And(fsE1)
			exp := bld.If(expABC.Or(expCD))

			res, err := exp.Evaluate(fuzzy.DataInput{
				fvA: 1,
				fvB: 2,
				fvC: 3,
				fvD: 4,
				fvE: 5,
			})

			abc := 1 * 2 * 3
			de := 4 * 5
			expected := abc + de - abc*de
			So(err, ShouldBeNil)
			So(res, ShouldEqual, expected)
		})

		Convey("when not-and", func() {
			bld := NewFuzzyLogic(
				fuzzy.OperatorZadeh,
				nil,
				nil,
				nil,
			)

			// A not-and B
			rule := bld.If(fsA1).NAnd(fsB1)
			res, err := rule.Evaluate(fuzzy.DataInput{
				fvA: 10,
				fvB: 20,
			})
			So(err, ShouldBeNil)
			So(res, ShouldEqual, -9) // 1-min(10,20)
		})

		Convey("when not-or", func() {
			bld := NewFuzzyLogic(
				fuzzy.OperatorZadeh,
				nil,
				nil,
				nil,
			)

			// A not-and B
			rule := bld.If(fsA1).NOr(fsB1)
			res, err := rule.Evaluate(fuzzy.DataInput{
				fvA: 10,
				fvB: 20,
			})
			So(err, ShouldBeNil)
			So(res, ShouldEqual, -19) // 1-max(10,20)
		})

		Convey("when x-or", func() {
			bld := NewFuzzyLogic(
				fuzzy.OperatorZadeh,
				nil,
				nil,
				nil,
			)

			// A x-or B
			rule := bld.If(fsA1).XOr(fsB1)
			res, err := rule.Evaluate(fuzzy.DataInput{
				fvA: 10,
				fvB: 20,
			})
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 10) // 10+20-2*min(10,20)

		})
	})
}

func TestAdd(t *testing.T) {
	Convey("explicit add rule", t, func() {
		bld := NewFuzzyLogic(fuzzy.Operator{}, nil, nil, nil)
		So(bld.rules, ShouldBeEmpty)

		// Add rule #1
		rule1 := fuzzy.Rule{}
		bld.add(rule1)
		So(bld.rules, ShouldResemble, []fuzzy.Rule{rule1})

		// Add rule #2
		rule2 := fuzzy.Rule{}
		bld.add(rule2)
		So(bld.rules, ShouldResemble, []fuzzy.Rule{rule1, rule2})
	})

	Convey("implicit add rule", t, func() {
		_, fsA1 := newTestVal("a", "a1")
		_, fsB1 := newTestVal("b", "b1")
		_, fsC1 := newTestVal("c", "c1")

		bld := Mamdani().FuzzyLogic()
		So(bld.rules, ShouldBeEmpty)

		// Add rule #1
		bld.If(fsA1).Then(fsC1)
		So(bld.rules, ShouldHaveLength, 1)

		// Add rule #2
		bld.If(fsB1).Then(fsC1)
		So(bld.rules, ShouldHaveLength, 2)
	})
}

func TestEngine(t *testing.T) {
	fvA, fsA1 := newTestVal("a", "a1")
	fvB, fsB1 := newTestVal("b", "b1")
	fvC, fsC1 := newTestVal("c", "c1")

	Convey("engine", t, func() {
		bld := Mamdani().FuzzyLogic()
		bld.If(fsA1).Then(fsC1)
		bld.If(fsB1).Then(fsC1)

		engine, err := bld.Engine()
		So(engine, ShouldNotBeEmpty)
		So(err, ShouldBeNil)

		result, err := engine.Evaluate(fuzzy.DataInput{
			fvA: 1,
			fvB: 2,
		})
		So(result, ShouldResemble, fuzzy.DataOutput{
			fvC: 0,
		})
		So(err, ShouldBeNil)
	})
}

func TestFlExpression(t *testing.T) {
	fvA, fsA1 := newTestVal("a", "a1")
	fvB, fsB1 := newTestVal("b", "b1")
	fvC, fsC1 := newTestVal("c", "c1")
	fvD, fsD1 := newTestVal("d", "d1")
	fvE, fsE1 := newTestVal("e", "e1")

	input := fuzzy.DataInput{
		fvA: 1,
		fvB: 2,
		fvC: 3,
		fvD: 4,
		fvE: 5,
	}

	Convey("evaluate", t, func() {
		bld := NewFuzzyLogic(
			fuzzy.Operator{
				And: fuzzy.OperatorZadeh.And,
				Or:  fuzzy.OperatorZadeh.Or,
			},
			fuzzy.ImplicationMin,
			fuzzy.AggregationUnion,
			fuzzy.DefuzzificationCentroid,
		)

		expAB := flExpression{
			fl:    &bld,
			fzExp: fuzzy.NewExpression([]fuzzy.Premise{fsA1, fsB1}, fuzzy.OperatorZadeh.And),
		}
		expCD := flExpression{
			fl:    &bld,
			fzExp: fuzzy.NewExpression([]fuzzy.Premise{fsC1, fsD1}, fuzzy.OperatorZadeh.And),
		}

		Convey("and", func() {
			exp := expAB.And(expCD)

			res, err := exp.Evaluate(input)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 1)

			Convey("then", func() {
				exp.Then(fsE1) // only checks the "then" call
				engine, _ := bld.Engine()
				res, err := engine.Evaluate(input)
				So(err, ShouldBeNil)
				fmt.Print(res)
				So(res, ShouldResemble, fuzzy.DataOutput{
					fvE: 0,
				})
			})
		})

		Convey("or", func() {
			exp := expAB.Or(expCD)
			res, err := exp.Evaluate(input)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 3)

			Convey("then", func() {
				exp.Then(fsE1) // only checks the "then" call
				engine, _ := bld.Engine()
				res, err := engine.Evaluate(input)
				So(err, ShouldBeNil)
				fmt.Print(res)
				So(res, ShouldResemble, fuzzy.DataOutput{
					fvE: 0,
				})
			})
		})
	})
}
