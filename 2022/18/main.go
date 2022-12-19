package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Boiling Boulders", part1, part2).Run()
}

func part1() error {
	cubes, _, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Println(SurfaceArea(cubes))
	return nil
}

func part2() error {
	shape, grid, err := parse(InputFile)
	if err != nil {
		return err
	}

	complement := grid.ComplementCubes()
	innerCubes := set.NewSet[Cube]()

	// remove outer cubes
	for len(complement) > 0 {
		for cube := range complement {
			// if NOT comppletely contained by complement inner and grid then it is outside the shape 
			// therefore remove it and all neighbours from innerCubes
			neighbours := cube.GetNeighbours()
			
			count := 0
			var complementNeighbours []Cube
			for _, n := range neighbours {
				if shape.Contains(n) {
					count += 1
				}
				if complement.Contains(n) || innerCubes.Contains(n) {
					complementNeighbours = append(complementNeighbours, n)
					count += 1
				}
			}
			if count < len(neighbours) {
				innerCubes.Remove(cube)
				innerCubes.Remove(neighbours...)

				// now we have seen this cube, we can remove it from complement
				// but add its neighbours as they may have already been seen and 
				// incorrectly categoried as inner
				complement.Remove(cube)
				complement.Add(complementNeighbours...)
				continue
			}

			// else add to innerCubes and remove from the complement
			innerCubes.Add(cube)
			complement.Remove(cube)
		}
	}

	fmt.Println(SurfaceArea(shape) - SurfaceArea(innerCubes))
	return nil
}

func parse(path string) (set.Set[Cube], Grid, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, nil, err
	}

	var maxX, maxY, maxZ int

	cubes := set.NewSet[Cube]()
	for _, line := range lines {
		if line == "" {
			continue
		}
		nums, err := csv.ParseInts(line)
		if err != nil {
			return nil, nil, err
		}

		x, y, z := nums[0], nums[1], nums[2]
		cubes.Add(Cube{center: Vector{float64(x), float64(y), float64(z)}})

		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		if z > maxZ {
			maxZ = z
		}
	}

	grid := make(Grid, maxX+1)
	for x := range grid {
		grid[x] = make([][]bool, maxY+1)
		for y := range grid[x] {
			grid[x][y] = make([]bool, maxZ+1)
		}
	}

	for cube := range cubes {
		grid.AddCube(cube)
	}

	return cubes, grid, nil
}

type Cube struct {
	center Vector
}

func (cube Cube) GetSides() []Vector {
	return []Vector{
		{cube.center.X + 0.5, cube.center.Y, cube.center.Z},
		{cube.center.X - 0.5, cube.center.Y, cube.center.Z},
		{cube.center.X, cube.center.Y + 0.5, cube.center.Z},
		{cube.center.X, cube.center.Y - 0.5, cube.center.Z},
		{cube.center.X, cube.center.Y, cube.center.Z + 0.5},
		{cube.center.X, cube.center.Y, cube.center.Z - 0.5},
	}
}
func (cube Cube) GetNeighbours() []Cube {
	return []Cube{
		{center: Vector{cube.center.X + 1, cube.center.Y, cube.center.Z}},
		{center: Vector{cube.center.X - 1, cube.center.Y, cube.center.Z}},
		{center: Vector{cube.center.X, cube.center.Y + 1, cube.center.Z}},
		{center: Vector{cube.center.X, cube.center.Y - 1, cube.center.Z}},
		{center: Vector{cube.center.X, cube.center.Y, cube.center.Z + 1}},
		{center: Vector{cube.center.X, cube.center.Y, cube.center.Z - 1}},
	}
}

func SurfaceArea(cubes set.Set[Cube]) int {
	sides := set.NewSet[Vector]()
	for cube := range cubes {
		for _, side := range cube.GetSides() {
			if sides.Contains(side) {
				sides.Remove(side)
			} else {
				sides.Add(side)
			}
		}
	}
	return len(sides)
}

type Vector struct {
	X, Y, Z float64
}

type Grid [][][]bool

func (grid Grid) AddCube(cube Cube) {
	grid[int(cube.center.X)][int(cube.center.Y)][int(cube.center.Z)] = true
}

func (grid Grid) ComplementCubes() set.Set[Cube] {
	cubes := set.NewSet[Cube]()

	for x := range grid {
		for y := range grid[x] {
			for z := range grid[x][y] {
				if !grid[x][y][z] {
					cubes.Add(Cube{center: Vector{float64(x), float64(y), float64(z)}})
				}
			}
		}
	}

	return cubes
}
