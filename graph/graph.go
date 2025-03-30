package graph

import (
	"fmt"
	"slices"
)

// Graph represents a root node linked to a list of other nodes
type Graph[T comparable] map[T][]T

// New initializes a new instance of Graph
func New[T comparable]() Graph[T] {
	return make(map[T][]T)
}

// Add a directed edge from a root node to a list of adjacent nodes
func (g Graph[T]) Add(from T, to ...T) Graph[T] {
	if len(to) == 0 {
		return g
	}

	list := g[from]
	list = append(list, to...)
	g[from] = list
	return g
}

var (
	ErrCyclicGraph = fmt.Errorf("cycle detected")
)

// TopologicalSort performs a topological sort
// Returns an error if a cycle is detected
func (g Graph[T]) TopologicalSort() ([]T, error) {
	var result []T
	err := runDFS(g, func(node T) error {
		result = append(result, node)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Reverse the result to get the correct topological order
	slices.Reverse(result)
	return result, nil
}

// IsCyclic checks if the graph contains a loop
// * edges: list of directed edges from one node to a list of nodes
// func IsCyclic[T comparable](edges map[T][]T) bool {
// 	return DFS(edges, nil) != nil
// }

// DFS performs a depth-first search on the graph represented by edges
// process function when a node is reached
func runDFS[T comparable](edges map[T][]T, process func(T) error) error {
	// Fetch all nodes and init them with white color (2 operations at once)
	coloredNodes := make(map[T]color)
	for n1, edge := range edges {
		coloredNodes[n1] = white
		for _, n2 := range edge {
			coloredNodes[n2] = white
		}
	}

	// Do a DFS traversal on each node
	for node := range coloredNodes { // do not use the direct color, it will be updated later
		if coloredNodes[node] == white {
			err := dfs(edges, node, coloredNodes, process)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// color is for search of cycles in a graph
type color int8

const (
	white color = iota // white is for a node not yet processed
	grey               // grey is for a node being processed
	black              // black is for a node already processed
)

// dfs (for depth-first search) finds if a back edge exists in the sub-graph rooted with node "from"
func dfs[T comparable](edges map[T][]T, from T, colors map[T]color, process func(T) error) error {
	// "from" is being processed
	colors[from] = grey

	for _, to := range edges[from] {
		// If "to" is also grey => cycle detected
		if colors[to] == grey {
			return ErrCyclicGraph
		}
		// If "to" is not processed and there is a back edge in subtree rooted with "to" => loop
		if colors[to] == white {
			err := dfs(edges, to, colors, process)
			if err != nil {
				return err
			}
		}
	}

	// "from" is fully processed
	colors[from] = black
	var err error
	if process != nil {
		err = process(from)
	}
	return err
}
