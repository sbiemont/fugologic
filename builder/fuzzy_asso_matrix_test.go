package builder

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/sbiemont/fugologic/crisp"
	"github.com/sbiemont/fugologic/fuzzy"
	"github.com/sbiemont/fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

func newTestFAM() FuzzyAssoMatrix {
	return NewFuzzyAssoMatrix(
		fuzzy.OperatorZadeh{}.And,
		fuzzy.ImplicationMin,
		fuzzy.AggregationUnion,
		fuzzy.DefuzzificationCentroid,
	)
}

// Create a fuzzy value, a list of fuzzy sets and link both
func newTestVals(val id.ID, sets ...id.ID) *fuzzy.IDVal {
	// Prepare input
	fuzzySet := func(x float64) float64 { return x }
	idSets := make(map[id.ID]fuzzy.Set)
	for _, set := range sets {
		idSets[set] = fuzzySet
	}

	// New value
	fv, err := fuzzy.NewIDVal(val, crisp.Set{}, idSets)
	So(err, ShouldBeNil)
	return fv
}

func ids(sets []fuzzy.IDSet) []string {
	ids := make([]string, len(sets))
	for i, set := range sets {
		ids[i] = string(set.ID())
	}
	return ids
}

// Compact id from id-sets into string "id1.id2.id3"
func compactIDs(rules []fuzzy.Rule) []string {
	result := make([]string, len(rules))
	for i, rule := range rules {
		inputs, outputs := rule.IO()
		result[i] = fmt.Sprintf("%s=>%s", strings.Join(ids(inputs), "."), strings.Join(ids(outputs), "."))
	}
	sort.Strings(result)
	return result
}

func TestFuzzyAssoMatrix(t *testing.T) {
	Convey("fam", t, func() {
		fvA := newTestVals("a", "a1", "a2")
		fvB := newTestVals("b", "b1", "b2", "b3")
		fvC := newTestVals("c", "c1", "c2", "c3", "c4")

		Convey("when ok", func() {
			Convey("when full definition", func() {
				bld := newTestFAM()
				err := bld.
					Asso(fvA, fvB, fvC).
					Matrix(
						[]id.ID{"a1", "a2"},
						map[id.ID][]id.ID{
							"b1": {"c1", "c2"},
							"b2": {"c3", "c4"},
							"b3": {"c3", "c2"},
						})
				So(err, ShouldBeNil)
				So(bld.rules, ShouldHaveLength, 6)

				// Sort result because of rules random order
				So(compactIDs(bld.rules), ShouldResemble, []string{
					"a1.b1=>c1",
					"a1.b2=>c3",
					"a1.b3=>c3",
					"a2.b1=>c2",
					"a2.b2=>c4",
					"a2.b3=>c2",
				})
			})

			Convey("when empty rules", func() {
				bld := newTestFAM()
				err := bld.
					Asso(fvA, fvB, fvC).
					Matrix(
						[]id.ID{"a1", "a2"},
						map[id.ID][]id.ID{
							"b1": {"c1", "c2"},
							"b2": {"c3", ""},
							"b3": {"", "c2"},
						})
				So(err, ShouldBeNil)
				So(bld.rules, ShouldHaveLength, 4)

				// Sort result because of rules random order
				So(compactIDs(bld.rules), ShouldResemble, []string{
					"a1.b1=>c1",
					"a1.b2=>c3",
					"a2.b1=>c2",
					"a2.b3=>c2",
				})
			})
		})

		Convey("when ko", func() {
			Convey("when not enought output", func() {
				bld := newTestFAM()
				err := bld.
					Asso(fvA, fvB, fvC).
					Matrix(
						[]id.ID{"a1", "a2"},
						map[id.ID][]id.ID{
							"b1": {"c1", "c2"},
							"b2": {"c3"},
							"b3": {"c3", "c2"},
						})
				So(err, ShouldBeError, "rule, sizes should be the same (found: 1, expected: 2)")
			})

			Convey("when duplicated id-set on 'if' statement", func() {
				bld := newTestFAM()
				err := bld.
					Asso(fvA, fvB, fvC).
					Matrix(
						[]id.ID{"a1", "a1"},
						map[id.ID][]id.ID{
							"b1": {"c1", "c2"},
							"b2": {"c3", "c4"},
							"b3": {"c3", "c2"},
						})
				So(err, ShouldBeError, "'if' statement, duplicated headers found")
			})

			Convey("when unknown id-set on 'if' statement", func() {
				bld := newTestFAM()
				err := bld.
					Asso(fvA, fvB, fvC).
					Matrix(
						[]id.ID{"a1", "a0"},
						map[id.ID][]id.ID{
							"b1": {"c1", "c2"},
							"b2": {"c3", "c4"},
							"b3": {"c3", "c2"},
						})
				So(err, ShouldBeError, "'if' statement, cannot find a0 from a")
			})

			Convey("when unknown id-set on 'and' statement", func() {
				bld := newTestFAM()
				err := bld.
					Asso(fvA, fvB, fvC).
					Matrix(
						[]id.ID{"a1", "a2"},
						map[id.ID][]id.ID{
							"b1": {"c1", "c2"},
							"b2": {"c3", "c4"},
							"b4": {"c3", "c2"},
						})
				So(err, ShouldBeError, "'and' statement, cannot find b4 from b")
			})

			Convey("when unknown id-set 'on 'then' statement", func() {
				bld := newTestFAM()
				err := bld.
					Asso(fvA, fvB, fvC).
					Matrix(
						[]id.ID{"a1", "a2"},
						map[id.ID][]id.ID{
							"b1": {"c1", "c2"},
							"b2": {"c3", "c7"},
							"b3": {"c3", "c2"},
						})
				So(err, ShouldBeError, "'then' statement, cannot find c7 from c")
			})
		})
	})
}
