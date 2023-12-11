package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/geometry/plane"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Cosmic Expansion", part1, part2).Run()
}

func part1() error {
	universe, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	universe.Expand(2)
	total := maths.Sum(universe.ShortestPaths()...)

	fmt.Printf("sum of all shortest paths between galaxies: %d\n", total)
	return nil
}

func part2() error {
	universe, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	universe.Expand(1000000)
	total := maths.Sum(universe.ShortestPaths()...)

	fmt.Printf("sum of all shortest paths between galaxies: %d\n", total)
	return nil
}

func parseInput(path string) (Universe, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return Universe{}, err
	}

	var universe Universe
	for y := range lines {
		for x := range lines[y] {
			if lines[y][x] == '#' {
				galaxy := plane.Vector{X: x, Y: y}
				universe.Galaxies = append(universe.Galaxies, galaxy)
			}
		}
	}
	return universe, nil
}

type Universe struct {
	Galaxies []plane.Vector
}

func (u Universe) ShortestPaths() []int {
	var paths []int
	for i := range u.Galaxies {
		for j := i+1; j < len(u.Galaxies); j++ {
			paths = append(paths, plane.ManhattanMetric(u.Galaxies[i], u.Galaxies[j]))
		}
	}
	return paths
}

func (u Universe) Expand(factor int) {
	var max plane.Vector
	galaxyRows := set.NewSet[int]()
	galaxyCols := set.NewSet[int]()
	for _, g := range u.Galaxies {
		max.X = maths.Max(max.X, g.X)
		max.Y = maths.Max(max.Y, g.Y)

		galaxyRows.Add(g.Y)
		galaxyCols.Add(g.X)
	}

	var expansionRows []int
	for y := 0; y < max.Y; y++ {
		if galaxyRows.Contains(y) {
			continue
		}
		expansionRows = append(expansionRows, y)
	}
	var expansionCols []int
	for x := 0; x < max.X; x++ {
		if galaxyCols.Contains(x) {
			continue
		}
		expansionCols = append(expansionCols, x)
	}

	for i := range u.Galaxies {
		dy := 0
		for _, row := range expansionRows {
			if u.Galaxies[i].Y > row {
				dy += factor-1
			}
		}
		u.Galaxies[i].Y += dy
		dx := 0
		for _, col := range expansionCols {
			if u.Galaxies[i].X > col {
				dx += factor-1
			}
		}
		u.Galaxies[i].X += dx
	}
}
