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
	puzzle.NewSolution("I Was Told There Would Be No Math", part1, part2).Run()
}

func part1() error {
	boxes, err := readBoxes(InputFile)
	if err != nil {
		return err
	}

	var requiredPaper int
	for _, box := range boxes {
		requiredPaper += box.GetRequiredPaper()
	}

	fmt.Printf("total required paper: %d\n", requiredPaper)
	return nil
}

func part2() error {
	boxes, err := readBoxes(InputFile)
	if err != nil {
		return err
	}

	var requiredLength int
	for _, box := range boxes {
		requiredLength += box.GetRequiredBowLength()
	}

	fmt.Printf("total required bow length: %d\n", requiredLength)
	return nil
}

type Box struct {
	l, w, h int
}

func (b Box) GetRequiredPaper() int {
	areas := []int{b.l * b.w, b.l * b.h, b.w * b.h}
	sort.Ints(areas)
	extra := areas[0]

	// 2*l*w + 2*w*h + 2*h*l + smallest_side_area
	return sum([]int{2*areas[0], 2*areas[1], 2*areas[2], extra})
}

func (b Box) GetRequiredBowLength() int {
	areas := []int{b.l, b.w, b.h}
	sort.Ints(areas)
	ribbon := 2*(areas[0] + areas[1])
	bow := b.l * b.w * b.h


	return ribbon + bow
}

func readBoxes(path string) ([]Box, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var boxes []Box
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "x")
		if len(parts) < 3 {
			return nil, fmt.Errorf("need 3 dimensions but got %d on line %d", len(parts), idx)
		}

		var box Box
		box.l, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("dimension not int on line %d: %w", idx, err)
		}
		box.w, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("dimension not int on line %d: %w", idx, err)
		}
		box.h, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("dimension not int on line %d: %w", idx, err)
		}

		boxes = append(boxes, box)
	}

	return boxes, nil
}

func sum(nums []int) int {
	var sum int
	for _, n := range nums {
		sum += n
	}

	return sum
}
