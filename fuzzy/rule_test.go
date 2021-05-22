package fuzzy

import (
	"testing"

	"fugologic.git/crisp"
	"fugologic.git/id"
	. "github.com/smartystreets/goconvey/convey"
)

func TestImplication(t *testing.T) {
	var plusOne Set = func(x float64) float64 {
		return x + 1
	}

	Convey("implication product", t, func() {
		So(ImplicationProd(plusOne, 0)(2), ShouldEqual, 0)  // (2+1)*0
		So(ImplicationProd(plusOne, 1)(2), ShouldEqual, 3)  // (2+1)*1
		So(ImplicationProd(plusOne, 5)(2), ShouldEqual, 15) // (2+1)*5
	})

	Convey("implication min", t, func() {
		So(ImplicationMin(plusOne, 0)(2), ShouldEqual, 0) // min(2+1, 0)
		So(ImplicationMin(plusOne, 1)(2), ShouldEqual, 1) // min(2+1, 1)
		So(ImplicationMin(plusOne, 5)(2), ShouldEqual, 3) // min(2+1, 5)
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
			exp := NewExpression([]Premise{}, ConnectorAnd)
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
				exp := NewExpression([]Premise{fsA1, fsB1, fsC1}, ConnectorAnd)
				result, err := exp.Evaluate(dataIn)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, 1*2) // min(1, 2, 3)*2
			})

			Convey("when connector OR", func() {
				exp := NewExpression([]Premise{fsA1, fsB1, fsC1}, ConnectorOr)
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

			expABC := NewExpression([]Premise{fsA1, fsB1, fsC1}, ConnectorAnd)
			expDE := NewExpression([]Premise{fsD1, fsE1}, ConnectorAnd)
			exp := NewExpression([]Premise{expABC, expDE}, ConnectorOr)

			result, err := exp.Evaluate(dataIn)
			So(err, ShouldBeNil)
			So(result, ShouldEqual, 8) // max(min(1, 2, 3)*2, min(4, 5)*2)
		})
	})
}

func TestFlattenIDSets(t *testing.T) {
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

	Convey("when only one set", t, func() {
		result := flattenIDSets(nil, []Premise{fsA1})
		So(result, ShouldHaveLength, 1)
		So(result[0].ID(), ShouldEqual, "a1")
	})

	Convey("flatten id sets", t, func() {
		expAB := NewExpression([]Premise{fsA1, fsB1}, ConnectorAnd)       // A and B
		expCD := NewExpression([]Premise{fsC1, fsD1}, ConnectorAnd)       // C and D
		expABCD := NewExpression([]Premise{expAB, expCD}, ConnectorOr)    // (A and B) or (C and D)
		expABCDE := NewExpression([]Premise{expABCD, fsE1}, ConnectorAnd) // ((A and B) or (C and D)) and E

		result := flattenIDSets(nil, []Premise{expABCDE})
		So(result, ShouldHaveLength, 5)

		// Extract data
		var parents []*IDVal
		var uuids []id.ID
		for _, idSet := range result {
			parents = append(parents, idSet.parent)
			uuids = append(uuids, idSet.uuid)
		}
		So(parents, ShouldResemble, []*IDVal{&fvA, &fvB, &fvC, &fvD, &fvE})
		So(uuids, ShouldResemble, []id.ID{"a1", "b1", "c1", "d1", "e1"})
	})
}

func TestRule(t *testing.T) {
	// helper: extract ids
	ids := func(idSets []IDSet) []id.ID {
		var result []id.ID
		for _, idSet := range idSets {
			result = append(result, idSet.uuid)
		}
		return result
	}

	Convey("evaluate", t, func() {
		var setA Set = func(x float64) float64 { return x }
		var setB Set = func(x float64) float64 { return x }

		fvA := NewIDValCustom("a", crisp.Set{})
		fsA1 := NewIDSetCustom("a1", setA, &fvA)

		fvB := NewIDValCustom("b", crisp.Set{})
		fsB1 := NewIDSetCustom("b1", setB, &fvB)

		// A => B
		rule := NewRule(fsA1, ImplicationProd, []IDSet{fsB1})

		Convey("when empty data", func() {
			dataIn := DataInput{}
			output, err := rule.evaluate(dataIn)
			So(err, ShouldBeError, "input: cannot find data for id val `a` (id set `a1`)")
			So(output, ShouldBeEmpty)
		})

		Convey("when ok", func() {
			dataIn := DataInput{
				"a": 1,
			}
			output, err := rule.evaluate(dataIn)
			So(err, ShouldBeNil)
			So(output, ShouldHaveLength, 1)
			So(output[0].parent, ShouldEqual, &fvB)
			So(output[0].set, ShouldNotEqual, setB) // Membership function should have been replaced
			So(output[0].uuid, ShouldEqual, "b1")
		})
	})

	Convey("inputs", t, func() {
		var set Set = func(x float64) float64 { return x }

		fvA := NewIDValCustom("a", crisp.Set{})
		fsA1 := NewIDSetCustom("a1", set, &fvA)

		fvB := NewIDValCustom("b", crisp.Set{})
		fsB1 := NewIDSetCustom("b1", set, &fvB)

		fvC := NewIDValCustom("c", crisp.Set{})
		fsC1 := NewIDSetCustom("c1", set, &fvC)

		fvD := NewIDValCustom("d", crisp.Set{})
		fsD1 := NewIDSetCustom("d1", set, &fvD)

		Convey("when one input", func() {
			// A => B
			rule := NewRule(fsA1, ImplicationProd, []IDSet{fsB1})
			So(ids(rule.Inputs()), ShouldResemble, []id.ID{fsA1.ID()})
		})

		Convey("when several inputs", func() {
			// (A and B) or C => D
			expAB := NewExpression([]Premise{fsA1, fsB1}, ConnectorAnd)
			expABC := NewExpression([]Premise{expAB, fsC1}, ConnectorAnd)
			rule := NewRule(expABC, ImplicationProd, []IDSet{fsD1})
			So(ids(rule.Inputs()), ShouldResemble, []id.ID{fsA1.ID(), fsB1.ID(), fsC1.ID()})
		})
	})
}
