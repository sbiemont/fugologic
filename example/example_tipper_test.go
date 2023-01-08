package example

import (
	"testing"

	"github.com/sbiemont/fugologic/builder"
	"github.com/sbiemont/fugologic/crisp"
	"github.com/sbiemont/fugologic/fuzzy"
	"github.com/sbiemont/fugologic/id"
	. "github.com/smartystreets/goconvey/convey"
)

// https://fr.mathworks.com/help/fuzzy/working-from-the-command-line.html

func TestTipper(t *testing.T) {
	Convey("tipper", t, func() {
		k := 100
		// Input service
		crispSvc, _ := crisp.NewSetN(0, 10, k)
		fsSvc, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"poor":      fuzzy.Gauss{Sigma: 1.5, C: 0},
			"good":      fuzzy.Gauss{Sigma: 1.5, C: 5},
			"excellent": fuzzy.Gauss{Sigma: 1.5, C: 10},
		})
		fvSvc, _ := fuzzy.NewIDVal("service", crispSvc, fsSvc)

		// Input food
		crispFood, _ := crisp.NewSetN(0, 10, k)
		fsFood, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"rancid":    fuzzy.Trapezoid{A: -2, B: 0, C: 1, D: 3},
			"delicious": fuzzy.Trapezoid{A: 7, B: 9, C: 10, D: 12},
		})
		fvFood, _ := fuzzy.NewIDVal("food", crispFood, fsFood)

		// Output tip
		crispTip, _ := crisp.NewSetN(0, 30, 3*k)
		fsTip, _ := fuzzy.NewIDSets(map[id.ID]fuzzy.SetBuilder{
			"cheap":    fuzzy.Triangular{A: 0, B: 5, C: 10},
			"average":  fuzzy.Triangular{A: 10, B: 15, C: 20},
			"generous": fuzzy.Triangular{A: 20, B: 25, C: 30},
		})
		fvTip, _ := fuzzy.NewIDVal("tip", crispTip, fsTip)

		// If (service is poor) or (food is rancid), then (tip is cheap)
		// If (service is good), then (tip is average)
		// If (service is excellent) or (food is delicious), then (tip is generous)
		bld := builder.Mamdani().FuzzyLogic()
		bld.If(fvSvc.Get("poor")).Or(fvFood.Get(("rancid"))).Then(fvTip.Get("cheap"))
		bld.If(fvSvc.Get("good")).Then(fvTip.Get("average"))
		bld.If(fvSvc.Get("excellent")).Or(fvFood.Get(("delicious"))).Then(fvTip.Get("generous"))

		// Evaluate
		eng, err := bld.Engine()
		So(err, ShouldBeNil)

		eval := func(service, food, tip float64) {
			out, err := eng.Evaluate(fuzzy.DataInput{
				fvSvc:  service,
				fvFood: food,
			})
			So(err, ShouldBeNil)
			So(out[fvTip], ShouldAlmostEqual, tip, 1e-2)
		}
		eval(1, 2, 5.5586)
		eval(3, 5, 12.2184)
		eval(2, 7, 7.7885)
		eval(3, 1, 8.9547)
	})
}
