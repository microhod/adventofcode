package main

import (
	"fmt"
	"sort"
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
	puzzle.NewSolution("Smoke Basin", part1, part2).Run()
}

func part1() error {
	hmap, err := readHeightMap(InputFile)
	if err != nil {
		return err
	}

	minima := hmap.GetLocalMinima()
	risk := 0
	for _, min := range minima {
		risk += hmap[min.row][min.column] + 1
	}

	fmt.Printf("risk score: %d\n", risk)

	return nil
}

func part2() error {
	hmap, err := readHeightMap(InputFile)
	if err != nil {
		return err
	}

	basins := hmap.GetBasins()
	sort.Ints(basins)

	if len(basins) < 3 {
		return fmt.Errorf("expected at least 3 basins, but got %d", len(basins))
	}

	b1 := basins[len(basins)-1]
	b2 := basins[len(basins)-2]
	b3 := basins[len(basins)-3]

	fmt.Printf("three largest basins product = %d * %d * %d = %d\n", b1, b2, b3, b1*b2*b3)

	return nil
}

func readHeightMap(path string) (HeightMap, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	heightMap := HeightMap{}
	for _, line := range lines {
		nums := []int{}
		for _, str := range strings.Split(line, "") {
			n, err := strconv.Atoi(str)
			if err != nil {
				return nil, err
			}
			nums = append(nums, n)
		}
		if len(heightMap) > 1 && len(nums) != len(heightMap[len(heightMap)-1]) {
			return nil, fmt.Errorf("all rows must be of equal length")
		}
		heightMap = append(heightMap, nums)
	}

	return heightMap, nil
}

type HeightMap [][]int

type Coordinate struct {
	row, column int
}

func (hmap HeightMap) GetLocalMinima() []Coordinate {
	minima := []Coordinate{}
	for row := range hmap {
		for column := range hmap[row] {
			neighbours := hmap.GetNeighbourHeights(row, column)
			height := hmap[row][column]
			if lessThanAll(height, neighbours) {
				minima = append(minima, Coordinate{
					row:    row,
					column: column,
				})
			}
		}
	}
	return minima
}

func (hmap HeightMap) GetNeighbourHeights(rowIndex, columnIndex int) []int {
	row := hmap.Row(rowIndex)
	column := hmap.Column(columnIndex)

	neighbours := []int{}
	// top
	if rowIndex > 0 {
		neighbours = append(neighbours, column[rowIndex-1])
	}
	// left
	if columnIndex > 0 {
		neighbours = append(neighbours, row[columnIndex-1])
	}
	// right
	if columnIndex < len(row)-1 {
		neighbours = append(neighbours, row[columnIndex+1])
	}
	// top
	if rowIndex < len(column)-1 {
		neighbours = append(neighbours, column[rowIndex+1])
	}

	return neighbours
}

func (hmap HeightMap) GetBasins() []int {
	minima := hmap.GetLocalMinima()
	basins := []int{}
	for _, min := range minima {
		basin := hmap.GetBasin(min)
		basins = append(basins, len(basin))
	}
	return basins
}

// GetBasin uses a flood fill algorithm
// https://en.wikipedia.org/wiki/Flood_fill#Stack-based_recursive_implementation_(four-way)
func (hmap HeightMap) GetBasin(minima Coordinate) []Coordinate {
	basin := []Coordinate{}
	queue := []Coordinate{minima}
	for len(queue) > 0 {
		coord := queue[0]
		queue = queue[1:]

		if hmap[coord.row][coord.column] < 9 {
			if contains(basin, coord) {
				continue
			}
			
			basin = append(basin, coord)
			
			// left
			if coord.column > 0 {
				queue = append(queue, Coordinate{coord.row, coord.column - 1})
			}
			// right
			if coord.column < len(hmap.Row(coord.row))-1 {
				queue = append(queue, Coordinate{coord.row, coord.column + 1})
			}
			// up
			if coord.row > 0 {
				queue = append(queue, Coordinate{coord.row - 1, coord.column})
			}
			// down
			if coord.row < len(hmap)-1 {
				queue = append(queue, Coordinate{coord.row + 1, coord.column})
			}
		}
	}

	return basin
}

func (hmap HeightMap) Row(index int) []int {
	return hmap[index]
}

func (hmap HeightMap) Column(index int) []int {
	column := []int{}
	for _, row := range hmap {
		column = append(column, row[index])
	}
	return column
}

func contains(coords []Coordinate, coord Coordinate) bool {
	for _, c := range coords {
		if coord.row == c.row && coord.column == c.column {
			return true
		}
	}
	return false
}

func lessThanAll(m int, nums []int) bool {
	less := true
	for _, n := range nums {
		if n <= m {
			less = false
		}
	}
	return less
}
