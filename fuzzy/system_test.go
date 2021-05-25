package fuzzy

import (
	"testing"

	"fugologic.git/crisp"
	"fugologic.git/id"
	. "github.com/smartystreets/goconvey/convey"
)

// Create a fuzzy value, a fuzzy set and link both
func newTestSet(name id.ID) IDSet {
	fuzzySet := func(x float64) float64 { return x }
	fv := NewIDValCustom(name, crisp.Set{})
	fs1 := NewIDSetCustom(name+"1", fuzzySet, &fv)
	return fs1
}

func TestSystem(t *testing.T) {
	fsA1 := newTestSet("a")
	fsB1 := newTestSet("b")
	fsC1 := newTestSet("c")
	fsD1 := newTestSet("d")
	fsE1 := newTestSet("e")
	fsF1 := newTestSet("f")
	fsG1 := newTestSet("g")

	defuzz := NewDefuzzer(defuzzificationNone)

	// A and B => C
	rulesEng1 := []Rule{
		NewRule(NewExpression([]Premise{fsA1, fsB1}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsC1}),
	}
	// D => E, F
	rulesEng2 := []Rule{
		NewRule(fsD1, ImplicationMin, []IDSet{fsE1, fsF1}),
	}
	// C and E => G
	rulesEng3 := []Rule{
		NewRule(NewExpression([]Premise{fsC1, fsE1}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsG1}),
	}

	Convey("evaluate", t, func() {
		eng1, err1 := NewEngine(rulesEng1, defuzz)
		So(err1, ShouldBeNil)
		eng2, err2 := NewEngine(rulesEng2, defuzz)
		So(err2, ShouldBeNil)
		eng3, err3 := NewEngine(rulesEng3, defuzz)
		So(err3, ShouldBeNil)

		Convey("when ok", func() {
			var system System = []Engine{eng1, eng2, eng3}
			output, errOut := system.Evaluate(DataInput{
				"a": 1,
				"b": 1,
				"d": 1,
			})
			So(errOut, ShouldBeNil)
			So(output, ShouldResemble, DataOutput{
				"c": 0,
				"e": 0,
				"f": 0,
				"g": 0,
			})
		})

		Convey("when missing input", func() {
			var system System = []Engine{eng1, eng2, eng3}
			output, errOut := system.Evaluate(DataInput{
				"a": 1,
				"b": 1,
			})
			So(errOut, ShouldBeError, "input: cannot find data for id val `d` (id set `d1`)")
			So(output, ShouldBeEmpty)
		})

		Convey("check", func() {
			Convey("duplicated outputs", func() {
				Convey("when output defined twice", func() {
					// D => E, F, G
					rulesEng2Bis := []Rule{NewRule(fsD1, ImplicationProd, []IDSet{fsE1, fsF1, fsG1})}
					eng2Bis, err2Bis := NewEngine(rulesEng2Bis, defuzz)
					So(err2Bis, ShouldBeNil)

					var system System = []Engine{eng1, eng2Bis, eng3}
					So(system.checkDuplicatedOutputs(), ShouldBeError, "output `g` detected twice")
				})

				Convey("when ok", func() {
					var system System = []Engine{eng1, eng2, eng3}
					So(system.checkDuplicatedOutputs(), ShouldBeNil)
				})
			})

			Convey("duplicated id", func() {
				Convey("when duplicated", func() {
					// D => E, A
					fsF1Bis := newTestSet("a")
					rulesEng2Bis := []Rule{NewRule(fsD1, ImplicationProd, []IDSet{fsE1, fsF1Bis})}
					eng2Bis, err2Bis := NewEngine(rulesEng2Bis, defuzz)
					So(err2Bis, ShouldBeNil)

					var system System = []Engine{eng1, eng2Bis, eng3}
					So(system.checkUnicity(), ShouldBeError, "values: id `a` already present")
				})

				Convey("when ok", func() {
					var system System = []Engine{eng1, eng2, eng3}
					So(system.checkUnicity(), ShouldBeNil)
				})
			})

			Convey("loop", func() {
				Convey("when loop", func() {
					// G => E, F
					rulesEng2Bis := []Rule{NewRule(fsG1, ImplicationProd, []IDSet{fsE1, fsF1})}
					eng2Bis, err2Bis := NewEngine(rulesEng2Bis, defuzz)
					So(err2Bis, ShouldBeNil)

					var system System = []Engine{eng1, eng2Bis, eng3}
					So(system.checkLoop(), ShouldBeError, "input `g` cannot become an output")
				})

				Convey("when ok", func() {
					var system System = []Engine{eng1, eng2, eng3}
					So(system.checkLoop(), ShouldBeNil)
				})
			})
		})
	})
}
