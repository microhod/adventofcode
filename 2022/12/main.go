package main

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Hill Climbing Algorithm", part1, part2).Run()
}

func part1() error {
	graph, start, _, end, err := parse(InputFile)
	if err != nil {
		return err
	}

	path, err := graph.ShortestPath(start, end)
	if err != nil {
		return err
	}

	// print the number of steps (we can skip the start as that's step 0)
	fmt.Println(len(path) - 1)
	return nil
}

func part2() error {
	graph, _, possibleStarts, end, err := parse(InputFile)
	if err != nil {
		return err
	}

	shortestMu := &sync.RWMutex{}
	wg := &sync.WaitGroup{}

	shortest := math.MaxInt
	for start := range possibleStarts {
		wg.Add(1)
		go func(start string) {
			defer wg.Done()

			path, err := graph.ShortestPath(start, end)
			// skip starts which don't give valid paths
			if err != nil {
				return
			}

			shortestMu.Lock()
			// the number of steps (we can skip the start as that's step 0)
			if len(path)-1 < shortest {
				shortest = len(path) - 1
			}
			shortestMu.Unlock()
		}(start)
	}

	wg.Wait()
	fmt.Println(shortest)
	return nil
}

func parse(path string) (*Graph[string], string, set.Set[string], string, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, "", nil, "", err
	}

	graph := NewGraph[string]()
	possibleStarts := set.NewSet[string]()
	var start, end string

	for i, line := range lines {
		if line == "" {
			continue
		}
		for j, ch := range line {
			letter := parseLetter(ch)
			node := nodeName(i, j, letter)
			if ch == 'a' {
				possibleStarts.Add(node)
			}
			if ch == 'S' {
				start = node
			}
			if ch == 'E' {
				end = node
			}

			graph.AddNode(node)

			if i > 0 {
				otherLetter := parseLetter(rune(lines[i-1][j]))
				addEdge(graph, letter, otherLetter, node, nodeName(i-1, j, otherLetter))
			}
			if i < len(lines)-1 {
				otherLetter := parseLetter(rune(lines[i+1][j]))
				addEdge(graph, letter, otherLetter, node, nodeName(i+1, j, otherLetter))
			}
			if j > 0 {
				otherLetter := parseLetter(rune(lines[i][j-1]))
				addEdge(graph, letter, otherLetter, node, nodeName(i, j-1, otherLetter))
			}
			if j < len(lines[i])-1 {
				otherLetter := parseLetter(rune(lines[i][j+1]))
				addEdge(graph, letter, otherLetter, node, nodeName(i, j+1, otherLetter))
			}
		}
	}
	return graph, start, possibleStarts, end, nil
}

func parseLetter(letter rune) int {
	if letter == 'S' {
		letter = 'a'
	}
	if letter == 'E' {
		letter = 'z'
	}
	return int(letter)
}

func nodeName(i, j int, letter int) string {
	return fmt.Sprintf("(%d,%d)=%s", i, j, string(rune(letter)))
}

func addEdge(g *Graph[string], letter1, letter2 int, node1, node2 string) {
	if letter2 <= letter1+1 {
		g.AddEdge(node1, node2, 1)
	}
}

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
func (g *Graph[T]) ShortestPath(start, target T) ([]T, error) {
	_, previous := g.Dijkstra(start)

	return g.pathFromPrevious(start, target, previous)
}

func (g *Graph[T]) pathFromPrevious(start, target T, previous map[T]*T) ([]T, error) {
	var path []T
	seen := set.NewSet[T]()
	node := &target

	for node != nil {
		if seen.Contains(*node) {
			return nil, fmt.Errorf("path contains a cycle")
		}
		seen.Add(*node)

		path = append([]T{*node}, path...)
		node = previous[*node]
	}

	return path, nil
}

// Dijkstra is an implementation of Dijkstra's Algorithm
// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
func (g *Graph[T]) Dijkstra(start T) (map[T]int, map[T]*T) {
	distances := map[T]int{}
	previous := map[T]*T{}
	unvisited := NewPriorityQueue[T]()

	// initialise
	for _, node := range g.nodes {
		priority := math.MaxInt
		if node == start {
			priority = 0
		}

		distances[node] = priority
		unvisited.AddWithPriority(node, priority)
	}

	for !unvisited.Empty() {
		u := unvisited.ExtractMin()

		for v, dist := range g.edges[u] {
			alt := distances[u] + dist
			if alt < distances[v] {
				distances[v] = alt
				previous[v] = &u
				unvisited.DecreasePriority(v, alt)
			}
		}
	}

	return distances, previous
}

type PriorityQueue[T comparable] struct {
	priorities map[T]int
	queue      []T
}

func NewPriorityQueue[T comparable]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		priorities: map[T]int{},
	}
}

func (queue *PriorityQueue[T]) ExtractMin() T {
	var item T
	if len(queue.queue) < 1 {
		return item
	}

	item = queue.queue[0]
	queue.queue = queue.queue[1:]
	delete(queue.priorities, item)

	return item
}

func (queue *PriorityQueue[T]) Empty() bool {
	return len(queue.queue) == 0
}

func (queue *PriorityQueue[T]) AddWithPriority(item T, priority int) {
	queue.queue = append(queue.queue, item)
	queue.priorities[item] = priority
	queue.sort()
}

func (queue *PriorityQueue[T]) DecreasePriority(item T, priority int) {
	queue.priorities[item] = priority
	queue.sort()
}

func (queue *PriorityQueue[T]) sort() {
	// this is a bit crap given we only ever need to move one item, but it works!
	sort.Slice(queue.queue, func(i, j int) bool {
		return queue.priorities[queue.queue[i]] < queue.priorities[queue.queue[j]]
	})
}
