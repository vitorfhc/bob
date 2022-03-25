package graphite

import (
	"errors"

	"github.com/samber/lo"
)

// Node is a vertice of the graph.
type Node interface {
	comparable
}

// Graph is a graph representation.
type Graph[T Node] struct {
	connections [][]bool
	reversed    [][]bool
	reference   []T
	size        int
}

// New creates a graph using nodes as references.
func New[T Node](nodes []T) (*Graph[T], error) {
	uniqueNodes := lo.Uniq(nodes)
	if len(uniqueNodes) != len(nodes) {
		return nil, errors.New("can't create graph with duplicated nodes")
	}

	graph := &Graph[T]{}
	size := len(nodes)
	graph.size = size

	graph.reference = make([]T, size)
	copy(graph.reference, nodes)

	graph.connections = make([][]bool, size)
	graph.reversed = make([][]bool, size)
	for i := 0; i < size; i++ {
		graph.connections[i] = make([]bool, size)
		graph.reversed[i] = make([]bool, size)
	}

	return graph, nil
}

// DFS performs a depth-first search on the graph from the given node.
// It returns all nodes visited from the farthest to the nearest, inclusing the given node.
// It returns an error when cycles are detected or when the node doesn't exist.
func (graph *Graph[T]) DFS(node T, reverse bool) ([]T, error) {
	nodeIndex := lo.IndexOf(graph.reference, node)
	if nodeIndex == -1 {
		return nil, errors.New("node is not part of the graph")
	}

	visited := make([]bool, graph.size)
	results, err := dfs(nodeIndex, visited, graph, reverse)
	if err != nil {
		return nil, err
	}

	nodes := make([]T, len(results))
	lo.ForEach(results, func(neighbor int, i int) {
		nodes[i] = graph.reference[neighbor]
	})
	return nodes, nil
}

// GetDepsFor returns all nodes that depend on the given node
// in order of which should be processed first.
// You may receive an error if the node doesn't exist in the graph,
// or if any cycle if found.
// Example:
// 1 -> 2 -> 3
// GetDepsFor(3) returns [[1,2]];
// 1 -> 2 <- 3
// GetDepsFor(2) returns [[1],[3]];
// 1 2 3
// GetDepsFor(1) returns nil;
func (graph *Graph[T]) GetDepsFor(node T) ([][]T, error) {
	nodeIndex := lo.IndexOf(graph.reference, node)
	if nodeIndex == -1 {
		return nil, errors.New("node is not part of the graph")
	}

	results := make([][]T, 0)
	var err error
	lo.ForEach(graph.reversed[nodeIndex], func(connected bool, i int) {
		if connected {
			var deps []T
			deps, err = graph.DFS(graph.reference[i], true)
			if err != nil {
				return
			}
			results = append(results, deps)
		}
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}

// AddDependency explicits that a depends on b.
func (graph *Graph[T]) AddDependency(a T, b T) error {
	indexA := lo.IndexOf(graph.reference, a)
	indexB := lo.IndexOf(graph.reference, b)
	if indexA == -1 || indexB == -1 {
		return errors.New("node is not part of the graph")
	}

	connect(indexA, indexB, graph)

	return nil
}
