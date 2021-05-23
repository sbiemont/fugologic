package fuzzy

import (
	"testing"

	"fugologic.git/crisp"
	"fugologic.git/id"
	. "github.com/smartystreets/goconvey/convey"
)

// Custom engine for testing
// Returns the 2 inputs and the output
func customEngine() (Engine, [2]IDVal, IDVal) {
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
	setCh, _ := crisp.NewSet(-4, 4, 0.05)
	fvCh := NewIDValCustom("force", setCh)
	fsChM2 := NewIDSetCustom("force.--", NewSetStepDown(-4, -1), &fvCh)
	fsChM1 := NewIDSetCustom("force.-", NewSetTriangular(-2, -1, 0), &fvCh)
	fsCh0 := NewIDSetCustom("force.0", NewSetTriangular(-1, 0, 1), &fvCh)
	fsChP1 := NewIDSetCustom("force.+", NewSetTriangular(0, 1, 2), &fvCh)
	fsChP2 := NewIDSetCustom("force.++", NewSetStepUp(1, 4), &fvCh)

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
		NewRule(NewExpression([]Premise{fsDiffM2, fsDtM2}, ConnectorAnd), ImplicationProd, []IDSet{fsChP2}),
		NewRule(NewExpression([]Premise{fsDiffM2, fsDtM1}, ConnectorAnd), ImplicationProd, []IDSet{fsChP2}),
		NewRule(NewExpression([]Premise{fsDiffM2, fsDt0}, ConnectorAnd), ImplicationProd, []IDSet{fsChP2}),
		NewRule(NewExpression([]Premise{fsDiffM2, fsDtP1}, ConnectorAnd), ImplicationProd, []IDSet{fsChP1}),
		NewRule(NewExpression([]Premise{fsDiffM2, fsDtP2}, ConnectorAnd), ImplicationProd, []IDSet{fsChP1}),

		NewRule(NewExpression([]Premise{fsDiffM1, fsDtM2}, ConnectorAnd), ImplicationProd, []IDSet{fsChP2}),
		NewRule(NewExpression([]Premise{fsDiffM1, fsDtM1}, ConnectorAnd), ImplicationProd, []IDSet{fsChP2}),
		NewRule(NewExpression([]Premise{fsDiffM1, fsDt0}, ConnectorAnd), ImplicationProd, []IDSet{fsChP1}),
		NewRule(NewExpression([]Premise{fsDiffM1, fsDtP1}, ConnectorAnd), ImplicationProd, []IDSet{fsCh0}),
		NewRule(NewExpression([]Premise{fsDiffM1, fsDtP2}, ConnectorAnd), ImplicationProd, []IDSet{fsCh0}),

		NewRule(NewExpression([]Premise{fsDiff0, fsDtM2}, ConnectorAnd), ImplicationProd, []IDSet{fsChP1}),
		NewRule(NewExpression([]Premise{fsDiff0, fsDtM1}, ConnectorAnd), ImplicationProd, []IDSet{fsChP1}),
		NewRule(NewExpression([]Premise{fsDiff0, fsDt0}, ConnectorAnd), ImplicationProd, []IDSet{fsCh0}),
		NewRule(NewExpression([]Premise{fsDiff0, fsDtP1}, ConnectorAnd), ImplicationProd, []IDSet{fsChM1}),
		NewRule(NewExpression([]Premise{fsDiff0, fsDtP2}, ConnectorAnd), ImplicationProd, []IDSet{fsChM1}),

		NewRule(NewExpression([]Premise{fsDiffP1, fsDtM2}, ConnectorAnd), ImplicationProd, []IDSet{fsCh0}),
		NewRule(NewExpression([]Premise{fsDiffP1, fsDtM1}, ConnectorAnd), ImplicationProd, []IDSet{fsCh0}),
		NewRule(NewExpression([]Premise{fsDiffP1, fsDt0}, ConnectorAnd), ImplicationProd, []IDSet{fsChM1}),
		NewRule(NewExpression([]Premise{fsDiffP1, fsDtP1}, ConnectorAnd), ImplicationProd, []IDSet{fsChM2}),
		NewRule(NewExpression([]Premise{fsDiffP1, fsDtP2}, ConnectorAnd), ImplicationProd, []IDSet{fsChM2}),

		NewRule(NewExpression([]Premise{fsDiffP2, fsDtM2}, ConnectorAnd), ImplicationProd, []IDSet{fsChM1}),
		NewRule(NewExpression([]Premise{fsDiffP2, fsDtM1}, ConnectorAnd), ImplicationProd, []IDSet{fsChM1}),
		NewRule(NewExpression([]Premise{fsDiffP2, fsDt0}, ConnectorAnd), ImplicationProd, []IDSet{fsChM2}),
		NewRule(NewExpression([]Premise{fsDiffP2, fsDtP1}, ConnectorAnd), ImplicationProd, []IDSet{fsChM2}),
		NewRule(NewExpression([]Premise{fsDiffP2, fsDtP2}, ConnectorAnd), ImplicationProd, []IDSet{fsChM2}),
	}
	defuzzer := NewDefuzzer(DefuzzificationCentroid)
	engine, _ := NewEngine(rules, defuzzer)
	return engine, [2]IDVal{fvDiff, fvDt}, fvCh
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
				NewRule(fsA1, ImplicationProd, []IDSet{fsB1}),
				NewRule(fsA1, ImplicationProd, []IDSet{fsC1}),
				NewRule(fsC1, ImplicationProd, []IDSet{fsD1}),
			}
			defuzzer := NewDefuzzer(DefuzzificationCentroid)
			_, err := NewEngine(rules, defuzzer)
			So(err, ShouldBeNil)
		})

		Convey("when error", func() {
			fvCBis := NewIDValCustom("c'", crispSet)
			fsC1Bis := NewIDSetCustom("c1", fuzzySet, &fvCBis)

			// a => b
			// a => c' [c and c' have the same id]
			// c => d
			rules := []Rule{
				NewRule(fsA1, ImplicationProd, []IDSet{fsB1}),
				NewRule(fsA1, ImplicationProd, []IDSet{fsC1Bis}),
				NewRule(fsC1, ImplicationProd, []IDSet{fsD1}),
			}
			defuzzer := NewDefuzzer(DefuzzificationCentroid)
			_, err := NewEngine(rules, defuzzer)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "sets: id `c1` already present (for val id `c")
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
			NewRule(NewExpression([]Premise{fsA1, fsB1}, ConnectorAnd), ImplicationProd, []IDSet{fsC1}),
			NewRule(NewExpression([]Premise{fsA2, fsB2}, ConnectorAnd), ImplicationProd, []IDSet{fsC2}),
		}
		defuzzer := NewDefuzzer(DefuzzificationCentroid)
		engine, errEngine := NewEngine(rules, defuzzer)
		So(errEngine, ShouldBeNil)

		dataIn := DataInput{
			fvA.ID(): 2.1,
			fvB.ID(): 3.9,
		}

		result, errEval := engine.Evaluate(dataIn)
		So(errEval, ShouldBeNil)

		So(result, ShouldNotBeNil)
		So(result, ShouldResemble, DataOutput{
			fvC.ID(): 12.5,
		})
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
				outputForce: -0.8,
			},
			{
				inputDiff:   -1,
				inputDt:     -0.1,
				outputForce: 3.1,
			},
			{
				inputDiff:   -1,
				inputDt:     0.1,
				outputForce: 0.3,
			},
			{
				inputDiff:   1,
				inputDt:     0.1,
				outputForce: -3.1,
			},
			{
				inputDiff:   1,
				inputDt:     -0.1,
				outputForce: -0.3,
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
				outputForce: -3.1,
			},
		}

		engine, inputs, output := customEngine()

		for _, tt := range tests {
			result, errEngine := engine.Evaluate(map[id.ID]float64{
				inputs[0].ID(): tt.inputDiff,
				inputs[1].ID(): tt.inputDt,
			})
			So(errEngine, ShouldBeNil)
			So(result, ShouldHaveLength, 1)
			So(result[output.ID()], ShouldAlmostEqual, tt.outputForce, 0.1)
		}
	})
}

func BenchmarkEngineNTimes(b *testing.B) {
	const evaluations = 1000
	engine, inputs, _ := customEngine()

	// Evaluate n times the system
	for i := 0; i < evaluations; i++ {
		engine.Evaluate(map[id.ID]float64{
			inputs[0].ID(): -1,
			inputs[1].ID(): -0.1,
		})
	}
}

func BenchmarkEngineMatrix(b *testing.B) {
	engine, inputs, _ := customEngine()

	// Evaluate all possible values of the system
	for _, diff := range inputs[0].u.Values() {
		for _, dt := range inputs[1].u.Values() {
			engine.Evaluate(DataInput{
				inputs[0].ID(): diff,
				inputs[1].ID(): dt,
			})
		}
	}
}
