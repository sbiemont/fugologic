package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNode(t *testing.T) {
	Convey("new node", t, func() {
		Convey("when simple new", func() {
			So(NewNode(42), ShouldResemble, &Node{data: 42})
		})

		Convey("when complex new", func() {
			type hello struct {
				value float64
			}
			So(NewNode(hello{value: 42}), ShouldResemble, &Node{data: hello{value: 42}})
			So(NewNode(&hello{value: 42}), ShouldResemble, &Node{data: &hello{value: 42}})
		})
	})

	Convey("get data", t, func() {
		So(Node{data: 42}.Data(), ShouldEqual, 42)
		So(Node{data: nil}.Data(), ShouldBeNil)
	})
}

func TestEdges(t *testing.T) {
	a := NewNode("a")
	b := NewNode("b")
	c := NewNode("c")
	d := NewNode("d")
	e := NewNode("e")

	Convey("add", t, func() {
		edges := NewDirectedEdges()
		edges.Add(a, b).Add(a, c) // a -> b, c
		edges.Add(c, d, e)        // c -> d, e

		So(edges, ShouldResemble, DirectedEdges{
			a: []*Node{b, c},
			c: []*Node{d, e},
		})
	})
}

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
			dg, err := NewDirectedGraph([]*Node{a, b, c, d, e, f, g, h}, nil)

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, []*Node{h, g, f, e, d, c, b, a})
			})
		})

		Convey("when mini graph", func() {
			dg, err := NewDirectedGraph([]*Node{a, b, c}, DirectedEdges{
				a: []*Node{c},
				b: []*Node{c},
			})

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, []*Node{b, a, c})
			})
		})

		Convey("when 2 graphs", func() {
			// a, b -> c
			// d -> e, f
			// f -> e, g, h
			dg, err := NewDirectedGraph([]*Node{a, b, c, d, e, f, g, h}, DirectedEdges{
				a: []*Node{c},
				b: []*Node{c},
				d: []*Node{e, f},
				f: []*Node{e, g, h},
			})

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, []*Node{d, f, h, g, e, b, a, c})
			})
		})

		// a -> c
		// c -> g
		// b -> c, e
		// d -> e, f
		// e -> g
		Convey("with order #1", func() {
			dg, err := NewDirectedGraph([]*Node{a, b, c, d, e, f, g, h}, DirectedEdges{
				a: []*Node{c},
				b: []*Node{c, e},
				c: []*Node{g},
				d: []*Node{e, f},
				e: []*Node{g},
			})

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, []*Node{h, d, f, b, e, a, c, g})
			})
		})

		Convey("with order #2", func() {
			dg, err := NewDirectedGraph([]*Node{a, b, c, d, e, f, g, h}, DirectedEdges{
				e: []*Node{g},
				d: []*Node{e, f},
				c: []*Node{g},
				b: []*Node{c, e},
				a: []*Node{c},
			})

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
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
			// a -> b -> c -> d
			// a -> c
			// c -> a
			// d -> d
			dg, err := NewDirectedGraph([]*Node{a, b, c, d}, DirectedEdges{
				a: []*Node{b, c},
				b: []*Node{c},
				c: []*Node{d, a},
				d: []*Node{d},
			})

			Convey("cycle", func() {
				So(err, ShouldBeError, "cycle(s) detected in directed graph")
				So(dg, ShouldResemble, DirectedGraph{})
			})
		})

		Convey("custom #2", func() {
			// a -> b -> c -> d -> a
			dg, err := NewDirectedGraph([]*Node{a, b, c, d}, DirectedEdges{
				a: []*Node{b},
				b: []*Node{c},
				c: []*Node{d},
				d: []*Node{a},
			})

			Convey("cycle", func() {
				So(err, ShouldBeError, "cycle(s) detected in directed graph")
				So(dg, ShouldResemble, DirectedGraph{})
			})
		})
	})
}
