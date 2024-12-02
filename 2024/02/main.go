package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("RedNosed Reports", part1, part2).Run()
}

func part1() error {
	reports, err := parse(InputFile)
	if err != nil {
		return err
	}

	var safe int
	for _, report := range reports {
		if report.Safe() {
			safe++
		}
	}

	fmt.Printf("safe reports %d\n", safe)
	return nil
}

func part2() error {
	reports, err := parse(InputFile)
	if err != nil {
		return err
	}

	var safe int
	for _, report := range reports {
		if report.SafeWithProblemDampener() {
			safe++
		}
	}

	fmt.Printf("safe reports with problem dampener %d\n", safe)
	return nil
}

func parse(path string) ([]Report, error) {
	nums, err := file.ReadAllCsvInts(path, " ")
	if err != nil {
		return nil, err
	}

	var reports []Report
	for _, r := range nums {
		reports = append(reports, Report(r))
	}
	return reports, nil
}

type Report []int

func (r Report) Safe() bool {
	less := func(a, b int) bool {
		return a < b
	}
	if r[0] > r[1] {
		less = func(a, b int) bool {
			return a > b
		}
	}

	for i := range len(r)-1 {
		if !less(r[i], r[i+1]) {
			return false
		}
		if maths.Abs(r[i+1] - r[i]) > 3 {
			return false
		}
	}
	return true
}

func (r Report) SafeWithProblemDampener() bool {
	for i := range r {
		if Report(append(append([]int{}, r[0:i]...), r[i+1:]...)).Safe() {
			return true
		}
	}
	return false
}
