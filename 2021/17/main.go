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
	puzzle.NewSolution("Trick Shot", part1, part2).Run()
}

func part1() error {
	box, err := readBox(InputFile)
	if err != nil {
		return err
	}

	y := box.Highest()
	fmt.Printf("highest possible y: %d\n", y)

	return nil
}

func part2() error {
	box, err := readBox(InputFile)
	if err != nil {
		return err
	}

	count := box.NumValid()
	fmt.Printf("num valid: %d\n", count)

	return nil
}

func readBox(path string) (*Box, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	if len(lines) < 1 {
		return nil, nil
	}

	box := &Box{}

	overallParts := strings.Split(lines[0], ", ")
	
	parts := strings.Split(overallParts[0], "=")
	parts = strings.Split(parts[1], "..")
	box.Xmin, err = strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	box.Xmax, err = strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	parts = strings.Split(overallParts[1], "=")
	parts = strings.Split(parts[1], "..")
	box.Ymin, err = strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	box.Ymax, err = strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return box, nil
}

type Box struct {
	Xmin, Xmax, Ymin, Ymax int
}

func (box Box) Highest() int {
	max := 0
	for x := 0; x < box.Xmax; x++ {
		for y := 1000; y > box.Ymin; y-- {
			coords, valid := box.Trajectory(x, y)
			if valid && maxY(coords) > max {
				max = maxY(coords)
			}
		}
	}
	return max
}

func (box Box) NumValid() int {
	count := 0
	for x := 0; x < 10*box.Xmax; x++ {
		for y := 1000; y > 10*box.Ymin; y-- {
			_, valid := box.Trajectory(x, y)
			if valid {
				count += 1
			}
		}
	}
	return count
}

func (box Box) Trajectory(vx, vy int) ([]Coord, bool) {
	prev := Coord{0,0}
	trajectory := []Coord{prev}
	for !box.IsPast(prev) {
		next := Coord{
			prev.X + vx,
			prev.Y + vy,
		}
		trajectory = append(trajectory, next)
		
		if box.IsIn(next) {
			return trajectory, true
		}

		if vx > 0 {
			vx -= 1
		}
		if vx < 0 {
			vx += 1
		}
		vy -= 1

		prev = next
	}

	return trajectory, false
}

func (box Box) IsIn(coord Coord) bool {
	return coord.X >= box.Xmin && coord.X <= box.Xmax && coord.Y >= box.Ymin && coord.Y <= box.Ymax
}

func (box Box) IsPast(coord Coord) bool {
	return coord.Y < box.Ymin
}

type Coord struct {
	X, Y int
}

func maxY(coords []Coord) int {
	max := 0
	for _, c := range coords {
		if c.Y > max {
			max = c.Y
		}
	}
	return max
}
