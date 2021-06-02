package fuzzy

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func checkSet(fs Set, expected map[float64]float64) {
	for x, exp := range expected {
		So(fs(x), ShouldAlmostEqual, exp, 0.01)
	}
}

func TestSetNot(t *testing.T) {
	fs1 := NewSetTrapezoid(10, 15, 25, 30)

	Convey("complement", t, func() {
		complement := fs1.Complement()
		checkSet(complement, map[float64]float64{
			5:    1,
			10:   1,
			12.5: 0.5,
			15:   0,
			20:   0,
			25:   0,
			27.5: 0.5,
			30:   1,
			35:   1,
		})
	})
}

func TestSetUnion(t *testing.T) {
	Convey("union", t, func() {
		Convey("when trapezoids", func() {
			fs1 := NewSetTrapezoid(10, 15, 25, 30)
			fs2 := NewSetTrapezoid(25, 30, 40, 45)

			fs3 := fs1.Union(fs2)
			checkSet(fs3, map[float64]float64{
				10:   0,
				15:   1,
				25:   1,
				27.5: 0.5,
				30:   1,
				40:   1,
				45:   0,
			})
		})
	})
}

func TestSetIntersection(t *testing.T) {
	Convey("intersection", t, func() {
		Convey("when trapezoids", func() {
			fs1 := NewSetTrapezoid(10, 15, 25, 30)
			fs2 := NewSetTrapezoid(25, 30, 40, 45)

			fs3 := fs1.Intersection(fs2)
			checkSet(fs3, map[float64]float64{
				10:    0,
				15:    0,
				25:    0,
				26.25: 0.25,
				27.5:  0.5,
				28.75: 0.25,
				30:    0,
				40:    0,
				45:    0,
			})
		})
	})
}

func TestSetMin(t *testing.T) {
	Convey("min", t, func() {
		Convey("when trapezoid", func() {
			fs := NewSetTrapezoid(10, 15, 25, 30).Min(0.42)

			checkSet(fs, map[float64]float64{
				10:   0,
				12.5: 0.42,
				15:   0.42,
				25:   0.42,
				27.5: 0.42,
				30:   0,
			})
		})

		Convey("when triangular", func() {
			newSet := func() Set {
				return NewSetTriangular(10, 20, 30)
			}

			Convey("when in ]0 ; 1[", func() {
				fs := newSet().Min(0.42)

				checkSet(fs, map[float64]float64{
					10: 0,
					20: 0.42,
					30: 0,
				})
			})

			Convey("when 0", func() {
				fs := newSet().Min(0)

				checkSet(fs, map[float64]float64{
					10: 0,
					20: 0,
					30: 0,
				})
			})

			Convey("when 1", func() {
				fs := newSet().Min(1)

				checkSet(fs, map[float64]float64{
					10: 0,
					20: 1,
					30: 0,
				})
			})
		})
	})
}

func TestSetMultiply(t *testing.T) {
	Convey("multiply", t, func() {
		Convey("when trapezoid", func() {
			fs := NewSetTrapezoid(10, 15, 25, 30).Multiply(0.42)

			checkSet(fs, map[float64]float64{
				10:   0,
				12.5: 0.5 * 0.42,
				15:   0.42,
				25:   0.42,
				27.5: 0.5 * 0.42,
				30:   0,
			})
		})

		Convey("when triangular", func() {
			newSet := func() Set {
				return NewSetTriangular(10, 20, 30)
			}

			Convey("when in ]0 ; 1[", func() {
				fs := newSet().Multiply(0.42)

				checkSet(fs, map[float64]float64{
					10: 0,
					20: 0.42,
					30: 0,
				})
			})

			Convey("when 0", func() {
				fs := newSet().Multiply(0)

				checkSet(fs, map[float64]float64{
					10: 0,
					20: 0,
					30: 0,
				})
			})

			Convey("when 1", func() {
				fs := newSet().Multiply(1)

				checkSet(fs, map[float64]float64{
					10: 0,
					20: 1,
					30: 0,
				})
			})
		})
	})
}

func TestSetAggregate(t *testing.T) {
	// Add 1
	var fs1 Set = func(x float64) float64 {
		return x + 1
	}

	// Add 100
	var fs2 Set = func(x float64) float64 {
		return x + 100
	}

	Convey("aggregate", t, func() {
		Convey("when max", func() {
			fs := fs1.aggregate(fs2, math.Max)
			So(fs(1), ShouldEqual, 101) // Max(2, 101)
		})

		Convey("when min", func() {
			fs := fs1.aggregate(fs2, math.Min)
			So(fs(1), ShouldEqual, 2) // Min(2, 101)
		})
	})
}

func TestNewSet(t *testing.T) {
	Convey("triangular", t, func() {
		Convey("when ok", func() {
			checkSet(NewSetTriangular(0.0, 0.5, 1.0), map[float64]float64{
				0.0:  0.0,
				0.25: 0.5,
				0.5:  1.0,
				0.75: 0.5,
				1.0:  0.0,
			})
		})

		Convey("when a=b", func() {
			checkSet(NewSetTriangular(0.0, 0.0, 1.0), map[float64]float64{
				0.0:  1.0,
				0.25: 0.75,
				0.5:  0.5,
				0.75: 0.25,
				1.0:  0.0,
			})
		})

		Convey("when b=c", func() {
			checkSet(NewSetTriangular(0.0, 1.0, 1.0), map[float64]float64{
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
			checkSet(NewSetTrapezoid(0.0, 0.25, 0.75, 1.0), map[float64]float64{
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
			checkSet(NewSetTrapezoid(0.0, 0.0, 0.75, 1.0), map[float64]float64{
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
			checkSet(NewSetTrapezoid(0.0, 0.25, 1.0, 1.0), map[float64]float64{
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
			checkSet(NewSetTrapezoid(0.0, 0.0, 1.0, 1.0), map[float64]float64{
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
			checkSet(NewSetGauss(1.0, 5.0), map[float64]float64{
				1.0: 0.0,
				5.0: 1.0,
				9.0: 0.0,
			})
		})
	})

	Convey("g bell", t, func() {
		Convey("when ok", func() {
			checkSet(NewSetGbell(2.0, 4.0, 6.0), map[float64]float64{
				1.0:  0.0,
				5.0:  1.0,
				7.0:  1.0,
				10.0: 0.0,
			})
		})
	})

	Convey("step up", t, func() {
		Convey("when ok", func() {
			checkSet(NewSetStepUp(2.0, 4.0), map[float64]float64{
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
			checkSet(NewSetStepUp(2.0, 2.0), map[float64]float64{
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
			checkSet(NewSetStepDown(2.0, 4.0), map[float64]float64{
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
			checkSet(NewSetStepDown(2.0, 2.0), map[float64]float64{
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
