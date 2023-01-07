package fuzzy

import (
	"testing"

	"github.com/sbiemont/fugologic/crisp"
	"github.com/sbiemont/fugologic/id"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIDSet(t *testing.T) {
	Convey("accessors", t, func() {
		Convey("when id()", func() {
			item := IDSet{
				uuid: "set #1",
			}
			So(item.ID(), ShouldEqual, "set #1")
		})

		Convey("when evaluate()", func() {
			parent := &IDVal{
				uuid: "val #1",
			}
			item := IDSet{
				uuid:   "set #1",
				set:    func(x float64) float64 { return x },
				parent: parent,
			}

			Convey("when ok", func() {
				eval, err := item.Evaluate(DataInput{
					parent: 42.42,
				})
				So(err, ShouldBeNil)
				So(eval, ShouldEqual, 42.42)
			})

			Convey("when ko", func() {
				parent2 := &IDVal{
					uuid: "val #2",
				}
				eval, err := item.Evaluate(DataInput{
					parent2: 42.42,
				})
				So(err, ShouldBeError, "input: cannot find data for id val `val #1` (id set `set #1`)")
				So(eval, ShouldBeZeroValue)
			})
		})
	})
}

func TestIDVal(t *testing.T) {
	Convey("new id val", t, func() {
		Convey("when empty", func() {
			val, _ := NewIDVal("value", crisp.Set{}, nil)
			So(val, ShouldResemble, &IDVal{
				uuid:   "value",
				u:      crisp.Set{},
				idSets: map[id.ID]IDSet{},
			})
		})

		Convey("when id sets", func() {
			u, err := crisp.NewSet(0, 1, 0.5)
			So(err, ShouldBeNil)

			f1 := func(float64) float64 { return 1 }
			f2 := func(float64) float64 { return 2 }
			val, _ := NewIDVal("value", u, map[id.ID]Set{
				"set #1": f1,
				"set #2": f2,
			})

			So(val.uuid, ShouldEqual, "value")
			So(val.u, ShouldResemble, u)
			So(val.idSets, ShouldHaveLength, 2)

			So(val.idSets["set #1"].parent, ShouldEqual, val)
			So(val.idSets["set #1"].set, ShouldEqual, f1)
			So(val.idSets["set #1"].uuid, ShouldEqual, "set #1")

			So(val.idSets["set #2"].parent, ShouldEqual, val)
			So(val.idSets["set #2"].set, ShouldEqual, f2)
			So(val.idSets["set #2"].uuid, ShouldEqual, "set #2")
		})

		Convey("when ok", func() {
			val, _ := NewIDVal("value", crisp.Set{}, map[id.ID]Set{
				"set #1": nil,
				"set #2": nil,
			})

			So(val.idSets, ShouldResemble, val.idSets) // just check for same address
		})
	})

	Convey("accessors", t, func() {
		cset, err := crisp.NewSet(0, 1, 0.1)
		So(err, ShouldBeNil)
		val, _ := NewIDVal("value", cset, map[id.ID]Set{
			"set #1": nil,
			"set #2": nil,
		})

		Convey("when id()", func() {
			So(val.ID(), ShouldEqual, "value")
		})

		Convey("when u()", func() {
			cset, err := crisp.NewSetN(0, 1, 11)
			So(err, ShouldBeNil)
			So(val.U(), ShouldResemble, cset)
		})

		Convey("when get()", func() {
			So(val.Get("set #1").uuid, ShouldEqual, "set #1")
			So(val.Get("set #2").uuid, ShouldEqual, "set #2")
			So(val.Get("set #3").uuid, ShouldBeEmpty)
		})

		Convey("when fetch()", func() {
			v1, ok1 := val.Fetch("set #1")
			v2, ok2 := val.Fetch("set #2")
			v3, ok3 := val.Fetch("set #3")

			So(ok1, ShouldBeTrue)
			So(v1.uuid, ShouldEqual, "set #1")
			So(ok2, ShouldBeTrue)
			So(v2.uuid, ShouldEqual, "set #2")
			So(ok3, ShouldBeFalse)
			So(v3.uuid, ShouldBeEmpty)
		})
	})
}

func TestIDSets(t *testing.T) {
	Convey("id sets", t, func() {
		Convey("extract id vals", func() {
			// Prepare data
			a, errA := NewIDVal("a", crisp.Set{}, map[id.ID]Set{
				"a1": nil,
				"a2": nil,
			})
			So(errA, ShouldBeNil)

			b, errB := NewIDVal("b", crisp.Set{}, map[id.ID]Set{
				"b1": nil,
			})
			So(errB, ShouldBeNil)

			// Check
			So(IDSets{}.IDVals(), ShouldBeEmpty)
			So(IDSets{a.Get("a1")}.IDVals(), ShouldResemble, map[*IDVal]struct{}{a: {}})
			So(IDSets{b.Get("b1")}.IDVals(), ShouldResemble, map[*IDVal]struct{}{b: {}})
			So(IDSets{a.Get("a1"), a.Get("a2"), b.Get("b1")}.IDVals(), ShouldResemble, map[*IDVal]struct{}{
				a: {},
				b: {},
			})
		})
	})
}
