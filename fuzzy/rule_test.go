package fuzzy

import (
	"testing"

	"fugologic/crisp"
	"fugologic/id"

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

func TestFlattenIDSets(t *testing.T) {
	var timesTwo Set = func(x float64) float64 { return x * 2 }

	fvA, _ := NewIDVal("a", crisp.Set{}, map[id.ID]Set{"a1": timesTwo})
	fsA1 := fvA.Get("a1")

	fvB, _ := NewIDVal("b", crisp.Set{}, map[id.ID]Set{"b1": timesTwo})
	fsB1 := fvB.Get("b1")

	fvC, _ := NewIDVal("c", crisp.Set{}, map[id.ID]Set{"c1": timesTwo})
	fsC1 := fvC.Get("c1")

	fvD, _ := NewIDVal("d", crisp.Set{}, map[id.ID]Set{"d1": timesTwo})
	fsD1 := fvD.Get("d1")

	fvE, _ := NewIDVal("e", crisp.Set{}, map[id.ID]Set{"e1": timesTwo})
	fsE1 := fvE.Get("e1")

	Convey("when only one set", t, func() {
		result := flattenIDSets(nil, []Premise{fsA1})
		So(result, ShouldHaveLength, 1)
		So(result[0].ID(), ShouldEqual, "a1")
	})

	Convey("flatten id sets", t, func() {
		expAB := NewExpression([]Premise{fsA1, fsB1}, OperatorZadeh.And)       // A and B
		expCD := NewExpression([]Premise{fsC1, fsD1}, OperatorZadeh.And)       // C and D
		expABCD := NewExpression([]Premise{expAB, expCD}, OperatorZadeh.Or)    // (A and B) or (C and D)
		expABCDE := NewExpression([]Premise{expABCD, fsE1}, OperatorZadeh.And) // ((A and B) or (C and D)) and E

		result := flattenIDSets(nil, []Premise{expABCDE})
		So(result, ShouldHaveLength, 5)

		// Extract data
		var parents []*IDVal
		var uuids []id.ID
		for _, idSet := range result {
			parents = append(parents, idSet.parent)
			uuids = append(uuids, idSet.uuid)
		}
		So(parents, ShouldResemble, []*IDVal{fvA, fvB, fvC, fvD, fvE})
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

		fvA, _ := NewIDVal("a", crisp.Set{}, map[id.ID]Set{"a1": setA})
		fsA1 := fvA.Get("a1")

		fvB, _ := NewIDVal("b", crisp.Set{}, map[id.ID]Set{"b1": setB})
		fsB1 := fvB.Get("b1")

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
				fvA: 1,
			}
			output, err := rule.evaluate(dataIn)
			So(err, ShouldBeNil)
			So(output, ShouldHaveLength, 1)
			So(output[0].parent, ShouldEqual, fvB)
			So(output[0].set, ShouldNotEqual, setB) // Membership function should have been replaced
			So(output[0].uuid, ShouldEqual, "b1")
		})
	})

	Convey("inputs", t, func() {
		_, fsA1 := newTestVal("a", "a1")
		_, fsB1 := newTestVal("b", "b1")
		_, fsC1 := newTestVal("c", "c1")
		_, fsD1 := newTestVal("d", "d1")

		Convey("when one input", func() {
			// A => B
			rule := NewRule(fsA1, ImplicationProd, []IDSet{fsB1})
			So(ids(rule.Inputs()), ShouldResemble, []id.ID{fsA1.ID()})
		})

		Convey("when several inputs", func() {
			// (A and B) or C => D
			expAB := NewExpression([]Premise{fsA1, fsB1}, OperatorZadeh.And)
			expABC := NewExpression([]Premise{expAB, fsC1}, OperatorZadeh.And)
			rule := NewRule(expABC, ImplicationProd, []IDSet{fsD1})
			So(ids(rule.Inputs()), ShouldResemble, []id.ID{fsA1.ID(), fsB1.ID(), fsC1.ID()})
		})
	})
}
