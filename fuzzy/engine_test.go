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
	fsDiff, _ := NewIDSets(map[id.ID]SetBuilder{
		"--": &StepDown{-2, -0.5},
		"-":  &Triangular{-2, -0.5, 0},
		"0":  &Triangular{-0.5, 0, 0.5},
		"+":  &Triangular{0, 0.5, 2},
		"++": &StepUp{0.5, 2},
	})
	fzDiff, _ := NewIDVal("diff/consigne", setDiff, fsDiff)

	sets, _ := NewIDSets(map[id.ID]SetBuilder{
		"--": &Gauss{},
	})
	fzDiff2, _ := NewIDVal("diff/consigne", setDiff, sets)
	_ = fzDiff2

	// Input #2
	setDt, _ := crisp.NewSet(-0.2, 0.2, 0.01)
	fsDt, _ := NewIDSets(map[id.ID]SetBuilder{
		"--": &StepDown{-0.2, -0.1},
		"-":  &Triangular{-0.2, -0.1, 0},
		"0":  &Triangular{-0.1, 0, 0.1},
		"+":  &Triangular{0, 0.1, 0.2},
		"++": &StepUp{0.1, 0.2},
	})
	fzDt, _ := NewIDVal("temp/dt", setDt, fsDt)

	// Output
	setCh, _ := crisp.NewSet(-4, 4, 0.1)
	fsCh, _ := NewIDSets(map[id.ID]SetBuilder{
		"--": &StepDown{-4, -1},
		"-":  &Triangular{-2, -1, 0},
		"0":  &Triangular{-1, 0, 1},
		"+":  &Triangular{0, 1, 2},
		"++": &StepUp{1, 4},
	})
	fzCh, _ := NewIDVal("force", setCh, fsCh)

	return fzDiff, fzDt, fzCh
}

// Custom engine for testing
func customEngine() (Engine, *IDVal, *IDVal, *IDVal, error) {
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

	// Helper: a and b => c
	newRule := func(a, b, c id.ID) Rule {
		return NewRule(
			NewExpression([]Premise{fvDiff.Get(a), fvDt.Get(b)}, OperatorZadeh.And), ImplicationMin, []IDSet{fvCh.Get(c)},
		)
	}
	rules := []Rule{
		newRule("--", "--", "++"),
		newRule("--", "-", "++"),
		newRule("--", "0", "++"),
		newRule("--", "+", "+"),
		newRule("--", "++", "+"),

		newRule("-", "--", "++"),
		newRule("-", "-", "++"),
		newRule("-", "0", "+"),
		newRule("-", "+", "0"),
		newRule("-", "++", "0"),

		newRule("0", "--", "+"),
		newRule("0", "-", "+"),
		newRule("0", "0", "0"),
		newRule("0", "+", "-"),
		newRule("0", "++", "-"),

		newRule("+", "--", "0"),
		newRule("+", "-", "0"),
		newRule("+", "0", "-"),
		newRule("+", "+", "--"),
		newRule("+", "++", "--"),

		newRule("++", "--", "-"),
		newRule("++", "-", "-"),
		newRule("++", "0", "--"),
		newRule("++", "+", "--"),
		newRule("++", "++", "--"),
	}

	engine, err := NewEngine(rules, AggregationUnion, DefuzzificationCentroid)
	return engine, fvDiff, fvDt, fvCh, err
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

func TestEngineEvaluate(t *testing.T) {
	Convey("rules with same output", t, func() {
		// Pressure
		setPressure, errPressure := crisp.NewSet(0, 5, 0.1)
		So(errPressure, ShouldBeNil)
		fsPressure, errFsPressure := NewIDSets(map[id.ID]SetBuilder{
			"high":    StepUp{2, 3},
			"average": Trapezoid{1.4, 1.9, 2.2, 2.6},
		})
		So(errFsPressure, ShouldBeNil)
		fvPressure, errFvPressure := NewIDVal("pressure", setPressure, fsPressure)
		So(errFvPressure, ShouldBeNil)

		// Temperature
		setTemperature, errTemperature := crisp.NewSet(12, 22, 0.1)
		So(errTemperature, ShouldBeNil)
		fsTemperature, errFsTemperature := NewIDSets(map[id.ID]SetBuilder{
			"high": StepUp{16, 19},
		})
		So(errFsTemperature, ShouldBeNil)
		fvTemperature, errFvTemperature := NewIDVal("temperature", setTemperature, fsTemperature)
		So(errFvTemperature, ShouldBeNil)

		// Valve
		setValve, errValve := crisp.NewSet(0, 70, 0.1)
		So(errValve, ShouldBeNil)
		fsValve, errFsValve := NewIDSets(map[id.ID]SetBuilder{
			"wide":    Triangular{30, 40, 50},
			"average": Triangular{20, 30, 40},
		})
		So(errFsValve, ShouldBeNil)
		fvValve, errFvValve := NewIDVal("valve wide open", setValve, fsValve)
		So(errFvValve, ShouldBeNil)

		// Rules
		// p.high & t.hight -> v.wide
		// p.average & t.hight -> v.average
		and := OperatorZadeh.And
		exp1 := NewExpression([]Premise{fvPressure.Get("high"), fvTemperature.Get("high")}, and)
		exp2 := NewExpression([]Premise{fvPressure.Get("average"), fvTemperature.Get("high")}, and)
		rules := []Rule{
			NewRule(exp1, ImplicationMin, []IDSet{fvValve.Get("wide")}),
			NewRule(exp2, ImplicationMin, []IDSet{fvValve.Get("average")}),
		}

		// Evaluate engine
		engine, errEngine := NewEngine(rules, AggregationUnion, DefuzzificationCentroid)
		So(errEngine, ShouldBeNil)

		// Evaluate engine
		result, errEval := engine.Evaluate(DataInput{
			fvPressure:    2.5,
			fvTemperature: 17,
		})
		So(errEval, ShouldBeNil)
		So(result, ShouldResemble, DataOutput{
			fvValve: 35.73264090043862,
		})
	})

	Convey("custom minimalistic test", t, func() {
		// Definitions
		setA, _ := crisp.NewSet(1, 4, 0.1)
		fsA, _ := NewIDSets(map[id.ID]SetBuilder{
			"a1": Triangular{1, 2, 3},
			"a2": Triangular{2, 3, 4},
		})
		fvA, _ := NewIDVal("a", setA, fsA)

		setB, _ := crisp.NewSet(2, 5, 0.1)
		fsB, _ := NewIDSets(map[id.ID]SetBuilder{
			"b1": Triangular{2, 3, 4},
			"b2": Triangular{3, 4, 5},
		})
		fvB, _ := NewIDVal("b", setB, fsB)

		setC, _ := crisp.NewSet(11, 14, 0.1)
		fsC, _ := NewIDSets(map[id.ID]SetBuilder{
			"c1": Triangular{11, 12, 13},
			"c2": Triangular{12, 13, 14},
		})
		fvC, _ := NewIDVal("c", setC, fsC)

		// Rules
		// a1 & b1 -> c1
		// a2 & b2 -> c2
		rules := []Rule{
			NewRule(NewExpression([]Premise{fvA.Get("a1"), fvB.Get("b1")}, OperatorZadeh.And), ImplicationMin, []IDSet{fvC.Get("c1")}),
			NewRule(NewExpression([]Premise{fvA.Get("a2"), fvB.Get("b2")}, OperatorZadeh.And), ImplicationMin, []IDSet{fvC.Get("c2")}),
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
		So(result, ShouldResemble, DataOutput{
			fvC: 12.5,
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
		check := func(engine Engine, fvDiff, fvDt, fvCh *IDVal, err error) {
			So(err, ShouldBeNil)
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

	// https://athena.ecs.csus.edu/~gordonvs/180/WeeklyNotes/03A_FuzzyLogic.pdf
	Convey("inverted pendulum", t, func() {
		// Definitions
		setA, _ := crisp.NewSet(-3, 3, 0.1)
		fsA, _ := NewIDSets(map[id.ID]SetBuilder{
			"-": StepDown{-2, 0},
			"0": Triangular{-2, 0, 2},
			"+": StepUp{0, 2},
		})
		fvA, _ := NewIDVal("t", setA, fsA)

		setB, _ := crisp.NewSet(-6, 6, 0.1)
		fsB, _ := NewIDSets(map[id.ID]SetBuilder{
			"-": StepDown{-5, 0},
			"0": Triangular{-5, 0, 5},
			"+": StepUp{0, 5},
		})
		fvB, _ := NewIDVal("dt", setB, fsB)

		setC, _ := crisp.NewSet(-18, 18, 0.1)
		fsC, _ := NewIDSets(map[id.ID]SetBuilder{
			"--": StepDown{-16, -8},
			"-":  Triangular{-16, -8, 0},
			"0":  Triangular{-8, 0, 8},
			"+":  Triangular{0, 8, 16},
			"++": StepUp{8, 16},
		})
		fvC, _ := NewIDVal("force", setC, fsC)

		// Rules
		// a and b => c
		newRule := func(a, b, c id.ID) Rule {
			return NewRule(
				NewExpression([]Premise{fvA.Get(a), fvB.Get(b)}, OperatorZadeh.And), ImplicationMin, []IDSet{fvC.Get(c)},
			)
		}
		rules := []Rule{
			newRule("+", "+", "++"),
			newRule("+", "0", "+"),
			newRule("+", "-", "0"),

			newRule("0", "+", "+"),
			newRule("0", "0", "0"),
			newRule("0", "-", "-"),

			newRule("-", "+", "0"),
			newRule("-", "0", "-"),
			newRule("-", "-", "--"),
		}

		engine, errEngine := NewEngine(rules, AggregationUnion, DefuzzificationCentroid)
		So(errEngine, ShouldBeNil)

		dataIn := DataInput{
			fvA: -1.5,
			fvB: 2.0,
		}

		// Evaluate engine
		result, errEval := engine.Evaluate(dataIn)
		So(errEval, ShouldBeNil)
		So(result, ShouldResemble, DataOutput{
			fvC: -2.0201342281879104,
		})
	})
}

func TestEngineIO(t *testing.T) {
	Convey("io", t, func() {
		Convey("when empty", func() {
			inputs, outputs := Engine{}.IO()
			So(inputs, ShouldBeEmpty)
			So(outputs, ShouldBeEmpty)
		})

		Convey("when ok", func() {
			eng, inA, inB, outC, err := customEngine()
			So(err, ShouldBeNil)

			inputs, outputs := eng.IO()
			So(inputs, ShouldHaveLength, 50)
			So(outputs, ShouldHaveLength, 25)
			So(IDSets(inputs).IDVals(), ShouldResemble, map[*IDVal]struct{}{
				inA: {},
				inB: {},
			})
			So(IDSets(outputs).IDVals(), ShouldResemble, map[*IDVal]struct{}{
				outC: {},
			})
		})
	})
}

func BenchmarkEngineNTimes(b *testing.B) {
	const evaluations = 2000
	engine, fvDiff, fvDt, _, err := customEngine()
	if err != nil {
		b.FailNow()
	}

	// Evaluate n times the system
	for i := 0; i < evaluations; i++ {
		engine.Evaluate(map[*IDVal]float64{
			fvDiff: -1,
			fvDt:   -0.1,
		})
	}
}

func BenchmarkEngine(b *testing.B) {
	engine, fvDiff, fvDt, _, err := customEngine()
	if err != nil {
		b.FailNow()
	}

	// Evaluate n times the system
	engine.Evaluate(map[*IDVal]float64{
		fvDiff: -1,
		fvDt:   -0.1,
	})
}

func BenchmarkEngineMatrix(b *testing.B) {
	engine, fvDiff, fvDt, _, err := customEngine()
	if err != nil {
		b.FailNow()
	}

	// Evaluate all possible values of the system
	diffValues := fvDiff.U().Values()
	dtValues := fvDt.U().Values()
	for _, diff := range diffValues {
		for _, dt := range dtValues {
			_, err := engine.Evaluate(DataInput{
				fvDiff: diff,
				fvDt:   dt,
			})
			if err != nil {
				b.FailNow()
			}
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
