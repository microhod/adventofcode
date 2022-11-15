package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("All in a Single Night", part1, part2).Run()
}

func part1() error {
	nodes, g, err := parse(InputFile)
	if err != nil {
		return err
	}

	permutations := getPermutations(len(nodes))

	var min int
	for _, p := range permutations {
		var length int
		for i := 1; i < len(p); i++ {
			length += g[nodes[p[i-1]]][nodes[p[i]]]
		}
		if min == 0 || length < min {
			min = length
		}
	}
	fmt.Println("Shortest Distance:", min)

	return nil
}

func part2() error {
	nodes, g, err := parse(InputFile)
	if err != nil {
		return err
	}

	permutations := getPermutations(len(nodes))

	var max int
	for _, p := range permutations {
		var length int
		for i := 1; i < len(p); i++ {
			length += g[nodes[p[i-1]]][nodes[p[i]]]
		}
		if max == 0 || length > max {
			max = length
		}
	}
	fmt.Println("Longest Distance:", max)

	return nil
}

func parse(path string) ([]string, Graph, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, nil, err
	}

	g := Graph{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		parts := strings.Fields(line)
		if len(parts) < 5 {
			return nil, nil, fmt.Errorf("expected '<from> to <to> = <weight>' but got %s", line)
		}
		from, to, w := parts[0], parts[2], parts[4]

		weight, err := strconv.Atoi(w)
		if err != nil {
			return nil, nil, err
		}

		if g[from] == nil {
			g[from] = map[string]int{}
		}
		g[from][to] = weight
		if g[to] == nil {
			g[to] = map[string]int{}
		}
		g[to][from] = weight
	}

	return g.Nodes(), g, nil
}

type Graph map[string]map[string]int

func (g Graph) Nodes() []string {
	seen := map[string]bool{}
	var nodes []string

	for n, edges := range g {
		for m := range edges {
			if !seen[n] {
				nodes = append(nodes, n)
				seen[n] = true
			}
			if !seen[m] {
				nodes = append(nodes, m)
				seen[m] = true
			}
		}
	}

	return nodes
}

func getPermutations(n int) [][]int {
	var permutations [][]int

	partials := [][]int{}
	for i := 0; i < n; i++ {
		partials = append(partials, []int{i})
	}

	for len(partials) > 0 {
		partial := partials[0]
		partials = partials[1:]

		if len(partial) == n {
			permutations = append(permutations, partial)
		}

		for i := 0; i < n; i++ {
			if contains(partial, i) {
				continue
			}
			next := append([]int{}, partial...)
			next = append(next, i)
			partials = append(partials, next)
		}
	}
	return permutations
}

func contains(nums []int, num int) bool {
	for _, n := range nums {
		if n == num {
			return true
		}
	}
	return false
}
