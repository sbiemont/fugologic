package builder

import (
	"testing"

	"fugologic.git/crisp"
	"fugologic.git/fuzzy"
	"fugologic.git/id"

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
			builder := Builder{
				fuzzy.ConnectorZadehAnd,
				fuzzy.ConnectorZadehOr,
				nil,
			}

			// (A and B and C) or (D and E)
			expABC := builder.If(fsA1).And(fsB1).And(fsC1)
			expCD := builder.If(fsD1).And(fsE1)
			rule := builder.If(expABC.Or(expCD))

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
			builder := Builder{
				fuzzy.ConnectorHyperbolicAnd,
				fuzzy.ConnectorHyperbolicOr,
				nil,
			}

			expABC := builder.If(fsA1).And(fsB1).And(fsC1)
			expCD := builder.If(fsD1).And(fsE1)
			rule := builder.If(expABC.Or(expCD))

			res, err := rule.Evaluate(fuzzy.DataInput{
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
