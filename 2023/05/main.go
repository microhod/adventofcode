package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("If You Give A Seed A Fertilizer", part1, part2).Run()
}

func part1() error {
	seeds, almanac, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	min := math.MaxInt
	for _, seed := range seeds {
		min = maths.Min(min, almanac.SeedToLocation(seed))
	}
	fmt.Printf("the closest seed location is: %d\n", min)
	return nil
}

func part2() error {
	seeds, almanac, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	var seedRanges []maths.Range
	for i := 0; i < len(seeds)-1; i += 2 {
		seedRanges = append(seedRanges, maths.Range{
			Left: seeds[i],
			Right: seeds[i] + seeds[i+1] - 1},
		)
	}

	destinations := almanac.SeedRangesToLocations(seedRanges)
	min := math.MaxInt
	for _, dst := range destinations {
		min = maths.Min(min, dst.Left)
	}

	fmt.Printf("the closest seed location is: %d\n", min)
	return nil
}

func parseInput(path string) ([]int, Almanac, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, nil, err
	}

	parts := strings.Split(input, "\n\n")
	seeds, err := csv.ParseInts(parts[0][7:], " ")
	if err != nil {
		return nil, nil, err
	}

	almanac := make(Almanac)
	for _, part := range parts[1:] {
		name, mapping, err := parseMap(part)
		if err != nil {
			return nil, nil, err
		}
		almanac[name] = mapping
	}
	return seeds, almanac, nil
}

func parseMap(input string) (string, AlmanacMap, error) {
	parts := strings.Split(input, " map:\n")
	name := parts[0]

	var mapping AlmanacMap
	for _, line := range strings.Split(parts[1], "\n") {
		if line == "" {
			continue
		}

		nums, err := csv.ParseInts(line, " ")
		if err != nil {
			return "", AlmanacMap{}, err
		}
		mapping.DstRanges = append(mapping.DstRanges, maths.Range{
			Left: nums[0], Right: nums[0] + nums[2],
		})
		mapping.SrcRanges = append(mapping.SrcRanges, maths.Range{
			Left: nums[1], Right: nums[1] + nums[2],
		})
	}
	return name, mapping, nil
}

type Almanac map[string]AlmanacMap

func (a Almanac) SeedToLocation(input int) int {
	input = a["seed-to-soil"].Dst(input)
	input = a["soil-to-fertilizer"].Dst(input)
	input = a["fertilizer-to-water"].Dst(input)
	input = a["water-to-light"].Dst(input)
	input = a["light-to-temperature"].Dst(input)
	input = a["temperature-to-humidity"].Dst(input)
	return a["humidity-to-location"].Dst(input)
}

func (a Almanac) SeedRangesToLocations(input []maths.Range) []maths.Range {
	input = a["seed-to-soil"].DstRange(input)
	input = a["soil-to-fertilizer"].DstRange(input)
	input = a["fertilizer-to-water"].DstRange(input)
	input = a["water-to-light"].DstRange(input)
	input = a["light-to-temperature"].DstRange(input)
	input = a["temperature-to-humidity"].DstRange(input)
	input = a["humidity-to-location"].DstRange(input)
	return input
}

type AlmanacMap struct {
	SrcRanges []maths.Range
	DstRanges []maths.Range
}

func (m AlmanacMap) Dst(src int) int {
	for i, s := range m.SrcRanges {
		if s.ContainsValue(src) {
			offset := src - s.Left
			return m.DstRanges[i].Left + offset
		}
	}
	return src
}

func (m AlmanacMap) DstRange(inputs []maths.Range) []maths.Range {
	var destinations []maths.Range

	for i := range m.SrcRanges {
		src := m.SrcRanges[i]
		dst := m.DstRanges[i]

		inputCount := len(inputs)
		for j := 0; j < inputCount; j++ {
			input := inputs[0]
			inputs = inputs[1:]

			intersection := src.Intersect(input)
			if !intersection.Valid() {
				// if no intersection, another source 
				// should be compared against the input
				inputs = append(inputs, input)
				continue
			}

			// compute the destination range of the intersection
			destinations = append(destinations, maths.Range{
				Left:  intersection.Left + dst.Left - src.Left,
				Right: intersection.Right + dst.Right - src.Right,
			})
			// keep all other input ranges to compare with other sources
			inputs = append(inputs, input.Diff(src)...)
		}
	}
	// all remaining inputs have no destination mapping
	// which means they are mapped to themselves
	return append(destinations, inputs...)
}
