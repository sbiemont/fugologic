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

func customValues() (*fuzzy.IDVal, *fuzzy.IDVal, *fuzzy.IDVal) {
	// Input #A
	crispDiff, _ := crisp.NewSetN(-2, 2, 40)
	fzDiff, _ := fuzzy.NewIDVal("diff/consigne", crispDiff, map[id.ID]fuzzy.Set{
		"--": fuzzy.NewSetStepDown(-2, -0.5),
		"-":  fuzzy.NewSetTriangular(-2, -0.5, 0),
		"0":  fuzzy.NewSetTriangular(-0.5, 0, 0.5),
		"+":  fuzzy.NewSetTriangular(0, 0.5, 2),
		"++": fuzzy.NewSetStepUp(0.5, 2),
	})

	// Input #B
	crispDt, _ := crisp.NewSetN(-0.2, 0.2, 40)
	fzDt, _ := fuzzy.NewIDVal("temp/dt", crispDt, map[id.ID]fuzzy.Set{
		"--": fuzzy.NewSetStepDown(-0.2, -0.1),
		"-":  fuzzy.NewSetTriangular(-0.2, -0.1, 0),
		"0":  fuzzy.NewSetTriangular(-0.1, 0, 0.1),
		"+":  fuzzy.NewSetTriangular(0, 0.1, 0.2),
		"++": fuzzy.NewSetStepUp(0.1, 0.2),
	})

	// Output #C
	crispForce, _ := crisp.NewSetN(-4, 4, 80)
	fzForce, _ := fuzzy.NewIDVal("force", crispForce, map[id.ID]fuzzy.Set{
		"--": fuzzy.NewSetStepDown(-4, -1),
		"-":  fuzzy.NewSetTriangular(-2, -1, 0),
		"0":  fuzzy.NewSetTriangular(-1, 0, 1),
		"+":  fuzzy.NewSetTriangular(0, 1, 2),
		"++": fuzzy.NewSetStepUp(1, 4),
	})

	return fzDiff, fzDt, fzForce
}

func TestExample(t *testing.T) {
	fvDiff, fvDt, fvForce := customValues()
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

	Convey("example", t, func() {
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

	Convey("full values", t, func() {
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
