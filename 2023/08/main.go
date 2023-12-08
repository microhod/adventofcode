package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Haunted Wasteland", part1, part2).Run()
}

func part1() error {
	directions, network, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	steps := network.FollowDirections("AAA", "ZZZ", directions...)

	fmt.Printf("the steps from AAA to ZZZ is: %d\n", steps)
	return nil
}

func part2() error {
	directions, network, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	steps := network.FollowGhostDirections(
		func(s string) bool {
			return s[2] == 'A'
		},
		func(s string) bool {
			return s[2] == 'Z'
		},
		directions...,
	)

	fmt.Printf("the steps from AAA to ZZZ is: %d\n", steps)
	return nil
}

func parseInput(path string) ([]Direction, Network, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, nil, err
	}

	parts := strings.Split(input, "\n\n")
	directions := []Direction(strings.TrimSpace(parts[0]))

	network := make(Network)
	for _, line := range strings.Split(parts[1], "\n") {
		if line == "" {
			continue
		}

		line = strings.ReplaceAll(line, " = (", ", ")
		line = strings.ReplaceAll(line, ")", "")
		parts = strings.Split(line, ", ")

		network[parts[0]] = map[Direction]string{
			'L': parts[1],
			'R': parts[2],
		}
	}

	return directions, network, nil
}

type Network map[string]map[Direction]string

type Direction rune

func (n Network) FollowDirections(start, end string, directions ...Direction) int {
	node := start
	steps := 0
	for node != end {
		direction := directions[maths.Mod(steps, len(directions))]
		node = n[node][direction]
		steps++
	}
	return steps
}

func (n Network) FollowGhostDirections(start, end func(string) bool, directions ...Direction) int {
	nodes := set.NewSet[string]()
	ends := set.NewSet[string]()
	for node := range n {
		if start(node) {
			nodes.Add(node)
		}
		if end(node) {
			ends.Add(node)
		}
	}

	steps := 0
	var finished []int

	for len(nodes) > 0 {
		direction := directions[maths.Mod(steps, len(directions))]
		steps++

		var nexts []string
		for node := range nodes {
			next := n[node][direction]
			if ends.Contains(next) {
				finished = append(finished, steps)
				continue
			}
			nexts = append(nexts, next)
		}
		nodes = set.NewSet[string](nexts...)
	}

	return maths.Lcm(finished...)
}
