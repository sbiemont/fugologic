package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNode(t *testing.T) {
	Convey("new node", t, func() {
		Convey("when simple new", func() {
			So(NewNode(42), ShouldResemble, &Node[int]{data: 42})
		})

		Convey("when complex new", func() {
			type hello struct {
				value float64
			}
			So(NewNode(hello{value: 42}), ShouldResemble, &Node[hello]{data: hello{value: 42}})
			So(NewNode(&hello{value: 42}), ShouldResemble, &Node[*hello]{data: &hello{value: 42}})
		})
	})

	Convey("get data", t, func() {
		So(Node[int]{data: 42}.Data(), ShouldEqual, 42)
		So(Node[*int]{data: nil}.Data(), ShouldBeNil)
	})
}

func TestEdges(t *testing.T) {
	a := NewNode("a")
	b := NewNode("b")
	c := NewNode("c")
	d := NewNode("d")
	e := NewNode("e")

	Convey("add", t, func() {
		edges := NewDirectedEdges[string]()
		edges.Add(a, b).Add(a, c) // a -> b, c
		edges.Add(c, d, e)        // c -> d, e
		edges.Add(c)              // c -> nil (no edge)

		So(edges, ShouldResemble, DirectedEdges[string]{
			a: Nodes[string]{b, c},
			c: Nodes[string]{d, e},
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

		Convey("when no edge (all independent nodes", func() {
			dg, err := NewDirectedGraph(Nodes[string]{a, b, c, d, e, f, g, h}, nil)

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, Nodes[string]{h, g, f, e, d, c, b, a})
			})
		})

		Convey("when mini graph", func() {
			dg, err := NewDirectedGraph(Nodes[string]{a, b, c}, DirectedEdges[string]{
				a: Nodes[string]{c},
				b: Nodes[string]{c},
			})

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, Nodes[string]{b, a, c})
			})
		})

		Convey("when 2 graphs", func() {
			// a, b -> c
			// d -> e, f
			// f -> e, g, h
			dg, err := NewDirectedGraph(Nodes[string]{a, b, c, d, e, f, g, h}, DirectedEdges[string]{
				a: Nodes[string]{c},
				b: Nodes[string]{c},
				d: Nodes[string]{e, f},
				f: Nodes[string]{e, g, h},
			})

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, Nodes[string]{d, f, h, g, e, b, a, c})
			})
		})

		// a -> c
		// c -> g
		// b -> c, e
		// d -> e, f
		// e -> g
		Convey("with order #1", func() {
			dg, err := NewDirectedGraph(Nodes[string]{a, b, c, d, e, f, g, h}, DirectedEdges[string]{
				a: Nodes[string]{c},
				b: Nodes[string]{c, e},
				c: Nodes[string]{g},
				d: Nodes[string]{e, f},
				e: Nodes[string]{g},
			})

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, Nodes[string]{h, d, f, b, e, a, c, g})
			})
		})

		Convey("with order #2", func() {
			dg, err := NewDirectedGraph(Nodes[string]{a, b, c, d, e, f, g, h}, DirectedEdges[string]{
				e: Nodes[string]{g},
				d: Nodes[string]{e, f},
				c: Nodes[string]{g},
				b: Nodes[string]{c, e},
				a: Nodes[string]{c},
			})

			Convey("cycle", func() {
				So(err, ShouldBeNil)
			})

			Convey("topo", func() {
				topo := dg.Flatten()
				So(topo, ShouldResemble, Nodes[string]{h, d, f, b, e, a, c, g})
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
			dg, err := NewDirectedGraph(Nodes[string]{a, b, c, d}, DirectedEdges[string]{
				a: Nodes[string]{b, c},
				b: Nodes[string]{c},
				c: Nodes[string]{d, a},
				d: Nodes[string]{d},
			})

			Convey("cycle", func() {
				So(err, ShouldBeError, "cycle(s) detected in directed graph")
				So(dg, ShouldResemble, DirectedGraph[string]{})
			})
		})

		Convey("custom #2", func() {
			// a -> b -> c -> d -> a
			dg, err := NewDirectedGraph(Nodes[string]{a, b, c, d}, DirectedEdges[string]{
				a: Nodes[string]{b},
				b: Nodes[string]{c},
				c: Nodes[string]{d},
				d: Nodes[string]{a},
			})

			Convey("cycle", func() {
				So(err, ShouldBeError, "cycle(s) detected in directed graph")
				So(dg, ShouldResemble, DirectedGraph[string]{})
			})
		})
	})
}

func TestNodes(t *testing.T) {
	Convey("when ok", t, func() {
		nodes := Nodes[string]{
			NewNode("a"),
			NewNode("b"),
			NewNode("c"),
		}
		So(nodes.Data(), ShouldResemble, []string{"a", "b", "c"})
	})

	Convey("when empty", t, func() {
		So(Nodes[int]{}.Data(), ShouldResemble, []int{})
	})
}
