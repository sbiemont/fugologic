package example

import (
	"testing"

	"fugologic/builder"
	"fugologic/crisp"
	"fugologic/fuzzy"
	"fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFuzzyAssoMatrix(t *testing.T) {
	// https://en.wikipedia.org/wiki/Fuzzy_associative_matrix
	//
	// HP: Hit Point
	// FP: Fire Power
	//
	// HP/FP        | Very low HP | Low HP   | Medium HP | High HP      | Very high HP
	// -------------|-------------|----------|-----------|--------------|-------------
	// Very weak FP | Retreat!    | Retreat! | Defend    | Defend       | Defend
	// Weak FP      | Retreat!    | Defend   | Defend    | Attack       | Attack
	// Medium FP    | Retreat!    | Defend   | Attack    | Attack       | Full attack!
	// High FP      | Retreat!    | Defend   | Attack    | Attack       | Full attack!
	// Very high FP | Defend      | Attack   | Attack    | Full attack! | Full attack!
	//
	Convey("rules", t, func() {
		// Input HP
		crispHP, _ := crisp.NewSetN(0, 100, 1000)
		fsHP, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"Very low HP":  fuzzy.StepDown{A: 0, B: 20},
			"Low HP":       fuzzy.Trapezoid{A: 0, B: 20, C: 40, D: 60},
			"Medium HP":    fuzzy.Triangular{A: 50, B: 50, C: 60},
			"High HP":      fuzzy.Trapezoid{A: 40, B: 60, C: 80, D: 100},
			"Very high HP": fuzzy.StepUp{A: 80, B: 100},
		})
		fvHP, _ := fuzzy.NewIDVal("HP", crispHP, fsHP)

		// Input FP
		crispFP, _ := crisp.NewSetN(0, 100, 1000)
		fsFP, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"Very weak FP": fuzzy.StepDown{A: 0, B: 20},
			"Weak FP":      fuzzy.Trapezoid{A: 0, B: 20, C: 40, D: 60},
			"Medium FP":    fuzzy.Triangular{A: 50, B: 50, C: 60},
			"High FP":      fuzzy.Trapezoid{A: 40, B: 60, C: 80, D: 100},
			"Very high FP": fuzzy.StepUp{A: 80, B: 100},
		})
		fvFP, _ := fuzzy.NewIDVal("FP", crispFP, fsFP)

		// Output Action
		crispAct, _ := crisp.NewSetN(-10, 10, 100)
		fsAct, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"Retreat!":     fuzzy.StepDown{A: -10, B: -5},
			"Defend":       fuzzy.Triangular{A: -10, B: -5, C: 5},
			"Attack":       fuzzy.Triangular{A: -5, B: 5, C: 10},
			"Full attack!": fuzzy.StepUp{A: 5, B: 10},
		})
		fvAct, _ := fuzzy.NewIDVal("Act", crispAct, fsAct)

		bld := builder.NewFuzzyAssoMatrixMamdani()
		bld.Asso(fvHP, fvFP, fvAct).Matrix(
			[]id.ID{"Very low HP", "Low HP", "Medium HP", "High HP", "Very high HP"},
			map[id.ID][]id.ID{
				"Very weak FP": {"Retreat!", "Retreat!", "Defend", "Defend", "Defend"},
				"Weak FP":      {"Retreat!", "Defend", "Defend", "Attack", "Attack"},
				"Medium FP":    {"Retreat!", "Defend", "Attack", "Attack", "Full attack!"},
				"High FP":      {"Retreat!", "Defend", "Attack", "Attack", "Full attack!"},
				"Very high FP": {"Defend", "Attack", "Attack", "Full attack!", "Full attack!"},
			},
		)

		eng, err := bld.Engine()
		So(err, ShouldBeNil)

		result, errEval := eng.Evaluate(fuzzy.DataInput{
			fvHP: 75,
			fvFP: 30,
		})
		So(errEval, ShouldBeNil)
		So(result, ShouldResemble, fuzzy.DataOutput{
			fvAct: 3.3326461897890463,
		})
	})
}
