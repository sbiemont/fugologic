package fuzzy

import (
	"testing"

	"fugologic/crisp"
	"fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

func customValues() (*IDVal, *IDVal, *IDVal) {
	// Input #1
	setDiff, _ := crisp.NewSet(-2, 2, 0.1)
	fzDiff, _ := NewIDVal("diff/consigne", setDiff, map[id.ID]Set{
		"--": NewSetStepDown(-2, -0.5),
		"-":  NewSetTriangular(-2, -0.5, 0),
		"0":  NewSetTriangular(-0.5, 0, 0.5),
		"+":  NewSetTriangular(0, 0.5, 2),
		"++": NewSetStepUp(0.5, 2),
	})

	// Input #2
	setDt, _ := crisp.NewSet(-0.2, 0.2, 0.01)
	fzDt, _ := NewIDVal("temp/dt", setDt, map[id.ID]Set{
		"--": NewSetStepDown(-0.2, -0.1),
		"-":  NewSetTriangular(-0.2, -0.1, 0),
		"0":  NewSetTriangular(-0.1, 0, 0.1),
		"+":  NewSetTriangular(0, 0.1, 0.2),
		"++": NewSetStepUp(0.1, 0.2),
	})

	// Output
	setCh, _ := crisp.NewSet(-4, 4, 0.1)
	fzCh, _ := NewIDVal("force", setCh, map[id.ID]Set{
		"--": NewSetStepDown(-4, -1),
		"-":  NewSetTriangular(-2, -1, 0),
		"0":  NewSetTriangular(-1, 0, 1),
		"+":  NewSetTriangular(0, 1, 2),
		"++": NewSetStepUp(1, 4),
	})

	return fzDiff, fzDt, fzCh
}

// Custom engine for testing
func customEngine() (Engine, *IDVal, *IDVal, *IDVal) {
	fvDiff, fvDt, fvCh := customValues()

	// Rules
	// diff & dt -> force
	//               |    temp/dt             |
	//               |------------------------|
	//               | -- |  - |  0 |  + | ++ |
	// ---------|----|----|----|----|----|----|
	// diff     | -- | ++ | ++ | ++ |  + |  + |
	// consigne |  - | ++ | ++ |  + |  0 |  0 |
	//          |  0 |  + |  + |  0 |  - |  - |
	//          |  + |  0 |  0 |  - | -- | -- |
	//          | ++ |  - |  - | -- | -- | -- |
	and := ConnectorZadehAnd
	implies := ImplicationMin
	rules := []Rule{
		NewRule(NewExpression([]Premise{fvDiff.Get("--"), fvDt.Get("--")}, and), implies, []IDSet{fvCh.Get("++")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("--"), fvDt.Get("-")}, and), implies, []IDSet{fvCh.Get("++")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("--"), fvDt.Get("0")}, and), implies, []IDSet{fvCh.Get("++")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("--"), fvDt.Get("+")}, and), implies, []IDSet{fvCh.Get("+")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("--"), fvDt.Get("++")}, and), implies, []IDSet{fvCh.Get("+")}),

		NewRule(NewExpression([]Premise{fvDiff.Get("-"), fvDt.Get("--")}, and), implies, []IDSet{fvCh.Get("++")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("-"), fvDt.Get("-")}, and), implies, []IDSet{fvCh.Get("++")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("-"), fvDt.Get("0")}, and), implies, []IDSet{fvCh.Get("+")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("-"), fvDt.Get("+")}, and), implies, []IDSet{fvCh.Get("0")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("-"), fvDt.Get("++")}, and), implies, []IDSet{fvCh.Get("0")}),

		NewRule(NewExpression([]Premise{fvDiff.Get("0"), fvDt.Get("--")}, and), implies, []IDSet{fvCh.Get("+")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("0"), fvDt.Get("-")}, and), implies, []IDSet{fvCh.Get("+")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("0"), fvDt.Get("0")}, and), implies, []IDSet{fvCh.Get("0")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("0"), fvDt.Get("+")}, and), implies, []IDSet{fvCh.Get("-")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("0"), fvDt.Get("++")}, and), implies, []IDSet{fvCh.Get("-")}),

		NewRule(NewExpression([]Premise{fvDiff.Get("+"), fvDt.Get("--")}, and), implies, []IDSet{fvCh.Get("0")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("+"), fvDt.Get("-")}, and), implies, []IDSet{fvCh.Get("0")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("+"), fvDt.Get("0")}, and), implies, []IDSet{fvCh.Get("-")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("+"), fvDt.Get("+")}, and), implies, []IDSet{fvCh.Get("--")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("+"), fvDt.Get("++")}, and), implies, []IDSet{fvCh.Get("--")}),

		NewRule(NewExpression([]Premise{fvDiff.Get("++"), fvDt.Get("--")}, and), implies, []IDSet{fvCh.Get("-")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("++"), fvDt.Get("-")}, and), implies, []IDSet{fvCh.Get("-")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("++"), fvDt.Get("0")}, and), implies, []IDSet{fvCh.Get("--")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("++"), fvDt.Get("+")}, and), implies, []IDSet{fvCh.Get("--")}),
		NewRule(NewExpression([]Premise{fvDiff.Get("++"), fvDt.Get("++")}, and), implies, []IDSet{fvCh.Get("--")}),
	}

	engine, err := NewEngine(rules, AggregationUnion, DefuzzificationCentroid)
	So(err, ShouldBeNil)
	return engine, fvDiff, fvDt, fvCh
}

func TestEngineCheck(t *testing.T) {
	Convey("check", t, func() {
		_, fsA1 := newTestVal("a", "a1")
		_, fsB1 := newTestVal("b", "b1")
		_, fsC1 := newTestVal("c", "c1")
		_, fsD1 := newTestVal("d", "d1")

		Convey("when ok", func() {
			// a => b
			// a => c
			// c => d
			rules := []Rule{
				NewRule(fsA1, ImplicationMin, []IDSet{fsB1}),
				NewRule(fsA1, ImplicationMin, []IDSet{fsC1}),
				NewRule(fsC1, ImplicationMin, []IDSet{fsD1}),
			}
			_, err := NewEngine(rules, AggregationUnion, DefuzzificationCentroid)
			So(err, ShouldBeNil)
		})

		Convey("when error", func() {
			fvCBis, _ := newTestVal("c", "c1")

			// a => b
			// a => c' [c and c' have the same id]
			// c => d
			rules := []Rule{
				NewRule(fsA1, ImplicationMin, []IDSet{fsB1}),
				NewRule(fsA1, ImplicationMin, []IDSet{fvCBis.Get("c1")}),
				NewRule(fsC1, ImplicationMin, []IDSet{fsD1}),
			}
			_, err := NewEngine(rules, AggregationUnion, DefuzzificationCentroid)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "values: id `c` already defined")
		})
	})
}

func TestEvaluate(t *testing.T) {
	Convey("custom minimalistic test", t, func() {
		// Definitions
		setA, _ := crisp.NewSet(1, 4, 0.1)
		fvA, _ := NewIDVal("a", setA, map[id.ID]Set{
			"a1": NewSetTriangular(1, 2, 3),
			"a2": NewSetTriangular(2, 3, 4),
		})

		setB, _ := crisp.NewSet(2, 5, 0.1)
		fvB, _ := NewIDVal("b", setB, map[id.ID]Set{
			"b1": NewSetTriangular(2, 3, 4),
			"b2": NewSetTriangular(3, 4, 5),
		})

		setC, _ := crisp.NewSet(11, 14, 0.1)
		fvC, _ := NewIDVal("c", setC, map[id.ID]Set{
			"c1": NewSetTriangular(11, 12, 13),
			"c2": NewSetTriangular(12, 13, 14),
		})

		// Rules
		// a1 & b1 -> c1
		// a2 & b2 -> c2
		rules := []Rule{
			NewRule(NewExpression([]Premise{fvA.Get("a1"), fvB.Get("b1")}, ConnectorZadehAnd), ImplicationMin, []IDSet{fvC.Get("c1")}),
			NewRule(NewExpression([]Premise{fvA.Get("a2"), fvB.Get("b2")}, ConnectorZadehAnd), ImplicationMin, []IDSet{fvC.Get("c2")}),
		}

		engine, errEngine := NewEngine(rules, AggregationUnion, DefuzzificationCentroid)
		So(errEngine, ShouldBeNil)

		dataIn := DataInput{
			fvA: 2.1,
			fvB: 3.9,
		}

		// Evaluate engine
		result, errEval := engine.Evaluate(dataIn)
		So(errEval, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(result[fvC], ShouldEqual, 12.5)
	})

	Convey("custom evaluate", t, func() {
		type test struct {
			inputDiff   float64
			inputDt     float64
			outputForce float64
		}

		tests := []test{
			{
				inputDiff:   -0.1,
				inputDt:     0.1,
				outputForce: -0.75,
			},
			{
				inputDiff:   -1,
				inputDt:     -0.1,
				outputForce: 2.94,
			},
			{
				inputDiff:   -1,
				inputDt:     0.1,
				outputForce: 0.36,
			},
			{
				inputDiff:   1,
				inputDt:     0.1,
				outputForce: -2.94,
			},
			{
				inputDiff:   1,
				inputDt:     -0.1,
				outputForce: -0.36,
			},
			{
				inputDiff:   0,
				inputDt:     0,
				outputForce: 0,
			},
			{
				inputDiff:   0,
				inputDt:     2,
				outputForce: -1,
			},
			{
				inputDiff:   20,
				inputDt:     0.1,
				outputForce: -3.03,
			},
		}

		// Check engine method
		check := func(engine Engine, fvDiff, fvDt, fvCh *IDVal) {
			for _, tt := range tests {
				result, errEngine := engine.Evaluate(map[*IDVal]float64{
					fvDiff: tt.inputDiff,
					fvDt:   tt.inputDt,
				})
				So(errEngine, ShouldBeNil)
				So(result, ShouldHaveLength, 1)
				So(result[fvCh], ShouldAlmostEqual, tt.outputForce, 0.01)
			}
		}

		check(customEngine())
	})
}

func BenchmarkEngineNTimes(b *testing.B) {
	const evaluations = 1000
	engine, fvDiff, fvDt, _ := customEngine()

	// Evaluate n times the system
	for i := 0; i < evaluations; i++ {
		engine.Evaluate(map[*IDVal]float64{
			fvDiff: -1,
			fvDt:   -0.1,
		})
	}
}

func BenchmarkEngineMatrix(b *testing.B) {
	engine, fvDiff, fvDt, _ := customEngine()

	// Evaluate all possible values of the system
	diffValues := fvDiff.u.Values()
	dtValues := fvDt.u.Values()
	for _, diff := range diffValues {
		for _, dt := range dtValues {
			engine.Evaluate(DataInput{
				fvDiff: diff,
				fvDt:   dt,
			})
		}
	}
}

func TestCheckIDs(t *testing.T) {
	fvA, _ := NewIDVal("a", crisp.Set{}, map[id.ID]Set{"a1": nil, "a2": nil})
	fsA1 := fvA.Get("a1")
	fsA2 := fvA.Get("a2")

	_, fsABis1 := newTestVal("a", "a-bis-1")
	_, fsATer1 := newTestVal("a-ter", "a1")

	Convey("id val", t, func() {
		Convey("when ok", func() {
			So(checkIDs([]IDSet{fsA1, fsA2}), ShouldBeNil)
		})

		Convey("when same id val uuid", func() {
			So(checkIDs([]IDSet{fsA1, fsA2, fsABis1}), ShouldBeError, "values: id `a` already defined")
		})

		Convey("when same id set uuid", func() {
			So(checkIDs([]IDSet{fsA1, fsA2, fsATer1}), ShouldBeNil)
		})

		Convey("when no id val uuid", func() {
			fv, err := NewIDVal("", crisp.Set{}, map[id.ID]Set{"b1": nil})
			So(err, ShouldBeError, "id val cannot be empty")
			So(fv, ShouldBeNil)
		})

		Convey("when no id set uuid", func() {
			fv, err := NewIDVal("c", crisp.Set{}, map[id.ID]Set{"": nil})
			So(err, ShouldBeError, "id set cannot be empty")
			So(fv, ShouldBeNil)
		})
	})
}
