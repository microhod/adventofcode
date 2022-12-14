package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Regolith Reservoir", part1, part2).Run()
}

func part1() error {
	cave, err := parse(InputFile)
	if err != nil {
		return err
	}

	cave.abyss = len(cave.slice) - 2

	for err == nil {
		err = cave.AddSand()
	}
	if !errors.Is(err, ErrFallenIntoAbyss) {
		return err
	}

	fmt.Println(cave.TotalSand)
	return nil
}

func part2() error {
	cave, err := parse(InputFile)
	if err != nil {
		return err
	}

	for err == nil {
		err = cave.AddSand()
	}
	if !errors.Is(err, ErrSourceBlocked) {
		return err
	}

	fmt.Println(cave.TotalSand)
	return nil
}

func parse(path string) (*Cave, error) {
	input, err := file.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	maxX, err := maxX(input)
	if err != nil {
		return nil, err
	}
	maxY, err := maxY(input)
	if err != nil {
		return nil, err
	}
	// allow space for sand to fall into abyss
	maxX += 1
	maxY += 1

	cave := &Cave{
		sandSouce: [2]int{500, 0},
		abyss:     math.MaxInt,
	}
	cave.slice = make([][]string, maxY+1)
	for i := range cave.slice {
		cave.slice[i] = make([]string, maxX+1)
	}

	for y := range cave.slice {
		for x := range cave.slice[y] {
			cave.slice[y][x] = " "
		}
	}

	for _, line := range strings.Split(string(input), "\n") {
		if line == "" {
			continue
		}
		rocks, err := parseRocks(cave.slice, line)
		if err != nil {
			return nil, err
		}
		for _, rock := range rocks {
			cave.slice[rock[1]][rock[0]] = "#"
		}
	}

	// add floor
	var floor []string
	for range cave.slice[0] {
		floor = append(floor, "#")
	}
	cave.slice = append(cave.slice, floor)

	return cave, nil
}

func maxX(input []byte) (int, error) {
	xCoords := regexp.MustCompile(`\d+,`).FindAll(input, -1)

	var max int
	for _, coord := range xCoords {
		x, err := strconv.Atoi(string(coord[:len(coord)-1]))
		if err != nil {
			return -1, err
		}
		if x > max {
			max = x
		}
	}
	return max, nil
}

func maxY(input []byte) (int, error) {
	yCoords := regexp.MustCompile(`,\d+`).FindAll(input, -1)

	var max int
	for _, coord := range yCoords {
		y, err := strconv.Atoi(string(coord[1:]))
		if err != nil {
			return -1, err
		}
		if y > max {
			max = y
		}
	}
	return max, nil
}

func parseRocks(slice [][]string, line string) ([][2]int, error) {
	var endpoints [][2]int
	for _, coords := range strings.Split(line, " -> ") {
		nums, err := csv.ParseInts(coords)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, [2]int{nums[0], nums[1]})
	}

	var coords [][2]int
	for i := 1; i < len(endpoints); i++ {
		left := endpoints[i-1]
		right := endpoints[i]

		// this is fine because one of the coordinates is always equal
		start := [2]int{min(left[0], right[0]), min(left[1], right[1])}
		end := [2]int{max(left[0], right[0]), max(left[1], right[1])}

		for x := start[0]; x <= end[0]; x++ {
			for y := start[1]; y <= end[1]; y++ {
				coords = append(coords, [2]int{x, y})
			}
		}
	}

	return coords, nil
}

type Cave struct {
	slice     [][]string
	sandSouce [2]int
	abyss     int

	TotalSand int
}

var (
	ErrFallenIntoAbyss = errors.New("fallen into abyss")
	ErrSourceBlocked = errors.New("source is blocked")
)

// AddSand adds a new sand particle and waits for it to come to rest or pass into the abyss,
// it returns true if it falls into the abyss
func (cave *Cave) AddSand() error {
	x, y := cave.sandSouce[0], cave.sandSouce[1]
	for y < cave.abyss {
		if x+1 >= len(cave.slice[0]) {
			cave.extendRight()
		}
		// down
		if cave.slice[y+1][x] == " " {
			y += 1
			continue
		}
		// down + left
		if cave.slice[y+1][x-1] == " " {
			x -= 1
			y += 1
			continue
		}
		// down + right
		if cave.slice[y+1][x+1] == " " {
			x += 1
			y += 1
			continue
		}
		// come to rest
		cave.slice[y][x] = "o"
		cave.TotalSand += 1

		// break if we've clogged the source
		if x == cave.sandSouce[0] && y == cave.sandSouce[1] {
			return ErrSourceBlocked
		}
		return nil
	}
	return ErrFallenIntoAbyss
}

func (cave *Cave) extendRight() {
	for i := range cave.slice {
		value := " "
		// if on the floor
		if i == len(cave.slice)-1 {
			value = "#"
		}
		cave.slice[i] = append(cave.slice[i], value)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
