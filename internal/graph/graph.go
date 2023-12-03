package graph

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/microhod/adventofcode/internal/queue"
)

type Graph[T comparable] map[T]map[T]int

func NewGraph[T comparable]() Graph[T] {
	return Graph[T]{}
}

func (g Graph[T]) AddEdge(u, v T, weight int) {
	if g[u] == nil {
		g[u] = map[T]int{}
	}

	g[u][v] = weight
}

func (g Graph[T]) DijkstraShortestPath(start, target T) ([]T, int, bool) {
	cost, cameFrom := g.Dijkstra(start)
	if _, exists := cost[target]; !exists {
		return nil, 0, false
	}

	var path []T
	node := target

	for node != start {
		path = append([]T{node}, path...)
		node = cameFrom[node]
	}

	return append([]T{start}, path...), cost[target], true
}

// Dijkstra is an implementation of Dijkstra's Algorithm
// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
//
// mainly implenented from this example:
// https://www.redblobgames.com/pathfinding/a-star/introduction.html#greedy-best-first
func (g Graph[T]) Dijkstra(start T) (map[T]int, map[T]T) {
	frontier := queue.NewPriorityQueue[T]()
	frontier.Push(start, 0)

	cameFrom := map[T]T{}
	costSoFar := map[T]int{}

	costSoFar[start] = 0

	for frontier.Size() > 0 {
		current := frontier.Pop()

		for next, cost := range g[current] {
			newCost := costSoFar[current] + cost

			if _, exists := costSoFar[next]; !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				frontier.Push(next, cost)
				cameFrom[next] = current
			}
		}
	}

	return costSoFar, cameFrom
}

func (g Graph[T]) Mermaid() string {
	builder := &strings.Builder{}
	builder.WriteString("graph LR")

	label := func(node T) string {
		l := fmt.Sprint(node)
		// mermaid labels cannot contain whitespace
		l = regexp.MustCompile(`\s+`).ReplaceAllString(l, "")
		return l
	}

	for from := range g {
		for to, weight := range g[from] {
			edge := fmt.Sprintf("\n  %s --%d--> %s", label(from), weight, label(to))
			builder.WriteString(edge)
		}
	}
	return builder.String()
}
