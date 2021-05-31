package graph

import (
	"errors"
)

// DirectedGraph represents a graph with directed edges
type DirectedGraph struct {
	edges map[*Node][]*Node
	nodes []*Node // Keep an ordered list to get a deterministic result when flatten
}

// NewDirectedGraph initialises a new DirectedGraph instance
func NewDirectedGraph(nodes []*Node) DirectedGraph {
	return DirectedGraph{
		nodes: nodes,
		edges: make(map[*Node][]*Node),
	}
}

// AddEdge adds a directed edge from `n1` to a list of adjacent nodes
func (dg *DirectedGraph) AddEdge(n1 *Node, adj ...*Node) *DirectedGraph {
	if len(adj) == 0 {
		return dg
	}

	list := dg.edges[n1]
	list = append(list, adj...)
	dg.edges[n1] = list
	return dg
}

// color is for search of cycles in a graph
type color int8

const (
	white color = iota // white is for a node not yet processed
	grey               // grey is for a node being processed
	black              // black is for a node already processed
)

// dfs (for depth-first search) finds if a back edge exists in the sub-graph rooted with node `from`
func (dg DirectedGraph) dfs(from *Node, colors map[*Node]color) bool {
	// `from` is being processed
	colors[from] = grey

	for _, to := range dg.edges[from] {
		// If `to`is also grey => loop
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
func (dg DirectedGraph) isCyclic() bool {
	// Initialize colors
	colors := make(map[*Node]color, len(dg.nodes))
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
func (dg DirectedGraph) flatten(from *Node, visited map[*Node]bool, flat *[]*Node) {
	// Mark the current node as visited.
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
func (dg DirectedGraph) Flatten() ([]*Node, error) {
	if dg.isCyclic() {
		return nil, errors.New("cannot flatten graph with cycle(s)")
	}

	// Init (all nodes are unvisited)
	visited := make(map[*Node]bool, len(dg.nodes))
	for _, node := range dg.nodes {
		visited[node] = false
	}

	// Call the recursive to flatten all sub-graphs
	var flat []*Node
	for _, node := range dg.nodes {
		if !visited[node] {
			dg.flatten(node, visited, &flat)
		}
	}

	// Reverse
	n := len(flat)
	result := make([]*Node, n)
	for i := 0; i < n; i++ {
		result[i] = flat[n-i-1]
	}
	return result, nil
}

// Node is a link to a real object
type Node struct {
	data interface{}
}

// NewNode builds a Node instance with data
func NewNode(data interface{}) *Node {
	return &Node{
		data: data,
	}
}

// Data returns the data contained in the node
func (node Node) Data() interface{} {
	return node.data
}
