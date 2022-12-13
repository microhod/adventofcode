package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/graph"
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

	path, _ := graph.DijkstraShortestPath(start, end)

	// print the number of steps (we can skip the start as that's step 0)
	fmt.Println(len(path) - 1)
	return nil
}

func part2() error {
	graph, _, possibleStarts, end, err := parse(InputFile)
	if err != nil {
		return err
	}

	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	shortest := math.MaxInt
	for start := range possibleStarts {
		wg.Add(1)
		go func(start string) {
			defer wg.Done()

			path, exists := graph.DijkstraShortestPath(start, end)
			// skip starts which don't have a path to the end
			if !exists {
				return
			}

			mu.Lock()
			// the number of steps (we can skip the start as that's step 0)
			if len(path)-1 < shortest {
				shortest = len(path) - 1
			}
			mu.Unlock()
		}(start)
	}

	wg.Wait()
	fmt.Println(shortest)
	return nil
}

func parse(path string) (*graph.Graph[string], string, set.Set[string], string, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, "", nil, "", err
	}

	g := graph.NewGraph[string]()
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

			g.AddNode(node)

			if i > 0 {
				otherLetter := parseLetter(rune(lines[i-1][j]))
				addEdge(g, letter, otherLetter, node, nodeName(i-1, j, otherLetter))
			}
			if i < len(lines)-1 {
				otherLetter := parseLetter(rune(lines[i+1][j]))
				addEdge(g, letter, otherLetter, node, nodeName(i+1, j, otherLetter))
			}
			if j > 0 {
				otherLetter := parseLetter(rune(lines[i][j-1]))
				addEdge(g, letter, otherLetter, node, nodeName(i, j-1, otherLetter))
			}
			if j < len(lines[i])-1 {
				otherLetter := parseLetter(rune(lines[i][j+1]))
				addEdge(g, letter, otherLetter, node, nodeName(i, j+1, otherLetter))
			}
		}
	}
	return g, start, possibleStarts, end, nil
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

func addEdge(g *graph.Graph[string], letter1, letter2 int, node1, node2 string) {
	if letter2 <= letter1+1 {
		g.AddEdge(node1, node2, 1)
	}
}
