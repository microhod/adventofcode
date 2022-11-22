package main

import (
	"bytes"
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Aunt Sue", part1, part2).Run()
}

func part1() error {
	sues, err := parse(InputFile)
	if err != nil {
		return err
	}

	for index, sue := range sues {
		if match(sue, target, equal) {
			fmt.Println("The real Aunt Sue is", index+1)
		}
	}
	return nil
}

func part2() error {
	sues, err := parse(InputFile)
	if err != nil {
		return err
	}

	for index, sue := range sues {
		if match(sue, target, rangeOrEqual) {
			fmt.Println("The real Aunt Sue is", index+1)
		}
	}
	return nil
}

var target = Sue{
	"children":    3,
	"cats":        7,
	"samoyeds":    2,
	"pomeranians": 3,
	"akitas":      0,
	"vizslas":     0,
	"goldfish":    5,
	"trees":       3,
	"cars":        2,
	"perfumes":    1,
}

func parse(path string) ([]Sue, error) {
	data, err := file.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	// convert file to yaml
	d := regexp.MustCompile(`Sue \d+: `).ReplaceAll(data, []byte("- "))
	d = bytes.ReplaceAll(d, []byte(", "), []byte("\n  "))

	var sues []Sue
	return sues, yaml.Unmarshal(d, &sues)
}

type Sue map[string]int

func match(sue, target Sue, f func(key string, sv, tv int) bool) bool {
	for tk, tv := range target {
		if sv, ok := sue[tk]; ok && !f(tk, sv, tv) {
			return false
		}
	}
	return true
}

func equal(_ string, sv, tv int) bool {
	return sv == tv
}

func rangeOrEqual(key string, sv, tv int) bool {
	switch key {
	case "cats":
		fallthrough
	case "trees":
		return sv > tv
	case "pomeranians":
		fallthrough
	case "goldfish":
		return sv < tv
	default:
		return sv == tv
	}
}
