package fuzzy

import (
	"testing"

	"fugologic.git/crisp"
	. "github.com/smartystreets/goconvey/convey"
)

var defuzzificationNone Defuzzification = func(_ Set, _ crisp.Set) float64 {
	return 0
}

func TestDefuzzerAdd(t *testing.T) {
	newSet := func() Set { return func(x float64) float64 { return x } }

	fvA := NewIDValCustom("a", crisp.Set{})
	fsA1 := NewIDSetCustom("a1", newSet(), &fvA)
	fsA2 := NewIDSetCustom("a2", newSet(), &fvA)
	fsA3 := NewIDSetCustom("a3", newSet(), &fvA)

	Convey("add", t, func() {
		Convey("when empty", func() {
			defuzzer := newDefuzzer(nil)
			So(defuzzer.results, ShouldBeEmpty)
		})

		Convey("when ok", func() {
			defuzzer := newDefuzzer(nil)
			defuzzer.add([]IDSet{fsA1, fsA3})
			defuzzer.add([]IDSet{fsA2})
			So(defuzzer.results, ShouldHaveLength, 3)
		})
	})
}

func TestDefuzzerDefuzz(t *testing.T) {
	setA, _ := crisp.NewSet(1, 4, 0.1)
	fvA := NewIDValCustom("a", setA)
	fsA1 := NewIDSetCustom("a1", NewSetTriangular(1, 2, 3).Min(0.25), &fvA)
	fsA2 := NewIDSetCustom("a2", NewSetTriangular(2, 3, 4).Min(0.75), &fvA)

	setB, _ := crisp.NewSet(11, 14, 0.1)
	fvB := NewIDValCustom("b", setB)
	fsB1 := NewIDSetCustom("b1", NewSetTriangular(11, 12, 13).Min(0.1), &fvB)
	fsB2 := NewIDSetCustom("b2", NewSetTriangular(12, 13, 14).Min(0.9), &fvB)

	Convey("add", t, func() {
		Convey("when empty", func() {
			defuzzer := newDefuzzer(nil)
			result, err := defuzzer.defuzz()
			So(err, ShouldBeNil)
			So(result, ShouldBeEmpty)
		})

		Convey("when ok", func() {
			defuzzer := newDefuzzer(defuzzificationNone)
			defuzzer.add([]IDSet{fsA1, fsA2, fsB1, fsB2})
			result, err := defuzzer.defuzz()
			So(err, ShouldBeNil)
			So(result, ShouldResemble, DataOutput{
				"a": 0,
				"b": 0,
			})
		})
	})
}

func TestDefuzzification(t *testing.T) {
	Convey("centroid", t, func() {
		dx := 0.01
		universe, _ := crisp.NewSet(0, 50, 0.1)

		fs1 := NewSetTrapezoid(10, 15, 25, 30)
		fs2 := NewSetTrapezoid(25, 30, 40, 45)

		Convey("when triangular", func() {
			fs := NewSetTriangular(10, 20, 30)

			defuzz := DefuzzificationCentroid(fs, universe)
			So(defuzz, ShouldAlmostEqual, 20, dx)
			So(fs(defuzz), ShouldAlmostEqual, 1, dx)
		})

		Convey("when trapezoids", func() {
			union := fs1.Union(fs2)
			defuzz := DefuzzificationCentroid(union, universe)
			So(defuzz, ShouldAlmostEqual, 27.5, dx)
			So(union(defuzz), ShouldAlmostEqual, 0.5, dx)
		})

		Convey("when first if truncated", func() {
			fsModif := fs1.Min(0.5)
			union := fsModif.Union(fs2)
			defuzz := DefuzzificationCentroid(union, universe)
			So(defuzz, ShouldAlmostEqual, 29.58, dx)
			So(union(defuzz), ShouldAlmostEqual, 0.92, dx)
		})

		Convey("when second if truncated", func() {
			fsModif := fs2.Min(0.5)
			union := fs1.Union(fsModif)
			defuzz := DefuzzificationCentroid(union, universe)
			So(defuzz, ShouldAlmostEqual, 25.42, dx)
			So(union(defuzz), ShouldAlmostEqual, 0.92, dx)
		})

		Convey("when custom #1", func() {
			Convey("when min", func() {
				fs1 := NewSetTrapezoid(0, 1, 4, 5).Min(0.3)
				fs2 := NewSetTrapezoid(3, 4, 6, 7).Min(0.5)
				fs3 := NewSetTrapezoid(5, 6, 7, 8)

				universe, _ := crisp.NewSet(0, 8, 0.1)
				union := fs1.Union(fs2).Union(fs3)

				defuzz := DefuzzificationCentroid(union, universe)
				So(defuzz, ShouldAlmostEqual, 4.81, dx)
				So(union(defuzz), ShouldEqual, 0.5)
			})

			Convey("when multiply", func() {

				fs1 := NewSetTrapezoid(0, 1, 4, 5).Multiply(0.3)
				fs2 := NewSetTrapezoid(3, 4, 6, 7).Multiply(0.5)
				fs3 := NewSetTrapezoid(5, 6, 7, 8)

				universe, _ := crisp.NewSet(0, 8, 0.1)
				union := fs1.Union(fs2).Union(fs3)

				defuzz := DefuzzificationCentroid(union, universe)
				So(defuzz, ShouldAlmostEqual, 4.9530, dx)
				So(union(defuzz), ShouldEqual, 0.5)
			})
		})

		Convey("when custom #2", func() {
			Convey("when truncate", func() {
				fs1 := NewSetTrapezoid(0, 2, 8, 12).Min(0.9)
				fs2 := NewSetTrapezoid(5, 7, 12, 14).Min(0.5)
				fs3 := NewSetTrapezoid(12, 13, 18, 19).Min(0.1)

				universe, _ := crisp.NewSet(0, 20, 0.1)
				union := fs1.Union(fs2).Union(fs3)

				defuzz := DefuzzificationCentroid(union, universe)
				So(defuzz, ShouldAlmostEqual, 6.9403, dx)
				So(union(defuzz), ShouldEqual, 0.9)
			})

			Convey("when multiply", func() {
				fs1 := NewSetTrapezoid(0, 2, 8, 12).Multiply(0.9)
				fs2 := NewSetTrapezoid(5, 7, 12, 14).Multiply(0.5)
				fs3 := NewSetTrapezoid(12, 13, 18, 19).Multiply(0.1)

				universe, _ := crisp.NewSet(0, 20, 0.1)
				union := fs1.Union(fs2).Union(fs3)

				defuzz := DefuzzificationCentroid(union, universe)
				So(defuzz, ShouldAlmostEqual, 6.7719, dx)
				So(union(defuzz), ShouldEqual, 0.9)
			})
		})
	})

	Convey("smallest, largest max", t, func() {
		Convey("when triangular", func() {
			fs1 := NewSetTriangular(1, 2, 3)
			fs2 := NewSetTriangular(2, 3, 4)
			universe, _ := crisp.NewSet(0, 5, 0.25)

			xsm, xlm := defuzzificationMaximums(fs1.Union(fs2), universe)
			So(xsm, ShouldEqual, 2)
			So(xlm, ShouldEqual, 3)
		})

		Convey("when trapezoid", func() {
			fs1 := NewSetTrapezoid(1, 2, 3, 4)
			universe, _ := crisp.NewSet(0, 5, 0.25)

			xsm, xlm := defuzzificationMaximums(fs1, universe)
			So(xsm, ShouldEqual, 2)
			So(xlm, ShouldEqual, 3)
		})

		Convey("when trapezoid ymax=0.5", func() {
			fs1 := NewSetTrapezoid(1, 2, 3, 4)
			universe, _ := crisp.NewSet(0, 5, 0.1)

			xsm, xlm := defuzzificationMaximums(fs1.Min(0.6), universe)
			So(xsm, ShouldEqual, 1.6)
			So(xlm, ShouldAlmostEqual, 3.3)
		})
	})
}
