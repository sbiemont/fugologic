package fuzzy

import (
	"testing"

	"github.com/sbiemont/fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSet(t *testing.T) {
	Convey("triangular", t, func() {
		Convey("when ok", func() {
			fs, err := Triangular{0.0, 0.5, 1.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				0.0:  0.0,
				0.25: 0.5,
				0.5:  1.0,
				0.75: 0.5,
				1.0:  0.0,
			})
		})

		Convey("when a=b", func() {
			fs, err := Triangular{0.0, 0.0, 1.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				0.0:  1.0,
				0.25: 0.75,
				0.5:  0.5,
				0.75: 0.25,
				1.0:  0.0,
			})
		})

		Convey("when b=c", func() {
			fs, err := Triangular{0.0, 1.0, 1.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				0.0:  0.0,
				0.25: 0.25,
				0.5:  0.5,
				0.75: 0.75,
				1.0:  1.0,
			})

			Convey("when ko", func() {
				fs, err := Triangular{0, 2, 1}.New()
				So(err, ShouldBeError, "tri: params shall be sorted")
				So(fs, ShouldBeNil)
			})
		})
	})

	Convey("trapezoid", t, func() {
		Convey("when ok", func() {
			fs, err := Trapezoid{0.0, 0.25, 0.75, 1.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				0.0:   0.0,
				0.125: 0.5,
				0.25:  1.0,
				0.5:   1.0,
				0.75:  1.0,
				0.875: 0.5,
				1.0:   0.0,
			})
		})

		Convey("when a=b", func() {
			fs, err := Trapezoid{0.0, 0.0, 0.75, 1.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				0.0:   1.0,
				0.125: 1.0,
				0.25:  1.0,
				0.5:   1.0,
				0.75:  1.0,
				0.875: 0.5,
				1.0:   0.0,
			})
		})

		Convey("when c=d", func() {
			fs, err := Trapezoid{0.0, 0.25, 1.0, 1.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				0.0:   0.0,
				0.125: 0.5,
				0.25:  1.0,
				0.5:   1.0,
				0.75:  1.0,
				0.875: 1.0,
				1.0:   1.0,
			})
		})

		Convey("when a=b and c=d", func() {
			fs, err := Trapezoid{0.0, 0.0, 1.0, 1.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				0.0:   1.0,
				0.125: 1.0,
				0.25:  1.0,
				0.5:   1.0,
				0.75:  1.0,
				0.875: 1.0,
				1.0:   1.0,
			})

			Convey("when ko", func() {
				fs, err := Trapezoid{0, 4, 1, 2}.New()
				So(err, ShouldBeError, "trap: params shall be sorted")
				So(fs, ShouldBeNil)
			})
		})
	})

	Convey("gauss", t, func() {
		Convey("when ok", func() {
			gauss, err := Gauss{1.0, 5.0}.New()
			So(err, ShouldBeNil)
			checkSet(gauss, map[float64]float64{
				1: 0.0,
				3: 0.13,
				5: 1.0,
				7: 0.13,
				9: 0.0,
			})
		})

		Convey("when sigma==0", func() {
			gauss, err := Gauss{0.0, 5.0}.New()
			So(err, ShouldBeError, "gauss: first parameter must be non zero")
			So(gauss, ShouldBeNil)
		})
	})

	Convey("g bell", t, func() {
		Convey("when ok", func() {
			gbell, err := Gbell{2.0, 4.0, 6.0}.New()
			So(err, ShouldBeNil)
			checkSet(gbell, map[float64]float64{
				1.0:  0.0,
				5.0:  1.0,
				7.0:  1.0,
				10.0: 0.0,
			})
		})

		Convey("when a=0", func() {
			gbell, err := Gbell{0.0, 4.0, 6.0}.New()
			So(err, ShouldBeError, "gbell: first parameter must be non zero")
			So(gbell, ShouldBeNil)
		})
	})

	Convey("step up", t, func() {
		Convey("when ok", func() {
			fs, err := StepUp{2.0, 4.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				1.0: 0.0,
				2.0: 0.0,
				2.5: 0.25,
				3.0: 0.5,
				3.5: 0.75,
				4.0: 1.0,
				5.0: 1.0,
			})
		})

		Convey("when a=b", func() {
			fs, err := StepUp{2.0, 2.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				1.0: 0.0,
				2.0: 1.0,
				2.5: 1.0,
				3.0: 1.0,
				3.5: 1.0,
				4.0: 1.0,
				5.0: 1.0,
			})

			Convey("when a>b", func() {
				fs, err := StepUp{4.0, 2.0}.New()
				So(err, ShouldBeError, "step-up: params shall be sorted")
				So(fs, ShouldBeNil)
			})
		})
	})

	Convey("step down", t, func() {
		Convey("when ok", func() {
			fs, err := StepDown{2.0, 4.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				1.0: 1.0,
				2.0: 1.0,
				2.5: 0.75,
				3.0: 0.5,
				3.5: 0.25,
				4.0: 0.0,
				5.0: 0.0,
			})
		})

		Convey("when a=b", func() {
			fs, err := StepDown{2.0, 2.0}.New()
			So(err, ShouldBeNil)
			checkSet(fs, map[float64]float64{
				1.0: 1.0,
				2.0: 1.0,
				2.5: 0.0,
				3.0: 0.0,
				3.5: 0.0,
				4.0: 0.0,
				5.0: 0.0,
			})
		})

		Convey("when a>b", func() {
			fs, err := StepDown{4, 2}.New()
			So(err, ShouldBeError, "step-down: params shall be sorted")
			So(fs, ShouldBeNil)
		})
	})

	Convey("sigmoid", t, func() {
		Convey("when S shape", func() {
			sig, err := Sigmoid{2.0, 6.0}.New()
			So(err, ShouldBeNil)
			checkSet(sig, map[float64]float64{
				2:  0.0,
				4:  0.018,
				6:  0.5,
				8:  0.982,
				10: 1.0,
			})
		})

		Convey("when Z shape", func() {
			sig, err := Sigmoid{-2.0, 6.0}.New()
			So(err, ShouldBeNil)
			checkSet(sig, map[float64]float64{
				2:  1.0,
				4:  0.982,
				6:  0.5,
				8:  0.018,
				10: 0.0,
			})
		})
	})
}

func TestCheckSorted(t *testing.T) {
	Convey("when empty", t, func() {
		So(checkSorted("fct"), ShouldBeNil)
	})

	Convey("when sorted", t, func() {
		So(checkSorted("fct", 1), ShouldBeNil)
		So(checkSorted("fct", 1, 1), ShouldBeNil)
		So(checkSorted("fct", 1, 2), ShouldBeNil)
		So(checkSorted("fct", 1, 2, 3), ShouldBeNil)
		So(checkSorted("fct", -3, -2, -1), ShouldBeNil)
	})

	Convey("when not sorted", t, func() {
		So(checkSorted("fct", 2, 1), ShouldBeError, "fct: params shall be sorted")
		So(checkSorted("fct", 3, 2, 1), ShouldBeError, "fct: params shall be sorted")
		So(checkSorted("fct", 4, 3, 2, 1), ShouldBeError, "fct: params shall be sorted")
	})
}

func TestNewIDSets(t *testing.T) {
	Convey("new id sets", t, func() {
		Convey("when empty", func() {
			sets, err := NewIDSets(nil)
			So(err, ShouldBeNil)
			So(sets, ShouldBeEmpty)
		})

		Convey("when ok", func() {
			sets, err := NewIDSets(map[id.ID]SetBuilder{
				"fs1": Triangular{1, 2, 3},
				"fs2": Triangular{4, 5, 6},
			})
			So(err, ShouldBeNil)
			So(sets, ShouldHaveLength, 2)
			So(sets["fs1"], ShouldNotBeNil)
			So(sets["fs2"], ShouldNotBeNil)
		})

		Convey("when ko", func() {
			sets, err := NewIDSets(map[id.ID]SetBuilder{
				"fs1": Triangular{1, 2, 3},
				"fs2": Triangular{3, 2, 1},
			})
			So(err, ShouldBeError, "fs2: tri: params shall be sorted")
			So(sets, ShouldBeNil)
		})
	})
}
