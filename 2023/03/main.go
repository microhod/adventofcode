package main

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/geometry/plane"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Gear Ratios", part1, part2).Run()
}

func part1() error {
	schematic, err := parseSchematic(InputFile)
	if err != nil {
		return err
	}

	var partNumbers int
	for _, number := range schematic.ParseNumbers() {
		if schematic.IsPartNumber(number) {
			value, err := schematic.NumberValue(number)
			if err != nil {
				return err
			}
			partNumbers += value
		}
	}
	fmt.Printf("the sum of all part numbers is: %d\n", partNumbers)
	return nil
}

func part2() error {
	schematic, err := parseSchematic(InputFile)
	if err != nil {
		return err
	}

	gearRatios := make(map[plane.Vector][]int)
	for _, number := range schematic.ParseNumbers() {
		for _, gear := range schematic.GetAdjacentGears(number) {
			value, err := schematic.NumberValue(number)
			if err != nil {
				return err
			}
			gearRatios[gear] = append(gearRatios[gear], value)
		}
	}

	var totalRatios int
	for _, ratios := range gearRatios {
		if len(ratios) != 2 {
			continue
		}
		totalRatios += ratios[0] * ratios[1]
	}

	fmt.Printf("the sum of all gear ratios is: %d\n", totalRatios)
	return nil
}

func parseSchematic(path string) (Schematic, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var schematic Schematic
	for _, line := range lines {
		schematic = append(schematic, []rune(line))
	}
	return schematic, nil
}

type Schematic [][]rune

func (s Schematic) ParseNumbers() []Number {
	var numbers []Number
	for i, row := range s {
		var current *Number
		for j, char := range row {
			if unicode.IsNumber(char) {
				if current == nil {
					current = &Number{Start: plane.Vector{Y: i, X: j}}
				}
				current.Length += 1
				continue
			}
			if current != nil {
				numbers = append(numbers, *current)
			}
			current = nil
		}
		if current != nil {
			numbers = append(numbers, *current)
		}
		current = nil
	}
	return numbers
}

func (s Schematic) IsPartNumber(number Number) bool {
	for _, n := range s.numberNeighbourhood(number) {
		if !unicode.IsNumber(s[n.Y][n.X]) && s[n.Y][n.X] != '.' {
			return true
		}
	}
	return false
}

func (s Schematic) GetAdjacentGears(number Number) []plane.Vector {
	var gears []plane.Vector
	for _, n := range s.numberNeighbourhood(number) {
		if s[n.Y][n.X] == '*' {
			gears = append(gears, n)
		}
	}
	return gears
}

func (s Schematic) numberNeighbourhood(number Number) []plane.Vector {
	var neighbours []plane.Vector
	// left
	if number.Start.X > 0 {
		neighbours = append(neighbours, plane.Vector{
			X: number.Start.X-1,
			Y: number.Start.Y,
		})
	}
	// right
	if number.Start.X + number.Length < len(s[0]) {
		neighbours = append(neighbours, plane.Vector{
			X: number.Start.X + number.Length,
			Y: number.Start.Y,
		})
	}
	// top
	if number.Start.Y > 0 {
		start := maths.Max(number.Start.X-1, 0)
		end := maths.Min(number.Start.X + number.Length + 1, len(s[0]))
		for x :=start; x < end; x++ {
			neighbours = append(neighbours, plane.Vector{
				X: x,
				Y: number.Start.Y-1,
			})
		}
	}
	// bottom
	if number.Start.Y < len(s[0])-1 {
		start := maths.Max(number.Start.X-1, 0)
		end := maths.Min(number.Start.X + number.Length + 1, len(s[0]))
		for x :=start; x < end; x++ {
			neighbours = append(neighbours, plane.Vector{
				X: x,
				Y: number.Start.Y+1,
			})
		}
	}
	return neighbours
}

func (s Schematic) NumberValue(number Number) (int, error) {
	return strconv.Atoi(string(s[number.Start.Y][number.Start.X:number.Start.X+number.Length]))
}

type Number struct {
	Start  plane.Vector
	Length int
}
