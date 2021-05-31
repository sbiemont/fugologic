package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGraph(t *testing.T) {
	Convey("when ok #1", t, func() {
		a := NewNode("a")
		b := NewNode("b")
		c := NewNode("c")
		d := NewNode("d")
		e := NewNode("e")
		f := NewNode("f")
		g := NewNode("g")
		h := NewNode("h")

		Convey("when no edge (all independant nodes", func() {
			dg := NewDirectedGraph([]*Node{a, b, c, d, e, f, g, h})

			Convey("cycle", func() {
				So(dg.isCyclic(), ShouldBeFalse)
			})

			Convey("topo", func() {
				topo, err := dg.Flatten()
				So(err, ShouldBeNil)
				So(topo, ShouldResemble, []*Node{h, g, f, e, d, c, b, a})
			})
		})

		Convey("when mini graph", func() {
			dg := NewDirectedGraph([]*Node{a, b, c})
			dg.AddEdge(a, c).AddEdge(b, c)

			Convey("cycle", func() {
				So(dg.isCyclic(), ShouldBeFalse)
			})

			Convey("topo", func() {
				topo, err := dg.Flatten()
				So(err, ShouldBeNil)
				So(topo, ShouldResemble, []*Node{b, a, c})
			})
		})

		// a -> c
		// c -> g
		// b -> c, e
		// d -> e, f
		// e -> g
		Convey("with order #1", func() {
			dg := NewDirectedGraph([]*Node{a, b, c, d, e, f, g, h})
			dg.
				AddEdge(a, c).
				AddEdge(b, c, e).
				AddEdge(c, g).
				AddEdge(d, e, f).
				AddEdge(e, g)

			Convey("cycle", func() {
				So(dg.isCyclic(), ShouldBeFalse)
			})

			Convey("topo", func() {
				topo, err := dg.Flatten()
				So(err, ShouldBeNil)
				So(topo, ShouldResemble, []*Node{h, d, f, b, e, a, c, g})
			})
		})

		Convey("with order #2", func() {
			dg := NewDirectedGraph([]*Node{a, b, c, d, e, f, g, h})
			dg.
				AddEdge(e, g).
				AddEdge(d, e, f).
				AddEdge(c, g).
				AddEdge(b, c, e).
				AddEdge(a, c)

			Convey("cycle", func() {
				So(dg.isCyclic(), ShouldBeFalse)
			})

			Convey("topo", func() {
				topo, err := dg.Flatten()
				So(err, ShouldBeNil)
				So(topo, ShouldResemble, []*Node{h, d, f, b, e, a, c, g})
			})
		})
	})

	Convey("when error", t, func() {
		a := NewNode("a")
		b := NewNode("b")
		c := NewNode("c")
		d := NewNode("d")

		Convey("custom #1", func() {
			dg := NewDirectedGraph([]*Node{a, b, c, d})

			// a -> b -> c -> d
			// a -> c
			// c -> a
			// d -> d
			dg.
				AddEdge(a, b, c).
				AddEdge(b, c).
				AddEdge(c, d, a).
				AddEdge(d, d)

			Convey("cycle", func() {
				So(dg.isCyclic(), ShouldBeTrue)
			})

			Convey("topo", func() {
				topo, err := dg.Flatten()
				So(err, ShouldNotBeNil)
				So(topo, ShouldBeEmpty)
			})
		})

		Convey("custom #2", func() {
			dg := NewDirectedGraph([]*Node{a, b, c, d})

			// a -> b -> c -> d -> a
			dg.
				AddEdge(a, b).
				AddEdge(b, c).
				AddEdge(c, d).
				AddEdge(d, a)

			Convey("cycle", func() {
				So(dg.isCyclic(), ShouldBeTrue)
			})

			Convey("topo", func() {
				topo, err := dg.Flatten()
				So(err, ShouldNotBeNil)
				So(topo, ShouldBeEmpty)
			})
		})
	})
}
