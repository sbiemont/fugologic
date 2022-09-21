package builder

import (
	"fmt"
	"testing"

	"fugologic/fuzzy"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExpression(t *testing.T) {
	fsA1 := newTestSet("a")
	fsB1 := newTestSet("b")
	fsC1 := newTestSet("c")
	fsD1 := newTestSet("d")
	fsE1 := newTestSet("e")

	input := fuzzy.DataInput{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
		"e": 5,
	}

	Convey("evaluate", t, func() {
		bld := NewBuilder(
			fuzzy.ConnectorZadehAnd,
			fuzzy.ConnectorZadehOr,
			fuzzy.ImplicationMin,
			fuzzy.AggregationUnion,
			fuzzy.DefuzzificationCentroid,
		)

		expAB := expression{
			bld:   &bld,
			fzExp: fuzzy.NewExpression([]fuzzy.Premise{fsA1, fsB1}, fuzzy.ConnectorZadehAnd),
		}
		expCD := expression{
			bld:   &bld,
			fzExp: fuzzy.NewExpression([]fuzzy.Premise{fsC1, fsD1}, fuzzy.ConnectorZadehAnd),
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
					"e": 0,
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
					"e": 0,
				})
			})
		})
	})
}
