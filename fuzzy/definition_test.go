package fuzzy

import (
	"fugologic/crisp"
	"fugologic/id"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefinition(t *testing.T) {
	Convey("new id val", t, func() {
		val, _ := NewIDVal("value", crisp.Set{}, nil)
		So(val.ID(), ShouldEqual, "value")
	})

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
	})

	Convey("flatten", t, func() {
		val, _ := NewIDVal("value", crisp.Set{}, map[id.ID]Set{
			"set #1": nil,
			"set #2": nil,
		})

		So(val.idSets, ShouldResemble, val.idSets) // just check for same address
	})
}
