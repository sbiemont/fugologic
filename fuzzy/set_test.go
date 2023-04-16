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

func TestSetMin(t *testing.T) {
	Convey("min", t, func() {
		Convey("when trapezoid", func() {
			fs, err := Trapezoid{10, 15, 25, 30}.New()
			So(err, ShouldBeNil)
			fs = fs.Min(0.42)

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
				set, _ := Triangular{10, 20, 30}.New()
				return set
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
			fs, err := Trapezoid{10, 15, 25, 30}.New()
			So(err, ShouldBeNil)
			fs = fs.Multiply(0.42)

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
				set, _ := Triangular{10, 20, 30}.New()
				return set
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
