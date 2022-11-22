package main

import (
	"fmt"
	"strconv"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("No Such Thing as Too Much", part1, part2).Run()
}

func part1() error {
	containers, err := parse(InputFile)
	if err != nil {
		return err
	}

	combinations := combinations(containers, 150)
	fmt.Printf("Combinations of containers that hold 150 litres of eggnog: %d\n", len(combinations))

	return nil
}

func part2() error {
	containers, err := parse(InputFile)
	if err != nil {
		return err
	}
	
	var min int
	var count int
	for _, combination := range combinations(containers, 150) {
		if min == 0 || len(combination) < min {
			min = len(combination)
			count = 0
		}
		if len(combination) == min {
			count++
		}
	}
	fmt.Printf("Minimum number of containers required to hold 150 litres of eggnog: %d\n", min)
	fmt.Printf("Combinations which satisfy that minimum: %d\n", count)
	return nil
}

func parse(path string) ([]int, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var containers []int
	for _, line := range lines {
		if line == "" {
			continue
		}
		c, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		containers = append(containers, c)
	}
	return containers, nil
}

func combinations(containers []int, target int) [][]int {
	type combination struct {
		containers []int
		left       int
	}

	var stack []combination
	for i, capacity := range containers {
		stack = append(stack, combination{
			containers: []int{i},
			left:       target - capacity,
		})
	}

	var combinations [][]int
	for len(stack) > 0 {
		partial := stack[0]
		stack = stack[1:]

		if partial.left < 0 {
			continue
		}
		if partial.left == 0 {
			combinations = append(combinations, partial.containers)
			continue
		}

		lastIndex := partial.containers[len(partial.containers)-1]

		for i := lastIndex + 1; i < len(containers); i++ {
			c := combination{
				containers: append([]int{}, partial.containers...),
			}
			c.containers = append(c.containers, i)
			c.left = partial.left - containers[i]

			stack = append(stack, c)
		}
	}

	return combinations
}
