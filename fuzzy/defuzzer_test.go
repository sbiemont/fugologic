package fuzzy

import (
	"math"
	"testing"

	"github.com/sbiemont/fugologic/crisp"
	"github.com/sbiemont/fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

var defuzzificationNone Defuzzification = func(_ Set, _ crisp.Set) float64 {
	return 0
}

func fsUnion(fs1 Set, fs ...Set) Set {
	result := fs1
	for _, fs2 := range fs {
		cpy := result
		result = func(x float64) float64 {
			return math.Min(cpy(x), fs2(x))
		}
	}
	return result
}

func TestDefuzzerDefuzz(t *testing.T) {
	setA, _ := crisp.NewSet(1, 4, 0.1)
	fsA1, _ := Triangular{1, 2, 3}.New()
	fsA2, _ := Triangular{2, 3, 4}.New()
	fvA, _ := NewIDVal("a", setA, map[id.ID]Set{
		"a1": fsA1.Min(0.25),
		"a2": fsA2.Min(0.75),
	})

	setB, _ := crisp.NewSet(11, 14, 0.1)
	fsB1, _ := Triangular{11, 12, 13}.New()
	fsB2, _ := Triangular{12, 13, 14}.New()
	fvB, _ := NewIDVal("b", setB, map[id.ID]Set{
		"b1": fsB1.Min(0.1),
		"b2": fsB2.Min(0.9),
	})

	Convey("add", t, func() {
		Convey("when empty", func() {
			defuzzer := newDefuzzer(nil, AggregationUnion)
			result := defuzzer.defuzz(nil)
			So(result, ShouldBeEmpty)
		})

		Convey("when ok", func() {
			Convey("when several values", func() {
				defuzzer := newDefuzzer(defuzzificationNone, AggregationUnion)
				result := defuzzer.defuzz([]IDSet{
					fvA.Get("a1"), fvA.Get("a2"), fvB.Get("b1"), fvB.Get("b2"),
				})
				So(result, ShouldResemble, DataOutput{
					fvA: 0,
					fvB: 0,
				})
			})

			Convey("when one value", func() {
				defuzzer := newDefuzzer(defuzzificationNone, AggregationUnion)
				result := defuzzer.defuzz([]IDSet{
					fvA.Get("a1"), fvA.Get("a2"),
				})
				So(result, ShouldResemble, DataOutput{
					fvA: 0,
				})
			})
		})
	})
}

func TestDefuzzification(t *testing.T) {
	Convey("centroid", t, func() {
		dx := 0.01
		universe, _ := crisp.NewSet(0, 50, 0.1)

		fs1, _ := Trapezoid{10, 15, 25, 30}.New()
		fs2, _ := Trapezoid{25, 30, 40, 45}.New()

		Convey("when zero", func() {
			fs := func(float64) float64 { return 42 }
			defuzz := DefuzzificationCentroid(fs, crisp.Set{})
			So(defuzz, ShouldEqual, 0)
			So(fs(defuzz), ShouldEqual, 42)
		})

		Convey("when triangular", func() {
			fs, _ := Triangular{10, 20, 30}.New()

			defuzz := DefuzzificationCentroid(fs, universe)
			So(defuzz, ShouldAlmostEqual, 20, dx)
			So(fs(defuzz), ShouldAlmostEqual, 1, dx)
		})

		Convey("when trapezoids", func() {
			union := fsUnion(fs1, fs2)
			defuzz := DefuzzificationCentroid(union, universe)
			So(defuzz, ShouldAlmostEqual, 27.5, dx)
			So(union(defuzz), ShouldAlmostEqual, 0.5, dx)
		})

		Convey("when first if truncated", func() {
			fsModif := fs1.Min(0.5)
			union := fsUnion(fsModif, fs2)
			defuzz := DefuzzificationCentroid(union, universe)
			So(defuzz, ShouldAlmostEqual, 29.58, dx)
			So(union(defuzz), ShouldAlmostEqual, 0.92, dx)
		})

		Convey("when second if truncated", func() {
			fsModif := fs2.Min(0.5)
			union := fsUnion(fs1, fsModif)
			defuzz := DefuzzificationCentroid(union, universe)
			So(defuzz, ShouldAlmostEqual, 25.42, dx)
			So(union(defuzz), ShouldAlmostEqual, 0.92, dx)
		})

		Convey("when custom #1", func() {
			Convey("when min", func() {
				fs1, _ := Trapezoid{0, 1, 4, 5}.New()
				fs2, _ := Trapezoid{3, 4, 6, 7}.New()
				fs3, _ := Trapezoid{5, 6, 7, 8}.New()
				fs1 = fs1.Min(0.3)
				fs2 = fs2.Min(0.5)

				universe, _ := crisp.NewSet(0, 8, 0.1)
				union := fsUnion(fs1, fs2, fs3)

				defuzz := DefuzzificationCentroid(union, universe)
				So(defuzz, ShouldAlmostEqual, 4.81, dx)
				So(union(defuzz), ShouldEqual, 0.5)
			})

			Convey("when multiply", func() {
				fs1, _ := Trapezoid{0, 1, 4, 5}.New()
				fs2, _ := Trapezoid{3, 4, 6, 7}.New()
				fs3, _ := Trapezoid{5, 6, 7, 8}.New()
				fs1 = fs1.Multiply(0.3)
				fs2 = fs2.Multiply(0.5)

				universe, _ := crisp.NewSet(0, 8, 0.1)
				union := fsUnion(fs1, fs2, fs3)

				defuzz := DefuzzificationCentroid(union, universe)
				So(defuzz, ShouldAlmostEqual, 4.9530, dx)
				So(union(defuzz), ShouldEqual, 0.5)
			})
		})

		Convey("when simple trapezoid", func() {
			u, _ := crisp.NewSet(-10, 10, 0.1)
			fs, _ := Trapezoid{-10, -8, -4, 7}.New()

			defuzz := DefuzzificationCentroid(fs, u)
			So(defuzz, ShouldAlmostEqual, -3.2857, dx)
			So(fs(defuzz), ShouldAlmostEqual, 0.935, dx)
		})

		Convey("when custom #2", func() {
			Convey("when min", func() {
				fs1, _ := Trapezoid{0, 2, 8, 12}.New()
				fs2, _ := Trapezoid{5, 7, 12, 14}.New()
				fs3, _ := Trapezoid{12, 13, 18, 19}.New()
				fs1 = fs1.Min(0.9)
				fs2 = fs2.Min(0.5)
				fs3 = fs3.Min(0.1)

				universe, _ := crisp.NewSet(0, 20, 0.1)
				union := fsUnion(fs1, fs2, fs3)

				defuzz := DefuzzificationCentroid(union, universe)
				So(defuzz, ShouldAlmostEqual, 6.9403, dx)
				So(union(defuzz), ShouldEqual, 0.9)
			})

			Convey("when multiply", func() {
				// https://www.mathworks.com/help/fuzzy/defuzzification-methods.html
				fs1, _ := Trapezoid{0, 2, 8, 12}.New()
				fs2, _ := Trapezoid{5, 7, 12, 14}.New()
				fs3, _ := Trapezoid{12, 13, 18, 19}.New()
				fs1 = fs1.Multiply(0.9)
				fs2 = fs2.Multiply(0.5)
				fs3 = fs3.Multiply(0.1)

				universe, _ := crisp.NewSet(0, 20, 0.1)
				union := fsUnion(fs1, fs2, fs3)

				defuzz := DefuzzificationCentroid(union, universe)
				So(defuzz, ShouldAlmostEqual, 6.7719, dx)
				So(union(defuzz), ShouldEqual, 0.9)
			})
		})
	})

	Convey("bisector", t, func() {
		// https://www.mathworks.com/help/fuzzy/defuzzification-methods.html
		dx := 0.01

		Convey("when multiply", func() {
			fs1, _ := Trapezoid{0, 2, 8, 12}.New()
			fs2, _ := Trapezoid{5, 7, 12, 14}.New()
			fs3, _ := Trapezoid{12, 13, 18, 19}.New()
			fs1 = fs1.Multiply(0.9)
			fs2 = fs2.Multiply(0.5)
			fs3 = fs3.Multiply(0.1)

			universe, _ := crisp.NewSet(0, 20, 0.1)
			union := fsUnion(fs1, fs2, fs3)

			defuzz := DefuzzificationBisector(union, universe)
			So(defuzz, ShouldAlmostEqual, 6.3, dx)
			So(union(defuzz), ShouldEqual, 0.9)
		})

		Convey("when truncate", func() {
			fs1, _ := Trapezoid{0, 2, 8, 12}.New()
			fs2, _ := Trapezoid{5, 7, 12, 14}.New()
			fs3, _ := Trapezoid{12, 13, 18, 19}.New()
			fs1 = fs1.Min(0.9)
			fs2 = fs2.Min(0.5)
			fs3 = fs3.Min(0.1)

			universe, _ := crisp.NewSet(0, 20, 0.1)
			union := fsUnion(fs1, fs2, fs3)

			defuzz := DefuzzificationBisector(union, universe)
			So(defuzz, ShouldAlmostEqual, 6.5, dx)
			So(union(defuzz), ShouldEqual, 0.9)
		})

		Convey("when constant", func() {
			constant := func(float64) float64 { return 0.42 }
			universe, _ := crisp.NewSet(0, 20, 0.1)

			defuzz := DefuzzificationBisector(constant, universe)
			So(defuzz, ShouldEqual, 10)             // middle point [0 ; 20]
			So(constant(defuzz), ShouldEqual, 0.42) // constant
		})

		Convey("when left", func() {
			fs1, _ := StepDown{0, 0.2}.New()
			universe, _ := crisp.NewSet(0, 20, 0.1)

			defuzz := DefuzzificationBisector(fs1, universe)
			So(defuzz, ShouldAlmostEqual, 0.1, dx)
			So(fs1(defuzz), ShouldAlmostEqual, 0.5, dx)
		})

		Convey("when right", func() {
			fs1, _ := StepUp{19.8, 20}.New()
			universe, _ := crisp.NewSet(0, 20, 0.1)

			defuzz := DefuzzificationBisector(fs1, universe)
			So(defuzz, ShouldAlmostEqual, 19.9, dx)
			So(fs1(defuzz), ShouldAlmostEqual, 0.5, dx)
		})
	})

	Convey("smallest, largest max", t, func() {
		Convey("when triangular", func() {
			fs1, err1 := Triangular{1, 2, 3}.New()
			fs2, err2 := Triangular{2, 3, 4}.New()
			So(err1, ShouldBeNil)
			So(err2, ShouldBeNil)
			universe, errU := crisp.NewSet(0, 5, 0.25)
			So(errU, ShouldBeNil)

			xsm, xlm := defuzzificationMaximums(fsUnion(fs1, fs2), universe)
			So(xsm, ShouldEqual, 2)
			So(xlm, ShouldEqual, 3)

			// Same checks
			xsm = DefuzzificationSmallestOfMaxs(fsUnion(fs1, fs2), universe)
			xmm := DefuzzificationMiddleOfMaxs(fsUnion(fs1, fs2), universe)
			xlm = DefuzzificationLargestOfMaxs(fsUnion(fs1, fs2), universe)
			So(xsm, ShouldEqual, 2)
			So(xmm, ShouldEqual, 2.5)
			So(xlm, ShouldEqual, 3)
		})

		Convey("when trapezoid", func() {
			fs1, _ := Trapezoid{1, 2, 3, 4}.New()
			universe, _ := crisp.NewSet(0, 5, 0.25)

			xsm, xlm := defuzzificationMaximums(fs1, universe)
			So(xsm, ShouldEqual, 2)
			So(xlm, ShouldEqual, 3)

			// Same checks
			xsm = DefuzzificationSmallestOfMaxs(fs1, universe)
			xmm := DefuzzificationMiddleOfMaxs(fs1, universe)
			xlm = DefuzzificationLargestOfMaxs(fs1, universe)
			So(xsm, ShouldEqual, 2)
			So(xmm, ShouldEqual, 2.5)
			So(xlm, ShouldEqual, 3)
		})

		Convey("when trapezoid with ymax", func() {
			fs1, _ := Trapezoid{1, 2, 3, 4}.New()
			fs1 = fs1.Min(0.6)
			universe, _ := crisp.NewSet(0, 5, 0.1)

			xsm, xlm := defuzzificationMaximums(fs1, universe)
			So(xsm, ShouldEqual, 1.6)
			So(xlm, ShouldAlmostEqual, 3.3)

			// Same checks
			xsm = DefuzzificationSmallestOfMaxs(fs1, universe)
			xmm := DefuzzificationMiddleOfMaxs(fs1, universe)
			xlm = DefuzzificationLargestOfMaxs(fs1, universe)
			So(xsm, ShouldEqual, 1.6)
			So(xmm, ShouldEqual, 2.45)
			So(xlm, ShouldAlmostEqual, 3.3)
		})
	})
}
