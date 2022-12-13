package graph

import (
	"math"

	"github.com/microhod/adventofcode/internal/queue"
)

type Graph[T comparable] struct {
	nodes []T
	edges map[T]map[T]int
}

func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{
		edges: map[T]map[T]int{},
	}
}

func (g *Graph[T]) AddNode(node T) {
	g.nodes = append(g.nodes, node)
	g.edges[node] = map[T]int{}
}

func (g *Graph[T]) AddEdge(u, v T, weight int) {
	if g.edges[u] == nil {
		g.edges[u] = map[T]int{}
	}

	g.edges[u][v] = weight
}

// ShortestPath is an implementation of Dijkstra's Algorithm
// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
func (g *Graph[T]) DijkstraShortestPath(start, target T) ([]T, bool) {
	dist, previous := g.Dijkstra(start)
	if dist[target] == math.MaxInt {
		return nil, false
	}

	var path []T
	node := target

	for node != start {
		path = append([]T{node}, path...)
		node = previous[node]
	}

	return append([]T{start}, path...), true
}

// Dijkstra is an implementation of Dijkstra's Algorithm
// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
func (g *Graph[T]) Dijkstra(start T) (map[T]int, map[T]T) {
	distances := map[T]int{}
	previous := map[T]T{}
	unvisited := queue.NewPriorityQueue[T]()

	// initialise
	for _, node := range g.nodes {
		distance := math.MaxInt
		if node == start {
			distance = 0
		}

		distances[node] = distance
		unvisited.AddWithPriority(node, distance)
	}

	for !unvisited.Empty() {
		u := unvisited.ExtractMin()
		// give up once we get to a node with infinite distance
		if distances[u] == math.MaxInt {
			break
		}

		for v, dist := range g.edges[u] {
			alt := distances[u] + dist
			if alt < distances[v] {
				distances[v] = alt
				previous[v] = u
				unvisited.DecreasePriority(v, alt)
			}
		}
	}

	return distances, previous
}
