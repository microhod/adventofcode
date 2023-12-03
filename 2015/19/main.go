package main

import (
	"fmt"
	"strings"

	"github.com/agnivade/levenshtein"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/queue"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Medicine for Rudolph", part1, part2).Run()
}

func part1() error {
	replacements, molecule, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	distinct := molecule.PossibleMolecules(replacements)
	fmt.Printf("Number of possible distinct molecules after one replacement: %d\n", len(distinct))
	return nil
}

func part2() error {
	replacements, molecule, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	// this in theory should work, but it takes way too long at the moment...
	steps := FindShortestPath("e", molecule, replacements)
	fmt.Printf("fewest number of steps to get molecule: %d\n", steps)
	return nil
}

func parseInput(path string) ([]Replacement, Molecule, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, "", err
	}

	parts := strings.Split(input, "\n\n")
	return parseReplacements(parts[0]), Molecule(strings.TrimSpace(parts[1])), nil
}

func parseReplacements(input string) []Replacement {
	var replacements []Replacement
	for _, line := range strings.Split(input, "\n") {
		parts := strings.Split(line, " => ")
		replacements = append(replacements, Replacement{parts[0], parts[1]})
	}
	return replacements
}

type Replacement [2]string

func FindShortestPath(start, target Molecule, replacements []Replacement) int {
	type Path struct {
		current Molecule
		length  int
	}
	stack := queue.NewPriorityQueue[Path]()
	startToTarget := levenshtein.ComputeDistance(string(start), string(target))
	stack.Push(Path{current: start, length: 0}, startToTarget)

	seen := map[Molecule]int{}
	for stack.Size() > 0 {
		path := stack.Pop()

		nexts := path.current.PossibleMolecules(replacements)
		if nexts.Contains(target) {
			return path.length + 1
		}
		for molecule := range nexts {
			next := Path{current: molecule, length: path.length + 1}
			// all replacement LHS are shorter or equal to RHS so there's no way to reduce the
			// molecule length by performing more replacements
			if len(next.current) >= len(target) {
				continue
			}
			// if we've already seen the 
			if length, exists := seen[next.current]; exists && length <= next.length {
				continue
			}
			distToTarget := levenshtein.ComputeDistance(string(next.current), string(target))
			// if we're further away from the target than the start then give up on this branch
			// (not sure this is always correct)
			if distToTarget >= startToTarget {
				continue
			}
			seen[next.current] = path.length
			stack.Push(next, distToTarget)
		}
	}
	return -1
}

type Molecule string

func (m Molecule) PossibleMolecules(replacements []Replacement) set.Set[Molecule] {
	possible := set.NewSet[Molecule]()
	seenIndexes := map[string][]int{}
	for _, r := range replacements {
		old, new := r[0], r[1]
		if _, exists := seenIndexes[old]; !exists {
			seenIndexes[old] = FindAllIndexes(string(m), old)
		}
		for _, index := range seenIndexes[old] {
			possible.Add(m[:index] + Molecule(new) + m[index+len(old):])
		}
	}
	return possible
}

func FindAllIndexes(s, substr string) []int {
	indexes := make([]int, 0, len(s))
	// var indexes []int
	var offset int
	for len(s) > 0 {
		index := strings.Index(s, substr)
		if index < 0 {
			break
		}
		indexes = append(indexes, index+offset)
		s = s[index+len(substr):]
		offset += index + len(substr)
	}
	return indexes
}
