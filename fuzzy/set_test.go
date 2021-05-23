package fuzzy

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetUnion(t *testing.T) {
	Convey("union", t, func() {
		Convey("when trapezoids", func() {
			fs1 := NewSetTrapezoid(10, 15, 25, 30)
			fs2 := NewSetTrapezoid(25, 30, 40, 45)

			fs3 := fs1.Union(fs2)
			So(fs3(10), ShouldEqual, 0)
			So(fs3(10), ShouldEqual, 0)
			So(fs3(15), ShouldEqual, 1)
			So(fs3(27.5), ShouldEqual, 0.5)
			So(fs3(30), ShouldEqual, 1)
			So(fs3(40), ShouldEqual, 1)
			So(fs3(45), ShouldEqual, 0)
		})
	})
}

func TestSetMin(t *testing.T) {
	Convey("min", t, func() {
		Convey("when trapezoid", func() {
			fs := NewSetTrapezoid(10, 15, 25, 30).Min(0.42)

			So(fs(10), ShouldEqual, 0)
			So(fs(12.5), ShouldEqual, 0.42)
			So(fs(15), ShouldEqual, 0.42)
			So(fs(25), ShouldEqual, 0.42)
			So(fs(27.5), ShouldEqual, 0.42)
			So(fs(30), ShouldEqual, 0)
		})

		Convey("when triangular", func() {
			newSet := func() Set {
				return NewSetTriangular(10, 20, 30)
			}

			Convey("when in ]0 ; 1[", func() {
				fs := newSet().Min(0.42)

				So(fs(10), ShouldEqual, 0)
				So(fs(20), ShouldEqual, 0.42)
				So(fs(30), ShouldEqual, 0)
			})

			Convey("when 0", func() {
				fs := newSet().Min(0)

				So(fs(10), ShouldEqual, 0)
				So(fs(20), ShouldEqual, 0)
				So(fs(30), ShouldEqual, 0)
			})

			Convey("when 1", func() {
				fs := newSet().Min(1)

				So(fs(10), ShouldEqual, 0)
				So(fs(20), ShouldEqual, 1)
				So(fs(30), ShouldEqual, 0)
			})
		})
	})
}

func TestSetMultiply(t *testing.T) {
	Convey("multiply", t, func() {
		Convey("when trapezoid", func() {
			fs := NewSetTrapezoid(10, 15, 25, 30).Multiply(0.42)

			So(fs(10), ShouldEqual, 0)
			So(fs(12.5), ShouldEqual, 0.5*0.42)
			So(fs(15), ShouldEqual, 0.42)
			So(fs(25), ShouldEqual, 0.42)
			So(fs(27.5), ShouldEqual, 0.5*0.42)
			So(fs(30), ShouldEqual, 0)
		})

		Convey("when triangular", func() {
			newSet := func() Set {
				return NewSetTriangular(10, 20, 30)
			}

			Convey("when in ]0 ; 1[", func() {
				fs := newSet().Multiply(0.42)

				So(fs(10), ShouldEqual, 0)
				So(fs(20), ShouldEqual, 0.42)
				So(fs(30), ShouldEqual, 0)
			})

			Convey("when 0", func() {
				fs := newSet().Multiply(0)

				So(fs(10), ShouldEqual, 0)
				So(fs(20), ShouldEqual, 0)
				So(fs(30), ShouldEqual, 0)
			})

			Convey("when 1", func() {
				fs := newSet().Multiply(1)

				So(fs(10), ShouldEqual, 0)
				So(fs(20), ShouldEqual, 1)
				So(fs(30), ShouldEqual, 0)
			})
		})
	})
}

func TestSetMerge(t *testing.T) {
	// Add 1
	var fs1 Set = func(x float64) float64 {
		return x + 1
	}

	// Add 100
	var fs2 Set = func(x float64) float64 {
		return x + 100
	}

	Convey("merge", t, func() {
		Convey("when max", func() {
			fs := fs1.chain(fs2, math.Max)
			So(fs(1), ShouldEqual, 101) // Max(2, 101)
		})

		Convey("when min", func() {
			fs := fs1.chain(fs2, math.Min)
			So(fs(1), ShouldEqual, 2) // Min(2, 101)
		})
	})
}

func TestNewSet(t *testing.T) {
	check := func(fs Set, expected map[float64]float64) {
		for x, exp := range expected {
			So(fs(x), ShouldAlmostEqual, exp, 0.01)
		}
	}

	Convey("triangular", t, func() {
		Convey("when ok", func() {
			check(NewSetTriangular(0.0, 0.5, 1.0), map[float64]float64{
				0.0:  0.0,
				0.25: 0.5,
				0.5:  1.0,
				0.75: 0.5,
				1.0:  0.0,
			})
		})

		Convey("when a=b", func() {
			check(NewSetTriangular(0.0, 0.0, 1.0), map[float64]float64{
				0.0:  1.0,
				0.25: 0.75,
				0.5:  0.5,
				0.75: 0.25,
				1.0:  0.0,
			})
		})

		Convey("when b=c", func() {
			check(NewSetTriangular(0.0, 1.0, 1.0), map[float64]float64{
				0.0:  0.0,
				0.25: 0.25,
				0.5:  0.5,
				0.75: 0.75,
				1.0:  1.0,
			})
		})
	})

	Convey("trapezoid", t, func() {
		Convey("when ok", func() {
			check(NewSetTrapezoid(0.0, 0.25, 0.75, 1.0), map[float64]float64{
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
			check(NewSetTrapezoid(0.0, 0.0, 0.75, 1.0), map[float64]float64{
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
			check(NewSetTrapezoid(0.0, 0.25, 1.0, 1.0), map[float64]float64{
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
			check(NewSetTrapezoid(0.0, 0.0, 1.0, 1.0), map[float64]float64{
				0.0:   1.0,
				0.125: 1.0,
				0.25:  1.0,
				0.5:   1.0,
				0.75:  1.0,
				0.875: 1.0,
				1.0:   1.0,
			})
		})
	})

	Convey("gauss", t, func() {
		Convey("when ok", func() {
			check(NewSetGauss(1.0, 5.0), map[float64]float64{
				1.0: 0.0,
				5.0: 1.0,
				9.0: 0.0,
			})
		})
	})

	Convey("g bell", t, func() {
		Convey("when ok", func() {
			check(NewSetGbell(2.0, 4.0, 6.0), map[float64]float64{
				1.0:  0.0,
				5.0:  1.0,
				7.0:  1.0,
				10.0: 0.0,
			})
		})
	})

	Convey("step up", t, func() {
		Convey("when ok", func() {
			check(NewSetStepUp(2.0, 4.0), map[float64]float64{
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
			check(NewSetStepUp(2.0, 2.0), map[float64]float64{
				1.0: 0.0,
				2.0: 1.0,
				2.5: 1.0,
				3.0: 1.0,
				3.5: 1.0,
				4.0: 1.0,
				5.0: 1.0,
			})
		})
	})

	Convey("step down", t, func() {
		Convey("when ok", func() {
			check(NewSetStepDown(2.0, 4.0), map[float64]float64{
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
			check(NewSetStepDown(2.0, 2.0), map[float64]float64{
				1.0: 1.0,
				2.0: 1.0,
				2.5: 0.0,
				3.0: 0.0,
				3.5: 0.0,
				4.0: 0.0,
				5.0: 0.0,
			})
		})
	})
}
