package example

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"encoding/csv"
	"fugologic/builder"
	"fugologic/crisp"
	"fugologic/fuzzy"
	"fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

func createIDVals() (*fuzzy.IDVal, *fuzzy.IDVal, *fuzzy.IDVal) {
	// Input #A
	crispDiff, errCrispDiff := crisp.NewSetN(-2, 2, 40)
	fsDiff, errFsDiff := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
		"--": &fuzzy.StepDown{A: -2, B: -0.5},
		"-":  &fuzzy.Triangular{A: -2, B: -0.5, C: 0},
		"0":  &fuzzy.Triangular{A: -0.5, B: 0, C: 0.5},
		"+":  &fuzzy.Triangular{A: 0, B: 0.5, C: 2},
		"++": &fuzzy.StepUp{A: 0.5, B: 2},
	})
	fvDiff, errFvDiff := fuzzy.NewIDVal("diff/consigne", crispDiff, fsDiff)

	So(errCrispDiff, ShouldBeNil)
	So(errFsDiff, ShouldBeNil)
	So(errFvDiff, ShouldBeNil)

	// Input #B
	crispDt, errCrispDt := crisp.NewSetN(-0.2, 0.2, 40)
	fsDt, errFsDt := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
		"--": &fuzzy.StepDown{A: -0.2, B: -0.1},
		"-":  &fuzzy.Triangular{A: -0.2, B: -0.1, C: 0},
		"0":  &fuzzy.Triangular{A: -0.1, B: 0, C: 0.1},
		"+":  &fuzzy.Triangular{A: 0, B: 0.1, C: 0.2},
		"++": &fuzzy.StepUp{A: 0.1, B: 0.2},
	})
	fvDt, errFvDt := fuzzy.NewIDVal("temp/dt", crispDt, fsDt)

	So(errCrispDt, ShouldBeNil)
	So(errFsDt, ShouldBeNil)
	So(errFvDt, ShouldBeNil)

	// Output #C
	crispForce, errCrispForce := crisp.NewSetN(-4, 4, 80)
	fsForce, errFsForce := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
		"--": &fuzzy.StepDown{A: -4, B: -1},
		"-":  &fuzzy.Triangular{A: -2, B: -1, C: 0},
		"0":  &fuzzy.Triangular{A: -1, B: 0, C: 1},
		"+":  &fuzzy.Triangular{A: 0, B: 1, C: 2},
		"++": &fuzzy.StepUp{A: 1, B: 4},
	})
	fvForce, errFvForce := fuzzy.NewIDVal("force", crispForce, fsForce)

	So(errCrispForce, ShouldBeNil)
	So(errFsForce, ShouldBeNil)
	So(errFvForce, ShouldBeNil)

	return fvDiff, fvDt, fvForce
}

func TestExample(t *testing.T) {
	Convey("example", t, func() {
		fvDiff, fvDt, fvForce := createIDVals()
		bld := builder.NewBuilderMamdani()

		// Rules
		// diff & dt -> force
		//               |  (b) temperature/dt    |
		//               |------------------------|
		//               | -- |  - |  0 |  + | ++ |
		// ---------|----|----|----|----|----|----|
		// (a) diff | -- | ++ | ++ | ++ |  + |  + |
		// consigne |  - | ++ | ++ |  + |  0 |  0 |
		//          |  0 |  + |  + |  0 |  - |  - |
		//          |  + |  0 |  0 |  - | -- | -- |
		//          | ++ |  - |  - | -- | -- | -- |
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

			var result [][3]float64
			for _, diff := range fvDiff.U().Values() {
				for _, dt := range fvDt.U().Values() {
					data, errEval := engine.Evaluate(fuzzy.DataInput{
						fvDiff: diff,
						fvDt:   dt,
					})
					So(errEval, ShouldBeNil)

					force := data[fvForce]
					result = append(result, [3]float64{diff, dt, force})
				}
			}

			So(export(result), ShouldBeNil)
		})
	})
}

func export(values [][3]float64) error {
	f, errCreate := os.Create("./data.csv")
	if errCreate != nil {
		return errCreate
	}

	fltToStr := func(flt float64) string {
		return strings.Replace(fmt.Sprintf("%.3f", flt), ".", ",", 1)
	}

	writer := csv.NewWriter(f)
	var data = [][]string{
		{"diff", "dt", "force"},
	}

	// Convert data into strings
	for _, row := range values {
		data = append(data, []string{
			fltToStr(row[0]),
			fltToStr(row[1]),
			fltToStr(row[2]),
		})
	}

	errWrite := writer.WriteAll(data)
	if errWrite != nil {
		return errWrite
	}

	return nil
}
