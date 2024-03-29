package example

import (
	"testing"

	"github.com/sbiemont/fugologic/builder"
	"github.com/sbiemont/fugologic/crisp"
	"github.com/sbiemont/fugologic/fuzzy"
	"github.com/sbiemont/fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMinimalist(t *testing.T) {
	Convey("mini", t, func() {
		// Input #A
		crispA, _ := crisp.NewSetN(-2, 2, 20)
		fsA, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"-": fuzzy.StepDown{A: -1, B: 1},
			"+": fuzzy.StepUp{A: -1, B: 1},
		})
		fvA, _ := fuzzy.NewIDVal("a", crispA, fsA)

		// Input #B
		crispB, _ := crisp.NewSetN(-0.2, 0.2, 20)
		fsB, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"N": fuzzy.Trapezoid{A: -0.2, B: -0.1, C: 0, D: 0.1},
			"P": fuzzy.Trapezoid{A: -0.1, B: 0, C: 0.1, D: 0.2},
		})
		fvB, _ := fuzzy.NewIDVal("b", crispB, fsB)

		// Output #C
		crispC, _ := crisp.NewSetN(-4, 4, 40)
		fsC, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"##": fuzzy.StepDown{A: -4, B: 2},
			"**": fuzzy.StepUp{A: -2, B: 4},
		})
		fvC, _ := fuzzy.NewIDVal("c", crispC, fsC)

		Convey("with builder", func() {
			// Rules
			// a & b -> c
			//       | b       |
			//       |---------|
			//       | N  | P  |
			// --|---|----|----|
			// a | - | ## | ** |
			//   | + | ** | ## |
			bld := builder.Mamdani().FuzzyLogic()
			bld.If(fvA.Get("-")).And(fvB.Get("N")).Then(fvC.Get("##"))
			bld.If(fvA.Get("-")).And(fvB.Get("P")).Then(fvC.Get("**"))
			bld.If(fvA.Get("+")).And(fvB.Get("N")).Then(fvC.Get("**"))
			bld.If(fvA.Get("+")).And(fvB.Get("P")).Then(fvC.Get("##"))

			// Evaluate
			eng, _ := bld.Engine()
			out, _ := eng.Evaluate(fuzzy.DataInput{
				fvA: -0.2,
				fvB: 0.05,
			})

			So(out, ShouldResemble, fuzzy.DataOutput{
				fvC: 0.17915132672502984,
			})
		})

		Convey("with fam", func() {
			bld := builder.Mamdani().FuzzyAssoMatrix()
			err := bld.
				Asso(fvA, fvB, fvC).
				Matrix(
					[]id.ID{"-", "+"},
					map[id.ID][]id.ID{
						"N": {"##", "**"},
						"P": {"**", "##"},
					})
			So(err, ShouldBeNil)

			eng, _ := bld.Engine()
			out, _ := eng.Evaluate(fuzzy.DataInput{
				fvA: -0.2,
				fvB: 0.05,
			})

			So(out, ShouldResemble, fuzzy.DataOutput{
				fvC: 0.17915132672502984,
			})
		})
	})
}

func TestExample(t *testing.T) {
	Convey("example", t, func() {
		// Input #A
		crispDiff, errCrispDiff := crisp.NewSetN(-2, 2, 40)
		fsDiff, errFsDiff := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"--": fuzzy.StepDown{A: -2, B: -0.5},
			"-":  fuzzy.Triangular{A: -2, B: -0.5, C: 0},
			"0":  fuzzy.Triangular{A: -0.5, B: 0, C: 0.5},
			"+":  fuzzy.Triangular{A: 0, B: 0.5, C: 2},
			"++": fuzzy.StepUp{A: 0.5, B: 2},
		})
		fvDiff, errFvDiff := fuzzy.NewIDVal("diff/consigne", crispDiff, fsDiff)

		So(errCrispDiff, ShouldBeNil)
		So(errFsDiff, ShouldBeNil)
		So(errFvDiff, ShouldBeNil)

		// Input #B
		crispDt, errCrispDt := crisp.NewSetN(-0.2, 0.2, 40)
		fsDt, errFsDt := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"--": fuzzy.StepDown{A: -0.2, B: -0.1},
			"-":  fuzzy.Triangular{A: -0.2, B: -0.1, C: 0},
			"0":  fuzzy.Triangular{A: -0.1, B: 0, C: 0.1},
			"+":  fuzzy.Triangular{A: 0, B: 0.1, C: 0.2},
			"++": fuzzy.StepUp{A: 0.1, B: 0.2},
		})
		fvDt, errFvDt := fuzzy.NewIDVal("temp/dt", crispDt, fsDt)

		So(errCrispDt, ShouldBeNil)
		So(errFsDt, ShouldBeNil)
		So(errFvDt, ShouldBeNil)

		// Output #C
		crispForce, errCrispForce := crisp.NewSetN(-4, 4, 80)
		fsForce, errFsForce := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"--": fuzzy.StepDown{A: -4, B: -1},
			"-":  fuzzy.Triangular{A: -2, B: -1, C: 0},
			"0":  fuzzy.Triangular{A: -1, B: 0, C: 1},
			"+":  fuzzy.Triangular{A: 0, B: 1, C: 2},
			"++": fuzzy.StepUp{A: 1, B: 4},
		})
		fvForce, errFvForce := fuzzy.NewIDVal("force", crispForce, fsForce)

		So(errCrispForce, ShouldBeNil)
		So(errFsForce, ShouldBeNil)
		So(errFvForce, ShouldBeNil)

		// Rules
		// (a) diff & (b) dt -> (c) force
		//
		//             |  (a) diff              |
		//             |------------------------|
		//             | -- |  - |  0 |  + | ++ |
		// -------|----|----|----|----|----|----|
		// (b) dt | -- | ++ | ++ |  + |  0 |  - |
		//        |  - | ++ | ++ |  + |  0 |  - |
		//        |  0 | ++ |  + |  0 |  - | -- |
		//        |  + |  + |  0 |  - | -- | -- |
		//        | ++ |  + |  0 |  - | -- | -- |
		Convey("when using a fuzzy associative matrix", func() {
			// Builder
			bld := builder.Mamdani().FuzzyAssoMatrix()

			// Rules
			errFAM := bld.Asso(fvDiff, fvDt, fvForce).
				Matrix(
					[]id.ID{"--", "-", "0", "+", "++"},
					map[id.ID][]id.ID{
						"--": {"++", "++", "+", "0", "-"},
						"-":  {"++", "++", "+", "0", "-"},
						"0":  {"++", "+", "0", "-", "--"},
						"+":  {"+", "0", "-", "--", "--"},
						"++": {"+", "0", "-", "--", "--"},
					})
			So(errFAM, ShouldBeNil)

			// Build engine
			engine, errEngine := bld.Engine()
			So(errEngine, ShouldBeNil)

			// Evaluate engine
			result, errEval := engine.Evaluate(fuzzy.DataInput{
				fvDiff: 0.05,
				fvDt:   0.3,
			})
			So(errEval, ShouldBeNil)

			// Control output
			So(result, ShouldResemble, fuzzy.DataOutput{
				fvForce: -1.3498557846383932,
			})
		})

		Convey("when using a builder", func() {
			// Builder
			bld := builder.Mamdani().FuzzyLogic()

			// Rules
			bld.If(fvDiff.Get("--")).And(fvDt.Get("--")).Then(fvForce.Get("++"))
			bld.If(fvDiff.Get("--")).And(fvDt.Get("-")).Then(fvForce.Get("++"))
			bld.If(fvDiff.Get("--")).And(fvDt.Get("0")).Then(fvForce.Get("++"))
			bld.If(fvDiff.Get("--")).And(fvDt.Get("+")).Then(fvForce.Get("+"))
			bld.If(fvDiff.Get("--")).And(fvDt.Get("++")).Then(fvForce.Get("+"))

			bld.If(fvDiff.Get("-")).And(fvDt.Get("--")).Then(fvForce.Get("++"))
			bld.If(fvDiff.Get("-")).And(fvDt.Get("-")).Then(fvForce.Get("++"))
			bld.If(fvDiff.Get("-")).And(fvDt.Get("0")).Then(fvForce.Get("+"))
			bld.If(fvDiff.Get("-")).And(fvDt.Get("+")).Then(fvForce.Get("0"))
			bld.If(fvDiff.Get("-")).And(fvDt.Get("++")).Then(fvForce.Get("0"))

			bld.If(fvDiff.Get("0")).And(fvDt.Get("--")).Then(fvForce.Get("+"))
			bld.If(fvDiff.Get("0")).And(fvDt.Get("-")).Then(fvForce.Get("+"))
			bld.If(fvDiff.Get("0")).And(fvDt.Get("0")).Then(fvForce.Get("0"))
			bld.If(fvDiff.Get("0")).And(fvDt.Get("+")).Then(fvForce.Get("-"))
			bld.If(fvDiff.Get("0")).And(fvDt.Get("++")).Then(fvForce.Get("-"))

			bld.If(fvDiff.Get("+")).And(fvDt.Get("--")).Then(fvForce.Get("0"))
			bld.If(fvDiff.Get("+")).And(fvDt.Get("-")).Then(fvForce.Get("0"))
			bld.If(fvDiff.Get("+")).And(fvDt.Get("0")).Then(fvForce.Get("-"))
			bld.If(fvDiff.Get("+")).And(fvDt.Get("+")).Then(fvForce.Get("--"))
			bld.If(fvDiff.Get("+")).And(fvDt.Get("++")).Then(fvForce.Get("--"))

			bld.If(fvDiff.Get("++")).And(fvDt.Get("--")).Then(fvForce.Get("-"))
			bld.If(fvDiff.Get("++")).And(fvDt.Get("-")).Then(fvForce.Get("-"))
			bld.If(fvDiff.Get("++")).And(fvDt.Get("0")).Then(fvForce.Get("--"))
			bld.If(fvDiff.Get("++")).And(fvDt.Get("+")).Then(fvForce.Get("--"))
			bld.If(fvDiff.Get("++")).And(fvDt.Get("++")).Then(fvForce.Get("--"))

			Convey("when evaluate", func() {
				// Build engine
				engine, errEngine := bld.Engine()
				So(errEngine, ShouldBeNil)

				// Evaluate engine
				result, errEval := engine.Evaluate(fuzzy.DataInput{
					fvDiff: 0.05,
					fvDt:   0.3,
				})
				So(errEval, ShouldBeNil)

				// Control output
				So(result, ShouldResemble, fuzzy.DataOutput{
					fvForce: -1.3498557846383932,
				})
			})

			Convey("when full values", func() {
				engine, errEngine := bld.Engine()
				So(errEngine, ShouldBeNil)

				var result [][]float64
				for _, diff := range fvDiff.U().Values() {
					for _, dt := range fvDt.U().Values() {
						data, errEval := engine.Evaluate(fuzzy.DataInput{
							fvDiff: diff,
							fvDt:   dt,
						})
						So(errEval, ShouldBeNil)

						force := data[fvForce]
						result = append(result, []float64{diff, dt, force})
					}
				}

				So(writeCSV("./example_test.csv", []string{"diff", "dt", "force"}, result), ShouldBeNil)
			})
		})
	})
}
