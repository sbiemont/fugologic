package builder

import (
	"fmt"
	"testing"

	"fugologic.git/fuzzy"
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
		builder := Builder{
			fuzzy.ConnectorZadehAnd,
			fuzzy.ConnectorZadehOr,
			fuzzy.ImplicationMin,
		}

		expAB := Expression{
			builder: builder,
			fzExp:   fuzzy.NewExpression([]fuzzy.Premise{fsA1, fsB1}, fuzzy.ConnectorZadehAnd),
		}
		expCD := Expression{
			builder: builder,
			fzExp:   fuzzy.NewExpression([]fuzzy.Premise{fsC1, fsD1}, fuzzy.ConnectorZadehAnd),
		}

		Convey("and", func() {
			exp := expAB.And(expCD)

			res, err := exp.Evaluate(input)
			So(err, ShouldBeNil)
			So(res, ShouldEqual, 1)

			Convey("then", func() {
				engine, _ := fuzzy.NewEngine(
					[]fuzzy.Rule{exp.Then([]fuzzy.IDSet{fsE1})}, // only checks the "then" call
					fuzzy.DefuzzificationCentroid,
				)
				res, err := engine.Evaluate(input)
				So(err, ShouldBeNil)
				fmt.Print(res)
				So(res, ShouldEqual, fuzzy.DataOutput{
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
				engine, _ := fuzzy.NewEngine(
					[]fuzzy.Rule{exp.Then([]fuzzy.IDSet{fsE1})}, // only checks the "then" call
					fuzzy.DefuzzificationCentroid,
				)
				res, err := engine.Evaluate(input)
				So(err, ShouldBeNil)
				fmt.Print(res)
				So(res, ShouldEqual, fuzzy.DataOutput{
					"e": 0,
				})
			})
		})
	})
}
