package graph

import "github.com/microhod/adventofcode/internal/queue"

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

func (g *Graph[T]) DijkstraShortestPath(start, target T) ([]T, bool) {
	cost, cameFrom := g.Dijkstra(start)
	if _, exists := cost[target]; !exists {
		return nil, false
	}

	var path []T
	node := target

	for node != start {
		path = append([]T{node}, path...)
		node = cameFrom[node]
	}

	return append([]T{start}, path...), true
}

// Dijkstra is an implementation of Dijkstra's Algorithm
// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
//
// mainly implenented from this example:
// https://www.redblobgames.com/pathfinding/a-star/introduction.html#greedy-best-first 
func (g *Graph[T]) Dijkstra(start T) (map[T]int, map[T]T) {
	frontier := queue.NewPriorityQueue[T]()
	frontier.Put(start, 0)

	cameFrom := map[T]T{}
	costSoFar := map[T]int{}

	costSoFar[start] = 0

	for frontier.Size() > 0 {
		current := frontier.Get()

		for next, cost := range g.edges[current] {
			newCost := costSoFar[current] + cost

			if _, exists := costSoFar[next]; !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				frontier.Put(next, cost)
				cameFrom[next] = current
			}
		}
	}

	return costSoFar, cameFrom
}
