package graph_test

import (
	"slices"
	"testing"

	"github.com/sbiemont/fugologic/graph"
	. "github.com/smartystreets/goconvey/convey"
)

type node struct {
	id string
}

// ShouldBeOrdered checks if the given elements are ordered in the slice
func ShouldBeOrdered[T comparable](s []T, items ...T) bool {
	lastPos := -1
	for _, item := range items {
		pos := slices.Index(s, item)
		if pos == -1 {
			return false
		}
		if pos < lastPos {
			return false
		}
		lastPos = pos
	}
	return true
}

func TestEdges(t *testing.T) {
	a := node{id: "a"}
	b := node{id: "b"}
	c := node{id: "c"}
	d := node{id: "d"}
	e := node{id: "e"}

	Convey("add", t, func() {
		edges := graph.New[node]()
		edges.Add(a, b).Add(a, c) // a -> b, c
		edges.Add(c, d, e)        // c -> d, e
		edges.Add(c)              // c -> nil (no edge)

		So(edges, ShouldResemble, graph.Graph[node]{
			a: []node{b, c},
			c: []node{d, e},
		})
	})
}

func TestGraph(t *testing.T) {
	Convey("when ok #1", t, func() {
		a := node{id: "a"}
		b := node{id: "b"}
		c := node{id: "c"}
		d := node{id: "d"}
		e := node{id: "e"}
		f := node{id: "f"}
		g := node{id: "g"}
		h := node{id: "h"}

		Convey("when no edge (all independent nodes", func() {
			dg := graph.Graph[node]{a: nil, b: nil, c: nil, d: nil, e: nil, f: nil, g: nil, h: nil}

			topo, err := dg.TopologicalSort()
			So(err, ShouldBeNil)
			So(topo, ShouldHaveLength, 8) // no order
		})

		Convey("when mini graph", func() {
			dg := graph.Graph[node]{
				a: []node{c},
				b: []node{c},
			}

			topo, err := dg.TopologicalSort()
			So(err, ShouldBeNil)
			So(ShouldBeOrdered(topo, a, c), ShouldBeTrue)
			So(ShouldBeOrdered(topo, b, c), ShouldBeTrue)
		})

		Convey("when 2 graphs", func() {
			// a, b -> c
			// d -> e, f
			// f -> e, g, h
			dg := graph.Graph[node]{
				a: []node{c},
				b: []node{c},
				d: []node{e, f},
				f: []node{e, g, h},
			}

			topo, err := dg.TopologicalSort()
			So(err, ShouldBeNil)
			So(ShouldBeOrdered(topo, a, c), ShouldBeTrue)
			So(ShouldBeOrdered(topo, b, c), ShouldBeTrue)
			So(ShouldBeOrdered(topo, d, e), ShouldBeTrue)
			So(ShouldBeOrdered(topo, d, f), ShouldBeTrue)
			So(ShouldBeOrdered(topo, f, e), ShouldBeTrue)
			So(ShouldBeOrdered(topo, f, g), ShouldBeTrue)
			So(ShouldBeOrdered(topo, f, h), ShouldBeTrue)
		})

		// a -> c
		// c -> g
		// b -> c, e
		// d -> e, f
		// e -> g
		Convey("with order", func() {
			dg := graph.Graph[node]{
				a: []node{c},
				b: []node{c, e},
				c: []node{g},
				d: []node{e, f},
				e: []node{g},
			}

			topo, err := dg.TopologicalSort()
			So(err, ShouldBeNil)
			So(ShouldBeOrdered(topo, a, c), ShouldBeTrue)
			So(ShouldBeOrdered(topo, b, c, g), ShouldBeTrue)
			So(ShouldBeOrdered(topo, b, e, g), ShouldBeTrue)
			So(ShouldBeOrdered(topo, d, e, g), ShouldBeTrue)
			So(ShouldBeOrdered(topo, d, f), ShouldBeTrue)
		})
	})

	Convey("when error", t, func() {
		a := node{id: "a"}
		b := node{id: "b"}
		c := node{id: "c"}
		d := node{id: "d"}

		Convey("custom #1", func() {
			// a -> b -> c -> d
			// a -> c
			// c -> a
			// d -> d
			dg := graph.Graph[node]{
				a: []node{b, c},
				b: []node{c},
				c: []node{d, a},
				d: []node{d},
			}
			topo, err := dg.TopologicalSort()
			So(err, ShouldBeError, "cycle detected")
			So(topo, ShouldBeEmpty)
		})

		Convey("custom #2", func() {
			// a -> b -> c -> d -> a
			dg := graph.Graph[node]{
				a: []node{b},
				b: []node{c},
				c: []node{d},
				d: []node{a},
			}
			topo, err := dg.TopologicalSort()
			So(err, ShouldBeError, "cycle detected")
			So(topo, ShouldBeEmpty)
		})
	})
}
