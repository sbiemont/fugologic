package builder

import (
	"fmt"
	"testing"

	"fugologic/fuzzy"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExpression(t *testing.T) {
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
		bld := NewBuilder(
			fuzzy.ConnectorZadehAnd,
			fuzzy.ConnectorZadehOr,
			nil,
			nil,
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
