package main

import (
	"fmt"
	"strconv"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	ReportFile = "sonar-sweep-report.txt"
	TestFile   = "test-report.txt"
)

func main() {
	puzzle.NewSolution("Sonar Sweep", part1, part2).Run()
}

func part1() error {
	var err error
	report, err := readReport(ReportFile)
	if err != nil {
		return err
	}

	fmt.Printf("# depth increases: %d\n", numIncreases(report))
	return nil
}

func part2() error {
	var err error
	report, err := readReport(ReportFile)
	if err != nil {
		return err
	}

	fmt.Printf("# depth increases (three-measurement sliding window): %d\n", numIncreases(slidingWindow(report, 3)))
	return nil
}

func numIncreases(report []int) int {
	if len(report) < 2 {
		return len(report)
	}

	increases := 0
	prev := report[0]
	for _, depth := range report[1:] {
		if depth > prev {
			increases += 1
		}
		prev = depth
	}

	return increases
}

func slidingWindow(report []int, size int) []int {
	if len(report) < size {
		return []int{}
	}

	slidingReport := []int{}
	for i := range report[:len(report)-(size-1)] {
		sum := 0
		for _, depth := range report[i:i+size] {
			sum += depth
		}
		slidingReport = append(slidingReport, sum)
	}

	return slidingReport
}

func readReport(path string) ([]int, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	report := []int{}
	for _, line := range(lines) {
		depth, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}

		report = append(report, depth)
	}

	return report, nil
}
