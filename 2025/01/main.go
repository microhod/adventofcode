package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Secret Entrance", part1, part2).Run()
}

func part1() error {
	rotations, err := parse(InputFile)
	if err != nil {
		return err
	}

	dial := 50
	zeros := 0
	for _, r := range rotations {
		dial = maths.Mod(dial+r, 100)
		if dial == 0 {
			zeros++
		}
	}
	fmt.Println("password:", zeros)
	return nil
}

func part2() error {
	rotations, err := parse(InputFile)
	if err != nil {
		return err
	}

	dial := 50
	zeros := 0
	for _, r := range rotations {
		if r > 0 {
			zeros += (dial + r) / 100
		} else {
			if maths.Abs(r) >= dial && dial > 0 {
				zeros += 1 // reversing past zero the 1st time
			}
			zeros += (maths.Abs(r) - dial) / 100 // full reverse rotations
	 	}
		dial = maths.Mod(dial+r, 100)
	}
	fmt.Println("password:", zeros)

	return nil
}

func parse(path string) ([]int, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}
	var rotations []int
	for i, line := range lines {
		line = strings.ReplaceAll(line, "R", "")
		line = strings.ReplaceAll(line, "L", "-")
		rotation, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid rotation on line %d: %s: %w", i, line, err)
		}
		rotations = append(rotations, rotation)
	}
	return rotations, nil
}
