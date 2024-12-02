package main

import (
	"fmt"
	"slices"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Historian Hysteria", part1, part2).Run()
}

func part1() error {
	left, right, err := parse(InputFile)
	if err != nil {
		return err
	}

	slices.Sort(left)
	slices.Sort(right)

	var diff int
	for i := range left {
		diff += int(maths.Abs(left[i] - right[i]))
	}
	
	fmt.Printf("total diff: %d\n", diff)
	return nil
}

func part2() error {
	left, right, err := parse(InputFile)
	if err != nil {
		return err
	}

	occurences := make(map[int]int)
	for _, id := range right {
		occurences[id]++
	}

	var similarity int
	for _, id := range left {
		similarity += id * occurences[id]
	}

	fmt.Printf("similarity score: %d\n", similarity)
	return nil
}

func parse(path string) ([]int, []int, error) {
	nums, err := file.ReadAllCsvInts(path, "   ")
	if err != nil {
		return nil, nil, err
	}

	var left, right []int
	for _, row := range nums {
		left = append(left, row[0])
		right = append(right, row[1])
	}

	return left, right, nil
}
