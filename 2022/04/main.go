package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Camp Cleanup", part1, part2).Run()
}

func part1() error {
	pairs, err := parse(InputFile)
	if err != nil {
		return err
	}

	var contained int
	for _, pair := range pairs {
		if pair[0].Contains(pair[1]) || pair[1].Contains(pair[0]) {
			contained += 1
		}
	}

	fmt.Printf("total contained: %d\n", contained)
	return nil
}

func part2() error {
	pairs, err := parse(InputFile)
	if err != nil {
		return err
	}

	var overlaps int
	for _, pair := range pairs {
		if pair[0].Overlaps(pair[1]) {
			overlaps += 1
		}
	}

	fmt.Printf("total overlaps: %d\n", overlaps)
	return nil
}

func parse(path string) ([]ElfPair, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var pairs []ElfPair
	for _, line := range lines {
		// replace dashes to make each line csv format
		line = strings.ReplaceAll(line, "-", ",")
		nums, err := csv.ParseInts(line)
		if err != nil {
			return nil, err
		}
		
		pairs = append(pairs, ElfPair{
			Range{nums[0], nums[1]},
			Range{nums[2], nums[3]},
		})
	}

	return pairs, nil
}

type ElfPair [2]Range

type Range [2]int

func (r Range) Contains(s Range) bool {
	return s[0] >= r[0] && s[1] <= r[1]
}

func (r Range) Overlaps(s Range) bool {
	sStartsInsideR := s[0] >= r[0] && s[0] <= r[1]
	rStartsInsideS := r[0] >= s[0] && r[0] <= s[1]

	return sStartsInsideR || rStartsInsideS
}
