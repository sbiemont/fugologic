package graph

import (
	"errors"
)

// DirectedEdges represents a root node linked to a list of other nodes
type DirectedEdges[T any] map[*Node[T]]Nodes[T]

// NewDirectedEdges initialises a new instance of DirectedEdges
func NewDirectedEdges[T any]() DirectedEdges[T] {
	return make(map[*Node[T]]Nodes[T])
}

// Add a directed edge from a root node to a list of adjacent nodes
func (edges DirectedEdges[T]) Add(from *Node[T], to ...*Node[T]) DirectedEdges[T] {
	if len(to) == 0 {
		return edges
	}

	list := edges[from]
	list = append(list, to...)
	edges[from] = list
	return edges
}

// DirectedGraph represents a graph with directed edges
type DirectedGraph[T any] struct {
	edges DirectedEdges[T] // edges can be empty
	nodes Nodes[T]         // Keep an ordered list to get a deterministic result when flatten
}

// NewDirectedGraph initialises a new DirectedGraph instance
func NewDirectedGraph[T any](nodes Nodes[T], edges DirectedEdges[T]) (DirectedGraph[T], error) {
	dg := DirectedGraph[T]{
		nodes: nodes,
		edges: edges,
	}

	if dg.isCyclic() {
		return DirectedGraph[T]{}, errors.New("cycle(s) detected in directed graph")
	}

	return dg, nil
}

// color is for search of cycles in a graph
type color int8

const (
	white color = iota // white is for a node not yet processed
	grey               // grey is for a node being processed
	black              // black is for a node already processed
)

// dfs (for depth-first search) finds if a back edge exists in the sub-graph rooted with node `from`
func (dg DirectedGraph[T]) dfs(from *Node[T], colors map[*Node[T]]color) bool {
	// `from` is being processed
	colors[from] = grey

	for _, to := range dg.edges[from] {
		// If `to` is also grey => loop
		// If `to` is not processed and there is a back edge in subtree rooted with `to` => loop
		if colors[to] == grey || (colors[to] == white && dg.dfs(to, colors)) {
			return true
		}
	}

	// `from` is fully processed
	colors[from] = black
	return false
}

// isCyclic returns true if there is a cycle in graph
func (dg DirectedGraph[T]) isCyclic() bool {
	// Initialize colors
	colors := make(map[*Node[T]]color, len(dg.nodes))
	for _, node := range dg.nodes {
		colors[node] = white
	}

	// Do a DFS traversal beginning with all vertices
	for _, node := range dg.nodes {
		if colors[node] == white && dg.dfs(node, colors) {
			return true
		}
	}

	return false
}

// flatten recursively builds the sub-graph rooted with `from`
func (dg DirectedGraph[T]) flatten(from *Node[T], visited map[*Node[T]]bool, flat *Nodes[T]) {
	// Mark the current node as visited
	visited[from] = true

	// Recur for all the vertices adjacent to this node
	for _, to := range dg.edges[from] {
		if !visited[to] {
			dg.flatten(to, visited, flat)
		}
	}

	// Push current node into result
	*flat = append(*flat, from)
}

// Flatten the graph using a topological sort
func (dg DirectedGraph[T]) Flatten() Nodes[T] {
	// Init (all nodes are unvisited)
	visited := make(map[*Node[T]]bool, len(dg.nodes))
	for _, node := range dg.nodes {
		visited[node] = false
	}

	// Call the recursive to flatten all sub-graphs
	var flat Nodes[T]
	for _, node := range dg.nodes {
		if !visited[node] {
			dg.flatten(node, visited, &flat)
		}
	}

	// Reverse
	n := len(flat)
	result := make(Nodes[T], n)
	for i := 0; i < n; i++ {
		result[i] = flat[n-i-1]
	}
	return result
}

// Node is a link to a real object
type Node[T any] struct {
	data T
}

// NewNode builds a Node instance with data
func NewNode[T any](data T) *Node[T] {
	return &Node[T]{
		data: data,
	}
}

// Data returns the data contained in the node
func (node Node[T]) Data() T {
	return node.data
}

// Nodes is a flat list
type Nodes[T any] []*Node[T]

// NewNodes creates a list of nodes using the ordered input list
func NewNodes[T any](input []T) Nodes[T] {
	result := make(Nodes[T], len(input))
	for i, in := range input {
		result[i] = NewNode(in)
	}
	return result
}

// Data extract an ordered list of nodes' data
func (nodes Nodes[T]) Data() []T {
	result := make([]T, len(nodes))
	for i, node := range nodes {
		result[i] = node.Data()
	}
	return result
}
