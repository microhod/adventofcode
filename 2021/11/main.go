package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Dumbo Octopus", part1, part2).Run()
}

func part1() error {
	grid, err := readOctopuses(InputFile)
	if err != nil {
		return err
	}

	steps := 100

	flashes := 0
	for i := 0; i < 100; i++ {
		flashes += grid.Step()
	}

	fmt.Printf("flashes after %d steps: %d\n", steps, flashes)

	return nil
}

func part2() error {
	grid, err := readOctopuses(InputFile)
	if err != nil {
		return err
	}

	step := 1
	for {
		if grid.Step() == 100 {
			fmt.Printf("flashes synchronise on step %d\n", step)
			return nil
		}
		step += 1
	}
}

func readOctopuses(path string) (Grid, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return Grid{}, err
	}

	octopuses := Grid{}
	for i, line := range lines {
		nums, err := parseInts(line)
		if err != nil {
			return Grid{}, err
		}

		row := []*Octopus{}
		for j, n := range nums {
			row = append(row, &Octopus{Energy: n, Coord: Coordinate{
				Row:    i,
				Column: j,
			}})
		}

		octopuses = append(octopuses, row)
	}

	return octopuses, nil
}

func parseInts(str string) ([]int, error) {
	nums := []int{}
	for _, s := range strings.Split(str, "") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}

		nums = append(nums, n)
	}

	return nums, nil
}

type Grid [][]*Octopus

type Octopus struct {
	Energy  int
	Coord   Coordinate
	Flashed bool
}

func (octopus *Octopus) Step() {
	octopus.Energy = (octopus.Energy + 1) % 10
}

func (octopus *Octopus) Flash() {
	octopus.Flashed = true
}

func (octopus *Octopus) String() string {
	return fmt.Sprint(octopus.Energy)
}

type Coordinate struct {
	Row, Column int
}

func (grid Grid) Step() int {
	flashStack := []*Octopus{}
	
	for _, row := range grid {
		for _, o := range row {
			o.Flashed = false
			o.Step()
			
			if o.Energy == 0 {
				flashStack = append(flashStack, o)
			}
		}
	}
	
	flashed := []*Octopus{}

	for len(flashStack) > 0 {
		octopus := flashStack[0]
		flashStack = flashStack[1:]

		octopus.Flash()
		flashed = append(flashed, octopus)

		for _, o := range grid.Neighbourhood(octopus) {
			o.Step()
			// ignore flashed neighbours
			if !o.Flashed && o.Energy == 0 {
				flashStack = append(flashStack, o)
			}
		}
	}

	for _, o := range flashed {
		o.Energy = 0
	}

	return len(flashed)
}

func (grid Grid) Neighbourhood(octopus *Octopus) []*Octopus {
	coord := octopus.Coord

	top := max(coord.Row-1, 0)
	left := max(coord.Column-1, 0)
	right := min(coord.Column+1, len(grid.Row(coord.Row))-1)
	bottom := min(coord.Row+1, len(grid)-1)

	neighbours := []*Octopus{}
	for row := top; row <= bottom; row++ {
		for column := left; column <= right; column++ {
			// skip octopus
			if row == coord.Row && column == coord.Column {
				continue
			}

			neighbours = append(neighbours, grid[row][column])
		}
	}

	return neighbours
}

func (grid Grid) Row(index int) []*Octopus {
	return grid[index]
}

func (grid Grid) Column(index int) []*Octopus {
	column := []*Octopus{}
	for _, row := range grid {
		column = append(column, row[index])
	}
	return column
}

func (grid Grid) String() string {
	lines := []string{}
	for _, row := range grid {
		line := ""
		for _, o := range row {
			line += o.String()
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
