package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/deckarep/golang-set"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	EntriesFile = "entries.txt"
	TestFile    = "test.txt"

	Digit0 = "abcefg"
	Digit1 = "cf"
	Digit2 = "acdeg"
	Digit3 = "acdfg"
	Digit4 = "bcdf"
	Digit5 = "abdfg"
	Digit6 = "abdefg"
	Digit7 = "acf"
	Digit8 = "abcdefg"
	Digit9 = "abcdfg"
)

var digitMapping = map[string]int{
	Digit0: 0,
	Digit1: 1,
	Digit2: 2,
	Digit3: 3,
	Digit4: 4,
	Digit5: 5,
	Digit6: 6,
	Digit7: 7,
	Digit8: 8,
	Digit9: 9,
}

func main() {
	puzzle.NewSolution("Seven Segment Search", part1, part2).Run()
}

func part1() error {
	entries, err := readEntries(EntriesFile)
	if err != nil {
		return err
	}

	count := 0
	for _, entry := range entries {
		count += len(filterByLength(entry.Output, len(Digit1)))
		count += len(filterByLength(entry.Output, len(Digit4)))
		count += len(filterByLength(entry.Output, len(Digit7)))
		count += len(filterByLength(entry.Output, len(Digit8)))
	}

	fmt.Printf("# 1, 4, 7 and 8 digits: %d\n", count)

	return nil
}

func part2() error {
	entries, err := readEntries(EntriesFile)
	if err != nil {
		return err
	}

	sum := 0
	for _, entry := range entries {
		value, err := getOutputValue(entry)
		if err != nil {
			return err
		}
		sum += value
	}
	fmt.Printf("total sum of output values: %d\n", sum)

	return nil
}

type Entry struct {
	Patterns []string
	Output   []string
}

func findLetterMapping(entry Entry) (map[string]string, error) {
	options := map[string]mapset.Set{}

	digitErr := func(length, expected, actual int) error {
		return fmt.Errorf("not enough information, expected %d digit/s of length %d but found %d", expected, length, actual)
	}

	// find digits based on the length of their strings
	digit1s := filterByLength(entry.Patterns, len(Digit1))
	if len(digit1s) < 1 {
		return nil, digitErr(len(Digit1), 1, len(digit1s))
	}
	digit4s := filterByLength(entry.Patterns, len(Digit4))
	if len(digit4s) < 1 {
		return nil, digitErr(len(Digit4), 1, len(digit4s))
	}
	digit7s := filterByLength(entry.Patterns, len(Digit7))
	if len(digit7s) < 1 {
		return nil, digitErr(len(Digit7), 1, len(digit7s))
	}
	digit8s := filterByLength(entry.Patterns, len(Digit8))
	if len(digit8s) < 1 {
		return nil, digitErr(len(Digit8), 1, len(digit8s))
	}
	digit069s := filterByLength(entry.Patterns, len(Digit0))
	if len(digit069s) < 3 {
		return nil, digitErr(len(Digit0), 3, len(digit069s))
	}
	digit069s = distinct(digit069s)
	if len(digit069s) < 3 {
		return nil, fmt.Errorf("not enough information, expected at least 3 distinct digits of length %d but got %d", len(Digit0), len(digit069s))
	}

	// sets of letters
	one := letterSet(digit1s[0])
	four := letterSet(digit4s[0])
	seven := letterSet(digit7s[0])
	eight := letterSet(digit8s[0])
	// these might not be in the correct order as we just go off length and 0, 6 and 9 have the same length
	// this is okay as the operations used don't depend on ordering, it's just easier to name them based on their numbers
	zero := letterSet(digit069s[0])
	six := letterSet(digit069s[1])
	nine := letterSet(digit069s[2])

	// 7 \ 1 => A
	options["a"] = seven.Difference(one)
	// 4 \ 1 => {B, D}
	options["b"] = four.Difference(one)
	options["d"] = four.Difference(one)
	// 1 => {C, F}
	options["c"] = one
	options["f"] = one
	// 8 \ 4 \ { A } => {E, G}
	options["e"] = eight.Difference(four).Difference(options["a"])
	options["g"] = eight.Difference(four).Difference(options["a"])
	// (0 union 6 union 9) \ (0 intersect 6 intersect 9) => {C, D, E}
	cde := (zero.Union(six).Union(nine)).Difference(zero.Intersect(six).Intersect(nine))
	// { B, D } intersect { C, D, E } => D
	options["d"] = options["d"].Intersect(cde)
	// { B, D } \ { D } => B
	options["b"] = options["b"].Difference(options["d"])
	// { C, F } intersect { C, D, E } => C
	options["c"] = options["c"].Intersect(cde)
	// { C, F } \ { C } => F
	options["f"] = options["f"].Difference(options["c"])
	// { E, G } intersect { C, D, E } => E
	options["e"] = options["e"].Intersect(cde)
	// { E, G } \ { E } => G
	options["g"] = options["g"].Difference(options["e"])

	mappings := map[string]string{}
	for letter, set := range options {
		// by this point, we should have only one option for all letters
		if set.Cardinality() != 1 {
			return nil, fmt.Errorf("expected one option for letter %s, but got %d options", letter, set.Cardinality())
		}
		mapping, _ := set.Pop().(string)
		// swap so that the map goes from from entry letter to standard letter
		mappings[mapping] = letter
	}

	return mappings, nil
}

func getOutputValue(entry Entry) (int, error) {
	mapping, err := findLetterMapping(entry)
	if err != nil {
		return 0, err
	}

	digits := []int{}
	for _, digit := range entry.Output {
		mappedDigit := ""
		// fix digit to have correct letters using mapping
		for _, letter := range strings.Split(digit, "") {
			mappedDigit += mapping[letter]
		}
		// alphabetise so that it matches one of the digit constants
		mappedDigit = alphabetise(mappedDigit)
		digits = append(digits, digitMapping[mappedDigit])
	}

	// construct the value using the digits
	value := 0
	for i := range digits {
		value += digits[len(digits)-1-i] * int(math.Pow10(i))
	}

	return value, nil
}

func readEntries(path string) ([]Entry, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	entries := []Entry{}
	for _, line := range lines {
		parts := strings.Split(line, " | ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected line to split into 2 parts by the '|' delimeter, but got %d parts", len(parts))
		}

		entries = append(entries, Entry{
			Patterns: strings.Fields(parts[0]),
			Output:   strings.Fields(parts[1]),
		})
	}

	return entries, nil
}

func letterSet(str string) mapset.Set {
	slice := []interface{}{}
	for _, letter := range strings.Split(str, "") {
		slice = append(slice, letter)
	}
	return mapset.NewSetFromSlice(slice)
}

func filterByLength(list []string, length int) []string {
	filtered := []string{}
	for _, str := range list {
		if len(str) == length {
			filtered = append(filtered, str)
		}
	}
	return filtered
}

func contains(strs []string, str string) bool {
	for _, s := range strs {
		if str == s {
			return true
		}
	}
	return false
}

func distinct(strs []string) []string {
	filtered := []string{}
	for _, str := range strs {
		if !contains(filtered, str) {
			filtered = append(filtered, str)
		}
	}
	return filtered
}

func alphabetise(str string) string {
	s := strings.Split(str, "")
	sort.Strings(s)
	return strings.Join(s, "")
}
