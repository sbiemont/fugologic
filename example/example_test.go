package example

import (
	"testing"

	"fugologic.git/builder"
	"fugologic.git/crisp"
	"fugologic.git/fuzzy"

	. "github.com/smartystreets/goconvey/convey"
)

func customValues() ([5]fuzzy.IDSet, [5]fuzzy.IDSet, [5]fuzzy.IDSet) {
	// Input #A
	setDiff, _ := crisp.NewSetN(-2, 2, 40)
	fvDiff := fuzzy.NewIDValCustom("diff/consigne", setDiff)
	fsDiffM2 := fuzzy.NewIDSetCustom("diff/consigne.--", fuzzy.NewSetStepDown(-2, -0.5), &fvDiff)
	fsDiffM1 := fuzzy.NewIDSetCustom("diff/consigne.-", fuzzy.NewSetTriangular(-2, -0.5, 0), &fvDiff)
	fsDiff0 := fuzzy.NewIDSetCustom("diff/consigne.0", fuzzy.NewSetTriangular(-0.5, 0, 0.5), &fvDiff)
	fsDiffP1 := fuzzy.NewIDSetCustom("diff/consigne.+", fuzzy.NewSetTriangular(0, 0.5, 2), &fvDiff)
	fsDiffP2 := fuzzy.NewIDSetCustom("diff/consigne.++", fuzzy.NewSetStepUp(0.5, 2), &fvDiff)

	// Input #B
	setDt, _ := crisp.NewSetN(-0.2, 0.2, 40)
	fvDt := fuzzy.NewIDValCustom("temp/dt", setDt)
	fsDtM2 := fuzzy.NewIDSetCustom("temp/dt.--", fuzzy.NewSetStepDown(-0.2, -0.1), &fvDt)
	fsDtM1 := fuzzy.NewIDSetCustom("temp/dt.-", fuzzy.NewSetTriangular(-0.2, -0.1, 0), &fvDt)
	fsDt0 := fuzzy.NewIDSetCustom("temp/dt.0", fuzzy.NewSetTriangular(-0.1, 0, 0.1), &fvDt)
	fsDtP1 := fuzzy.NewIDSetCustom("temp/dt.+", fuzzy.NewSetTriangular(0, 0.1, 0.2), &fvDt)
	fsDtP2 := fuzzy.NewIDSetCustom("temp/dt.++", fuzzy.NewSetStepUp(0.1, 0.2), &fvDt)

	// Output #C
	setForce, _ := crisp.NewSetN(-4, 4, 80)
	fvForce := fuzzy.NewIDValCustom("force", setForce)
	fsForceM2 := fuzzy.NewIDSetCustom("force.--", fuzzy.NewSetStepDown(-4, -1), &fvForce)
	fsForceM1 := fuzzy.NewIDSetCustom("force.-", fuzzy.NewSetTriangular(-2, -1, 0), &fvForce)
	fsForce0 := fuzzy.NewIDSetCustom("force.0", fuzzy.NewSetTriangular(-1, 0, 1), &fvForce)
	fsForceP1 := fuzzy.NewIDSetCustom("force.+", fuzzy.NewSetTriangular(0, 1, 2), &fvForce)
	fsForceP2 := fuzzy.NewIDSetCustom("force.++", fuzzy.NewSetStepUp(1, 4), &fvForce)

	return [5]fuzzy.IDSet{fsDiffM2, fsDiffM1, fsDiff0, fsDiffP1, fsDiffP2},
		[5]fuzzy.IDSet{fsDtM2, fsDtM1, fsDt0, fsDtP1, fsDtP2},
		[5]fuzzy.IDSet{fsForceM2, fsForceM1, fsForce0, fsForceP1, fsForceP2}
}

func TestExample(t *testing.T) {
	fsDiff, fsDt, fsForce := customValues()
	bld := builder.NewBuilderMamdani()

	// Rules
	// a & b -> c
	//               |  (b) temperature/dt    |
	//               |------------------------|
	//               | -- |  - |  0 |  + | ++ |
	// ---------|----|----|----|----|----|----|
	// (a) diff | -- | ++ | ++ | ++ |  + |  + |
	// consigne |  - | ++ | ++ |  + |  0 |  0 |
	//          |  0 |  + |  + |  0 |  - |  - |
	//          |  + |  0 |  0 |  - | -- | -- |
	//          | ++ |  - |  - | -- | -- | -- |
	bld.If(fsDiff[0]).And(fsDt[0]).Then(fsForce[4])
	bld.If(fsDiff[0]).And(fsDt[1]).Then(fsForce[4])
	bld.If(fsDiff[0]).And(fsDt[2]).Then(fsForce[4])
	bld.If(fsDiff[0]).And(fsDt[3]).Then(fsForce[3])
	bld.If(fsDiff[0]).And(fsDt[4]).Then(fsForce[3])

	bld.If(fsDiff[1]).And(fsDt[0]).Then(fsForce[4])
	bld.If(fsDiff[1]).And(fsDt[1]).Then(fsForce[4])
	bld.If(fsDiff[1]).And(fsDt[2]).Then(fsForce[3])
	bld.If(fsDiff[1]).And(fsDt[3]).Then(fsForce[2])
	bld.If(fsDiff[1]).And(fsDt[4]).Then(fsForce[2])

	bld.If(fsDiff[2]).And(fsDt[0]).Then(fsForce[3])
	bld.If(fsDiff[2]).And(fsDt[1]).Then(fsForce[3])
	bld.If(fsDiff[2]).And(fsDt[2]).Then(fsForce[2])
	bld.If(fsDiff[2]).And(fsDt[3]).Then(fsForce[1])
	bld.If(fsDiff[2]).And(fsDt[4]).Then(fsForce[1])

	bld.If(fsDiff[3]).And(fsDt[0]).Then(fsForce[2])
	bld.If(fsDiff[3]).And(fsDt[1]).Then(fsForce[2])
	bld.If(fsDiff[3]).And(fsDt[2]).Then(fsForce[1])
	bld.If(fsDiff[3]).And(fsDt[3]).Then(fsForce[0])
	bld.If(fsDiff[3]).And(fsDt[4]).Then(fsForce[0])

	bld.If(fsDiff[4]).And(fsDt[0]).Then(fsForce[1])
	bld.If(fsDiff[4]).And(fsDt[1]).Then(fsForce[1])
	bld.If(fsDiff[4]).And(fsDt[2]).Then(fsForce[0])
	bld.If(fsDiff[4]).And(fsDt[3]).Then(fsForce[0])
	bld.If(fsDiff[4]).And(fsDt[4]).Then(fsForce[0])

	Convey("example", t, func() {
		// Build engine
		engine, errEngine := bld.Engine()
		So(errEngine, ShouldBeNil)

		// Evaluate engine
		result, errEval := engine.Evaluate(fuzzy.DataInput{
			"diff/consigne": 0.05,
			"temp/dt":       0.3,
		})
		So(errEval, ShouldBeNil)

		// Control output
		So(result, ShouldResemble, fuzzy.DataOutput{
			"force": -1.3498557846383932,
		})
	})
}
