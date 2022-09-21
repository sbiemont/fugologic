package builder

import (
	"testing"

	"fugologic/crisp"
	"fugologic/fuzzy"
	"fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

// Create a fuzzy value, a fuzzy set and link both
func newTestSet(name id.ID) fuzzy.IDSet {
	fuzzySet := func(x float64) float64 { return x }
	fv := fuzzy.NewIDValCustom(name, crisp.Set{})
	fs1 := fuzzy.NewIDSetCustom(name+"1", fuzzySet, &fv)
	return fs1
}

func TestIf(t *testing.T) {
	fsA1 := newTestSet("a")
	fsB1 := newTestSet("b")
	fsC1 := newTestSet("c")
	fsD1 := newTestSet("d")
	fsE1 := newTestSet("e")

	Convey("if", t, func() {
		Convey("when zadeh connectors", func() {
			bld := NewBuilder(
				fuzzy.ConnectorZadehAnd,
				fuzzy.ConnectorZadehOr,
				nil,
				nil,
				nil,
			)

			// (A and B and C) or (D and E)
			expABC := bld.If(fsA1).And(fsB1).And(fsC1)
			expCD := bld.If(fsD1).And(fsE1)
			rule := bld.If(expABC.Or(expCD))

			res, err := rule.Evaluate(fuzzy.DataInput{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
				"e": 5,
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
			)

			expABC := bld.If(fsA1).And(fsB1).And(fsC1)
			expCD := bld.If(fsD1).And(fsE1)
			exp := bld.If(expABC.Or(expCD))

			res, err := exp.Evaluate(fuzzy.DataInput{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
				"e": 5,
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
		bld := NewBuilder(nil, nil, nil, nil, nil)
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
		fsA1 := newTestSet("a")
		fsB1 := newTestSet("b")
		fsC1 := newTestSet("c")

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
	fsA1 := newTestSet("a")
	fsB1 := newTestSet("b")
	fsC1 := newTestSet("c")

	Convey("engine", t, func() {
		bld := NewBuilderMamdani()
		bld.If(fsA1).Then(fsC1)
		bld.If(fsB1).Then(fsC1)

		engine, err := bld.Engine()
		So(engine, ShouldNotBeEmpty)
		So(err, ShouldBeNil)

		result, err := engine.Evaluate(fuzzy.DataInput{
			"a": 1,
			"b": 2,
		})
		So(result, ShouldResemble, fuzzy.DataOutput{
			"c": 0,
		})
		So(err, ShouldBeNil)
	})
}
