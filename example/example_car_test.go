package example

import (
	"fmt"
	"fugologic/builder"
	"fugologic/crisp"
	"fugologic/fuzzy"
	"fugologic/id"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Converters
const (
	TO_MS  = 3.6       // to m/s
	TO_KMH = 1.0 / 3.6 // to km/h
)

// car simple model
type car struct {
	useBreaks   bool // allow negative force
	speed       float64
	wantedSpeed float64
	err         float64
	derr        float64
}

// compute next speed when "force" is applied
func (c *car) next(force float64) {
	if !c.useBreaks && force < 0 {
		force = 0
	}

	dt := 0.5  // s
	m := 1400. // kgs
	d := 200.  // N (natural car speed down)
	prev := c.speed
	errPrev := c.err
	c.speed = (force-d)*dt/m + prev
	if c.speed < 0 {
		c.speed = 0
	}
	c.err = c.speed - c.wantedSpeed
	c.derr = (c.err - errPrev) / dt
}

// in:  err   = diff from current speed to wanted speed
// in:  derr  = err/dt
// out: force = force required
func TestCar(t *testing.T) {
	Convey("rules", t, func() {
		// Symetrical definition
		newSymIDVal := func(name id.ID, nb int, v [4]float64, k float64) *fuzzy.IDVal {
			a, b, c, d := v[0]*k, v[1]*k, v[2]*k, v[3]*k
			dd := 2*a - d
			cc := 2*a - c
			bb := 2*a - b
			crispCfg, _ := crisp.NewSetN(dd, d, nb)
			fsCfg, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
				"---": fuzzy.StepDown{A: dd, B: cc},          // ▔\▁
				"--":  fuzzy.Triangular{A: dd, B: cc, C: bb}, // ▁/\▁
				"-":   fuzzy.Triangular{A: cc, B: bb, C: a},  // ▁/\▁
				"0":   fuzzy.Triangular{A: bb, B: a, C: b},   // ▁/\▁
				"+":   fuzzy.Triangular{A: a, B: b, C: c},    // ▁/\▁
				"++":  fuzzy.Triangular{A: b, B: c, C: d},    // ▁/\▁
				"+++": fuzzy.StepUp{A: c, B: d},              // ▁/▔
			})
			fvCfg, _ := fuzzy.NewIDVal(name, crispCfg, fsCfg)
			return fvCfg
		}

		// Fuzzy values
		fvErr := newSymIDVal("err", 100, [4]float64{0, 0.05, 0.2, 1}, 10)
		fvDErr := newSymIDVal("derr", 100, [4]float64{0, 0.05, 0.4, 1}, 10)
		fvFrc := newSymIDVal("force", 100, [4]float64{0, 0.15, 0.6, 1}, 2000)

		// Rules
		mx := builder.Mamdani().FuzzyAssoMatrix()
		_ = mx.Asso(fvErr, fvDErr, fvFrc).Matrix(
			[]id.ID{"---", "--", "-", "0", "+", "++", "+++"},
			map[id.ID][]id.ID{
				"---": {"+++", "+++", "+++", "+++", "++", "+", "0"},
				"--":  {"+++", "++", "++", "++", "+", "0", "-"},
				"-":   {"+++", "++", "+", "+", "0", "-", "--"},
				"0":   {"+++", "++", "+", "0", "-", "--", "---"},
				"+":   {"++", "+", "0", "-", "-", "--", "---"},
				"++":  {"+", "0", "-", "--", "--", "--", "---"},
				"+++": {"0", "-", "--", "---", "---", "---", "---"},
			},
		)

		// Engine
		eng, _ := mx.Engine()
		var values [][]float64

		// Evaluate a car in the engine
		evaluate := func(c *car) {
			result, err := eng.Evaluate(fuzzy.DataInput{
				fvErr:  c.err,  // diff with set-point
				fvDErr: c.derr, // current diff derive
			})
			So(err, ShouldBeNil)
			c.next(result[fvFrc])
		}

		// Iterations
		iter := 600

		// Desired speed (in km/h) at iteration #i
		wantedSpeed := map[int]float64{
			0:   5,
			100: 9,
			200: 7.5,
			300: 2,
			400: 6,
			500: 7.5,
			525: 10,
			550: 7.5,
			575: 10,
		}

		car1 := car{
			wantedSpeed: wantedSpeed[0] * TO_MS,
			useBreaks:   false,
		}
		car2 := car{
			wantedSpeed: wantedSpeed[0] * TO_MS,
			useBreaks:   true,
		}

		// Go!
		for i := 0; i < iter; i++ {
			// Fetch desired speed set-point
			ws, ok := wantedSpeed[i]
			if ok {
				car1.wantedSpeed = ws * TO_MS
				car2.wantedSpeed = ws * TO_MS
			}

			// Evaluate engines
			evaluate(&car1)
			evaluate(&car2)
			fmt.Printf(
				"[%4d] set-point: %3.2f, c1: %3.2f km/h, c2: %3.2f km/h, derr1: %f, derr2: %f\n",
				i, car1.wantedSpeed*TO_KMH, car1.speed*TO_KMH, car2.speed*TO_KMH, car1.derr, car2.derr,
			)

			// Save values
			values = append(values, []float64{car1.wantedSpeed * TO_KMH, car1.speed * TO_KMH, car2.speed * TO_KMH, car1.derr, car2.derr})
		}

		// Export
		So(writeCSV("./example_car_test.csv", []string{"wanted", "speed car 1", "speed car 2", "derr 1", "derr 2"}, values), ShouldBeNil)
	})
}
