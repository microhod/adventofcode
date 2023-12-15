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
	puzzle.NewSolution("Lens Library", part1, part2).Run()
}

func part1() error {
	sequence, _, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	var total int
	for _, step := range sequence {
		total += Hash(step)
	}

	fmt.Printf("the sum of all sequence hashes is: %d\n", total)
	return nil
}

func part2() error {
	_, operations, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	boxes := make(Boxes, 256)
	for _, op := range operations {
		op.Run(boxes)
	}

	fmt.Printf("focusing power: %d\n", boxes.FocusingPower())
	return nil
}

func parseInput(path string) ([]string, []BoxOperation, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, nil, err
	}

	input = strings.TrimSpace(input)
	parts := strings.Split(input, ",")

	var operations []BoxOperation
	for _, part := range parts {
		op, err := parseOperation(part)
		if err != nil {
			return nil, nil, err
		}
		operations = append(operations, op)
	}

	return strings.Split(input, ","), operations, nil
}

func parseOperation(input string) (BoxOperation, error) {
	if strings.Contains(input, "-") {
		return Remove{Label: input[:len(input)-1]}, nil
	}

	parts := strings.Split(input, "=")
	focalLength, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	return Add{Label: parts[0], FocalLength: focalLength}, nil
}

func Hash(input string) int {
	var hash int
	for _, ch := range input {
		hash += int(ch)
		hash *= 17
		hash %= 256
	}
	return hash
}

type Boxes []Lenses

func (b Boxes) FocusingPower() int {
	var power int
	for box := range b {
		for slot, lens := range b[box] {
			power += (box+1) * (slot+1) * lens.FocalLength
		}
	}
	return power
}

type BoxOperation interface {
	Run(Boxes)
}

type Remove struct {
	Label string
}

func (r Remove) Run(boxes Boxes) {
	box := Hash(r.Label)
	index := boxes[box].IndexOfLabel(r.Label)
	if index < 0 {
		return
	}

	boxes[box] = append(boxes[box][:index], boxes[box][index+1:]...)
}

type Add Lens

func (a Add) Run(boxes Boxes) {
	box := Hash(a.Label)
	index := boxes[box].IndexOfLabel(a.Label)

	if index < 0 {
		boxes[box] = append(boxes[box], Lens(a))
		return
	}
	boxes[box][index] = Lens(a)
}

type Lenses []Lens

func (l Lenses) IndexOfLabel(label string) int {
	for i, lens := range l {
		if lens.Label == label {
			return i
		}
	}
	return -1
}

type Lens struct {
	Label       string
	FocalLength int
}
