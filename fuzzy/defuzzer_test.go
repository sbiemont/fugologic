package fuzzy

import (
	"testing"

	"fugologic/crisp"
	"fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

var defuzzificationNone Defuzzification = func(_ Set, _ crisp.Set) float64 {
	return 0
}

func TestDefuzzerAdd(t *testing.T) {
	newSet := func() Set { return func(x float64) float64 { return x } }

	fvA, _ := NewIDVal("a", crisp.Set{}, map[id.ID]Set{
		"a1": newSet(),
		"a2": newSet(),
	})

	Convey("new", t, func() {
		Convey("when empty", func() {
			defuzzer := newDefuzzer(nil, AggregationUnion, nil)
			So(defuzzer.results, ShouldBeEmpty)
		})

		Convey("when ok", func() {
			defuzzer := newDefuzzer(DefuzzificationCentroid, AggregationUnion, []IDSet{
				fvA.Get("a1"),
				fvA.Get("a2"),
			})
			So(defuzzer.fct, ShouldEqual, DefuzzificationCentroid)
			So(defuzzer.agg, ShouldEqual, AggregationUnion)
			So(defuzzer.results, ShouldHaveLength, 2)
			So(defuzzer.results[0].uuid, ShouldEqual, "a1")
			So(defuzzer.results[1].uuid, ShouldEqual, "a2")
		})
	})
}

func TestDefuzzerDefuzz(t *testing.T) {
	setA, _ := crisp.NewSet(1, 4, 0.1)
	fvA, _ := NewIDVal("a", setA, map[id.ID]Set{
		"a1": NewSetTriangular(1, 2, 3).Min(0.25),
		"a2": NewSetTriangular(2, 3, 4).Min(0.75),
	})

	setB, _ := crisp.NewSet(11, 14, 0.1)
	fvB, _ := NewIDVal("b", setB, map[id.ID]Set{
		"b1": NewSetTriangular(11, 12, 13).Min(0.1),
		"b2": NewSetTriangular(12, 13, 14).Min(0.9),
	})

	Convey("add", t, func() {
		Convey("when empty", func() {
			defuzzer := newDefuzzer(nil, AggregationUnion, nil)
			result, err := defuzzer.defuzz()
			So(err, ShouldBeNil)
			So(result, ShouldBeEmpty)
		})

		Convey("when ok", func() {
			defuzzer := newDefuzzer(defuzzificationNone, AggregationUnion, []IDSet{
				fvA.Get("a1"), fvA.Get("a2"), fvB.Get("b1"), fvB.Get("b2"),
			})
			result, err := defuzzer.defuzz()
			So(err, ShouldBeNil)
			So(result, ShouldResemble, DataOutput{
				fvA: 0,
				fvB: 0,
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

		Convey("when trapezoid with ymax", func() {
			fs1 := NewSetTrapezoid(1, 2, 3, 4).Min(0.6)
			universe, _ := crisp.NewSet(0, 5, 0.1)

			xsm, xlm := defuzzificationMaximums(fs1, universe)
			So(xsm, ShouldEqual, 1.6)
			So(xlm, ShouldAlmostEqual, 3.3)
		})

		Convey("when smallest/largest of maxs", func() {
			fs1 := NewSetTrapezoid(1, 2, 3, 4).Min(0.6)
			universe, _ := crisp.NewSet(0, 5, 0.1)

			xsm := DefuzzificationSmallestOfMaxs(fs1, universe)
			xlm := DefuzzificationLargestOfMaxs(fs1, universe)
			So(xsm, ShouldEqual, 1.6)
			So(xlm, ShouldAlmostEqual, 3.3)
		})
	})
}
