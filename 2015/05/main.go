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
	puzzle.NewSolution("Doesn't He Have Intern-Elves For This?", part1, part2).Run()
}

func part1() error {
	lines, err := file.ReadLines(InputFile)
	if err != nil {
		return err
	}

	var numNice int
	for _, line := range lines {
		if NiceV1(line) {
			numNice += 1
		}
	}

	fmt.Printf("number of nice strings: %d\n", numNice)
	return nil
}

func part2() error {
	lines, err := file.ReadLines(InputFile)
	if err != nil {
		return err
	}

	var numNice int
	for _, line := range lines {
		// fmt.Printf("%s: %v\n", line, NiceV2(line))
		if NiceV2(line) {
			numNice += 1
		}
	}

	fmt.Printf("number of nice strings: %d\n", numNice)
	return nil
}

func NiceV1(input string) bool {
	vowels := map[rune]bool{
		'a': true,
		'e': true,
		'i': true,
		'o': true,
		'u': true,
	}
	var numVowels int

	notAllowedPairs := map[string]bool{
		"ab": true,
		"cd": true,
		"pq": true,
		"xy": true,
	}

	var hasRepeatedLetters bool

	for idx, ch := range input {
		if vowels[ch] {
			numVowels += 1
		}

		if idx < len(input)-1 {
			pair := input[idx : idx+2]
			if notAllowedPairs[pair] {
				return false
			}

			if pair[0] == pair[1] {
				hasRepeatedLetters = true
			}
		}
	}

	return numVowels >= 3 && hasRepeatedLetters
}

func NiceV2(input string) bool {
	pairs := map[string]map[int]bool{}
	for i := 0; i < len(input)-1; i++ {
		pair := input[i:i+2]
		if _, ok := pairs[pair]; !ok {
			pairs[pair] = map[int]bool{}
		}

		pairs[pair][i] = true
		pairs[pair][i+1] = true
	}
	has2Pair := false
	for _, indexes := range pairs {
		// two distinct pairs => 4 distinct indexes
		if len(indexes) >= 4 {
			has2Pair = true
		}
	}

	var triples int
	for i := 0; i < len(input)-2; i++ {
		if input[i] == input[i+2] {
			triples += 1
			break
		}
	}

	return has2Pair && triples >= 1
}
