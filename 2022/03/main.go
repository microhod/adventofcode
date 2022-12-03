package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Rucksack Reorganization", part1, part2).Run()
}

func part1() error {
	rucksacks, err := parse(InputFile)
	if err != nil {
		return err
	}

	var totalPriority int
	for _, rs := range rucksacks {
		for _, c := range rs.CommonItems() {
			totalPriority += priority(c)
		}
	}

	fmt.Printf("The total priority is: %d\n", totalPriority)
	return nil
}

func part2() error {
	rucksacks, err := parse(InputFile)
	if err != nil {
		return err
	}

	var totalPriority int
	for i := 2; i < len(rucksacks); i += 3 {
		common := Common(rucksacks[i-2], rucksacks[i-1], rucksacks[i])
		if len(common) != 1 {
			return fmt.Errorf("expected 1 common item but got %d", len(common))
		}
		totalPriority += priority(common[0])
	}

	fmt.Printf("The total priority is: %d\n", totalPriority)
	return nil
}

func parse(path string) ([]Rucksack, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var rs []Rucksack
	for _, line := range lines {
		if line == "" {
			continue
		}
		mid := len(line) / 2
		rs = append(rs, Rucksack{
			Compartment1: line[:mid],
			Compartment2: line[mid:],
		})
	}

	return rs, nil
}

type Rucksack struct {
	Compartment1 string
	Compartment2 string
}

func (rs Rucksack) CommonItems() []rune {
	set1 := set.NewSet([]rune(rs.Compartment1)...)
	set2 := set.NewSet([]rune(rs.Compartment2)...)

	var common []rune
	for r := range set.Intersect(set1, set2) {
		common = append(common, r)
	}
	return common
}

func Common(rucksacks ...Rucksack) []rune {
	var sets []set.Set[rune]

	for _, rs := range rucksacks {
		set := set.NewSet([]rune(rs.Compartment1+rs.Compartment2)...)
		sets = append(sets, set)
	}

	var common []rune
	for r := range set.Intersect(sets...) {
		common = append(common, r)
	}

	return common
}

func priority(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r-'a') + 1
	}
	if r >= 'A' && r <= 'Z' {
		return int(r-'A') + 27
	}
	panic(fmt.Errorf("invalid rune: %s", string(r)))
}
