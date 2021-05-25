package fuzzy

import (
	"testing"

	"fugologic.git/crisp"
	"fugologic.git/id"
	. "github.com/smartystreets/goconvey/convey"
)

func customValues() ([5]IDSet, [5]IDSet, [5]IDSet) {
	// Input #1
	setDiff, _ := crisp.NewSet(-2, 2, 0.1)
	fvDiff := NewIDValCustom("diff/consigne", setDiff)
	fsDiffM2 := NewIDSetCustom("diff/consigne.--", NewSetStepDown(-2, -0.5), &fvDiff)
	fsDiffM1 := NewIDSetCustom("diff/consigne.-", NewSetTriangular(-2, -0.5, 0), &fvDiff)
	fsDiff0 := NewIDSetCustom("diff/consigne.0", NewSetTriangular(-0.5, 0, 0.5), &fvDiff)
	fsDiffP1 := NewIDSetCustom("diff/consigne.+", NewSetTriangular(0, 0.5, 2), &fvDiff)
	fsDiffP2 := NewIDSetCustom("diff/consigne.++", NewSetStepUp(0.5, 2), &fvDiff)

	// Input #2
	setDt, _ := crisp.NewSet(-0.2, 0.2, 0.01)
	fvDt := NewIDValCustom("temp/dt", setDt)
	fsDtM2 := NewIDSetCustom("temp/dt.--", NewSetStepDown(-0.2, -0.1), &fvDt)
	fsDtM1 := NewIDSetCustom("temp/dt.-", NewSetTriangular(-0.2, -0.1, 0), &fvDt)
	fsDt0 := NewIDSetCustom("temp/dt.0", NewSetTriangular(-0.1, 0, 0.1), &fvDt)
	fsDtP1 := NewIDSetCustom("temp/dt.+", NewSetTriangular(0, 0.1, 0.2), &fvDt)
	fsDtP2 := NewIDSetCustom("temp/dt.++", NewSetStepUp(0.1, 0.2), &fvDt)

	// Output
	setCh, _ := crisp.NewSet(-4, 4, 0.1)
	fvCh := NewIDValCustom("force", setCh)
	fsChM2 := NewIDSetCustom("force.--", NewSetStepDown(-4, -1), &fvCh)
	fsChM1 := NewIDSetCustom("force.-", NewSetTriangular(-2, -1, 0), &fvCh)
	fsCh0 := NewIDSetCustom("force.0", NewSetTriangular(-1, 0, 1), &fvCh)
	fsChP1 := NewIDSetCustom("force.+", NewSetTriangular(0, 1, 2), &fvCh)
	fsChP2 := NewIDSetCustom("force.++", NewSetStepUp(1, 4), &fvCh)

	return [5]IDSet{fsDiffM2, fsDiffM1, fsDiff0, fsDiffP1, fsDiffP2},
		[5]IDSet{fsDtM2, fsDtM1, fsDt0, fsDtP1, fsDtP2},
		[5]IDSet{fsChM2, fsChM1, fsCh0, fsChP1, fsChP2}
}

// Custom engine for testing
func customEngine() (Engine, [5]IDSet, [5]IDSet, [5]IDSet) {
	fsDiff, fsDt, fsCh := customValues()

	// Rules
	// a & b -> c
	//               |    temp/dt             |
	//               |------------------------|
	//               | -- |  - |  0 |  + | ++ |
	// ---------|----|----|----|----|----|----|
	// diff     | -- | ++ | ++ | ++ |  + |  + |
	// consigne |  - | ++ | ++ |  + |  0 |  0 |
	//          |  0 |  + |  + |  0 |  - |  - |
	//          |  + |  0 |  0 |  - | -- | -- |
	//          | ++ |  - |  - | -- | -- | -- |
	rules := []Rule{
		NewRule(NewExpression([]Premise{fsDiff[0], fsDt[0]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[4]}),
		NewRule(NewExpression([]Premise{fsDiff[0], fsDt[1]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[4]}),
		NewRule(NewExpression([]Premise{fsDiff[0], fsDt[2]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[4]}),
		NewRule(NewExpression([]Premise{fsDiff[0], fsDt[3]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[3]}),
		NewRule(NewExpression([]Premise{fsDiff[0], fsDt[4]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[3]}),

		NewRule(NewExpression([]Premise{fsDiff[1], fsDt[0]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[4]}),
		NewRule(NewExpression([]Premise{fsDiff[1], fsDt[1]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[4]}),
		NewRule(NewExpression([]Premise{fsDiff[1], fsDt[2]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[3]}),
		NewRule(NewExpression([]Premise{fsDiff[1], fsDt[3]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[2]}),
		NewRule(NewExpression([]Premise{fsDiff[1], fsDt[4]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[2]}),

		NewRule(NewExpression([]Premise{fsDiff[2], fsDt[0]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[3]}),
		NewRule(NewExpression([]Premise{fsDiff[2], fsDt[1]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[3]}),
		NewRule(NewExpression([]Premise{fsDiff[2], fsDt[2]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[2]}),
		NewRule(NewExpression([]Premise{fsDiff[2], fsDt[3]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[1]}),
		NewRule(NewExpression([]Premise{fsDiff[2], fsDt[4]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[1]}),

		NewRule(NewExpression([]Premise{fsDiff[3], fsDt[0]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[2]}),
		NewRule(NewExpression([]Premise{fsDiff[3], fsDt[1]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[2]}),
		NewRule(NewExpression([]Premise{fsDiff[3], fsDt[2]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[1]}),
		NewRule(NewExpression([]Premise{fsDiff[3], fsDt[3]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[0]}),
		NewRule(NewExpression([]Premise{fsDiff[3], fsDt[4]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[0]}),

		NewRule(NewExpression([]Premise{fsDiff[4], fsDt[0]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[1]}),
		NewRule(NewExpression([]Premise{fsDiff[4], fsDt[1]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[1]}),
		NewRule(NewExpression([]Premise{fsDiff[4], fsDt[2]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[0]}),
		NewRule(NewExpression([]Premise{fsDiff[4], fsDt[3]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[0]}),
		NewRule(NewExpression([]Premise{fsDiff[4], fsDt[4]}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsCh[0]}),
	}

	engine, _ := NewEngine(rules, DefuzzificationCentroid)
	return engine, fsDiff, fsDt, fsCh
}

func TestEngineCheck(t *testing.T) {
	Convey("check", t, func() {
		var fuzzySet Set = func(x float64) float64 { return x }
		var crispSet crisp.Set = crisp.Set{}

		fvA := NewIDValCustom("a", crispSet)
		fsA1 := NewIDSetCustom("a1", fuzzySet, &fvA)

		fvB := NewIDValCustom("b", crispSet)
		fsB1 := NewIDSetCustom("b1", fuzzySet, &fvB)

		fvC := NewIDValCustom("c", crispSet)
		fsC1 := NewIDSetCustom("c1", fuzzySet, &fvC)

		fvD := NewIDValCustom("d", crispSet)
		fsD1 := NewIDSetCustom("d1", fuzzySet, &fvD)

		Convey("when ok", func() {
			// a => b
			// a => c
			// c => d
			rules := []Rule{
				NewRule(fsA1, ImplicationMin, []IDSet{fsB1}),
				NewRule(fsA1, ImplicationMin, []IDSet{fsC1}),
				NewRule(fsC1, ImplicationMin, []IDSet{fsD1}),
			}
			_, err := NewEngine(rules, DefuzzificationCentroid)
			So(err, ShouldBeNil)
		})

		Convey("when error", func() {
			fvCBis := NewIDValCustom("c'", crispSet)
			fsC1Bis := NewIDSetCustom("c1", fuzzySet, &fvCBis)

			// a => b
			// a => c' [c and c' have the same id]
			// c => d
			rules := []Rule{
				NewRule(fsA1, ImplicationMin, []IDSet{fsB1}),
				NewRule(fsA1, ImplicationMin, []IDSet{fsC1Bis}),
				NewRule(fsC1, ImplicationMin, []IDSet{fsD1}),
			}
			_, err := NewEngine(rules, DefuzzificationCentroid)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "sets: id `c1` already present (for val id `c")
		})
	})
}

func TestEvaluate(t *testing.T) {
	Convey("custom minimalistic test", t, func() {
		// Definitions
		setA, _ := crisp.NewSet(1, 4, 0.1)
		fvA := NewIDValCustom("a", setA)
		fsA1 := NewIDSetCustom("a1", NewSetTriangular(1, 2, 3), &fvA)
		fsA2 := NewIDSetCustom("a2", NewSetTriangular(2, 3, 4), &fvA)

		setB, _ := crisp.NewSet(2, 5, 0.1)
		fvB := NewIDValCustom("b", setB)
		fsB1 := NewIDSetCustom("b1", NewSetTriangular(2, 3, 4), &fvB)
		fsB2 := NewIDSetCustom("b2", NewSetTriangular(3, 4, 5), &fvB)

		setC, _ := crisp.NewSet(11, 14, 0.1)
		fvC := NewIDValCustom("c", setC)
		fsC1 := NewIDSetCustom("c1", NewSetTriangular(11, 12, 13), &fvC)
		fsC2 := NewIDSetCustom("c2", NewSetTriangular(12, 13, 14), &fvC)

		// Rules
		// a1 & b1 -> c1
		// a2 & b2 -> c2
		rules := []Rule{
			NewRule(NewExpression([]Premise{fsA1, fsB1}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsC1}),
			NewRule(NewExpression([]Premise{fsA2, fsB2}, ConnectorZadehAnd), ImplicationMin, []IDSet{fsC2}),
		}

		engine, errEngine := NewEngine(rules, DefuzzificationCentroid)
		So(errEngine, ShouldBeNil)

		dataIn := DataInput{
			fvA.ID(): 2.1,
			fvB.ID(): 3.9,
		}

		// Evaluate engine
		result, errEval := engine.Evaluate(dataIn)
		So(errEval, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(result[fvC.ID()], ShouldEqual, 12.5)
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
		check := func(engine Engine, fsDiff, fsDt, fsCh [5]IDSet) {
			for _, tt := range tests {
				result, errEngine := engine.Evaluate(map[id.ID]float64{
					fsDiff[0].parent.ID(): tt.inputDiff,
					fsDt[0].parent.ID():   tt.inputDt,
				})
				So(errEngine, ShouldBeNil)
				So(result, ShouldHaveLength, 1)
				So(result[fsCh[0].parent.ID()], ShouldAlmostEqual, tt.outputForce, 0.01)
			}
		}

		check(customEngine())
	})
}

func BenchmarkEngineNTimes(b *testing.B) {
	const evaluations = 1000
	engine, fsDiff, fsDt, _ := customEngine()

	// Evaluate n times the system
	for i := 0; i < evaluations; i++ {
		engine.Evaluate(map[id.ID]float64{
			fsDiff[0].parent.ID(): -1,
			fsDt[0].parent.ID():   -0.1,
		})
	}
}

func BenchmarkEngineMatrix(b *testing.B) {
	engine, fsDiff, fsDt, _ := customEngine()

	// Evaluate all possible values of the system
	diffValues := fsDiff[0].parent.u.Values()
	dtValues := fsDt[0].parent.u.Values()
	for _, diff := range diffValues {
		for _, dt := range dtValues {
			engine.Evaluate(DataInput{
				fsDiff[0].parent.ID(): diff,
				fsDt[0].parent.ID():   dt,
			})
		}
	}
}
