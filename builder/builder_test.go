package builder

import (
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
			bld := NewBuilder(
				fuzzy.ConnectorZadehAnd,
				fuzzy.ConnectorZadehOr,
				nil,
				nil,
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
			So(res, ShouldEqual, 4)
		})

		Convey("when connectors hyberbolic", func() {
			// (A and B and C) or (D and E)
			bld := NewBuilder(
				fuzzy.ConnectorHyperbolicAnd,
				fuzzy.ConnectorHyperbolicOr,
				nil,
				nil,
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
	})
}

func TestAdd(t *testing.T) {
	Convey("explicit add rule", t, func() {
		bld := NewBuilder(nil, nil, nil, nil, nil, nil, nil)
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

		bld := NewBuilderMamdani()
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
	fvA,fsA1 := newTestVal("a", "a1")
	fvB,fsB1 := newTestVal("b", "b1")
	fvC,fsC1 := newTestVal("c", "c1")

	Convey("engine", t, func() {
		bld := NewBuilderMamdani()
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
