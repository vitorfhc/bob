package graphite

import (
	"errors"

	"github.com/samber/lo"
)

func dfs[T Node](node int, visited []bool, graph *Graph[T], reverse bool) ([]int, error) {
	if visited[node] {
		return nil, errors.New("cycle detected while traversing graph")
	}
	visited[node] = true

	var connections []bool
	if reverse {
		connections = graph.reversed[node]
	} else {
		connections = graph.connections[node]
	}

	results := make([]int, 0)
	var err error
	lo.ForEach(connections, func(connected bool, i int) {
		if connected {
			var dfsResult []int
			dfsResult, err = dfs(i, visited, graph, reverse)
			if err != nil {
				return
			}
			results = append(results, dfsResult...)
		}
	})
	if err != nil {
		return nil, err
	}
	results = append(results, node)
	return results, nil
}

func connect[T Node](from, to int, graph *Graph[T]) {
	graph.reversed[to][from] = true
	graph.connections[from][to] = true
}
