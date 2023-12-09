package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Mirage Maintenance", part1, part2).Run()
}

func part1() error {
	sequences, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	var sum int
	for _, s := range sequences {
		sum += s.ExtrapolateNext()
	}

	fmt.Printf("the sum of all extrapolations is: %d\n", sum)
	return nil
}

func part2() error {
	sequences, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	var sum int
	for _, s := range sequences {
		sum += s.ExtrapolatePrevious()
	}

	fmt.Printf("the sum of all backwards extrapolations is: %d\n", sum)
	return nil
}

func parseInput(path string) ([]Sequence, error) {
	nums, err := file.ReadAllCsvInts(path, " ")
	if err != nil {
		return nil, err
	}

	var sequences []Sequence
	for _, row := range nums {
		sequences = append(sequences, Sequence(row))
	}
	return sequences, nil
}

type Sequence []int

func (s Sequence) ExtrapolateNext() int {
	sequences := s.AllDiffs()

	// add extra zero to end of 'all zeros' sequence
	sequences[len(sequences)-1] = append(sequences[len(sequences)-1], 0)
	// extrapolate each layer until we get to the original sequence
	for i := len(sequences)-2; i >= 0; i-- {
		extrapolation := sequences[i].Last() + sequences[i+1].Last()
		sequences[i] = append(sequences[i], extrapolation)
	}

	return sequences[0].Last()
}

func (s Sequence) ExtrapolatePrevious() int {
	sequences := s.AllDiffs()

	// add extra zero to start of 'all zeros' sequence
	sequences[len(sequences)-1] = append(Sequence{0}, sequences[len(sequences)-1]...)
	// extrapolate each layer until we get to the original sequence
	for i := len(sequences)-2; i >= 0; i-- {
		extrapolation := sequences[i].First() - sequences[i+1].First()
		sequences[i] = append(Sequence{extrapolation}, sequences[i]...)
	}

	return sequences[0].First()
}

func (s Sequence) AllDiffs() []Sequence {
	diffs := []Sequence{s}
	allZeros := false

	for !allZeros {
		var next Sequence
		next, allZeros = diffs[len(diffs)-1].Diff()
		diffs = append(diffs, next)
	}
	return diffs
}

func (s Sequence) Diff() (Sequence, bool) {
	diff := make(Sequence, len(s)-1)
	allZeros := true
	for i := 1; i < len(s); i++ {
		diff[i-1] = s[i]-s[i-1]
		allZeros = allZeros && diff[i-1] == 0
	}
	return diff, allZeros
}

func (s Sequence) First() int {
	return s[0]
}

func (s Sequence) Last() int {
	return s[len(s)-1]
}
