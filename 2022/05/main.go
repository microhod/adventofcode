package main

import (
	"bytes"
	"fmt"
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
	puzzle.NewSolution("Supply Stacks", part1, part2).Run()
}

func part1() error {
	stacks, moves, err := parse(InputFile)
	if err != nil {
		return err
	}

	for _, move := range moves {
		CrateMover3000{}.Rearrange(stacks, move)
	}

	var message string
	for _, stack := range stacks {
		message += stack[0]
	}

	fmt.Printf("top of stacks after rearrangement: %s\n", message)
	return nil
}

func part2() error {
	stacks, moves, err := parse(InputFile)
	if err != nil {
		return err
	}

	for _, move := range moves {
		CrateMover3001{}.Rearrange(stacks, move)
	}

	var message string
	for _, stack := range stacks {
		message += stack[0]
	}

	fmt.Printf("top of stacks after rearrangement: %s\n", message)
	return nil
}

func parse(path string) (Stacks, []Move, error) {
	data, err := file.ReadBytes(path)
	if err != nil {
		return nil, nil, err
	}

	parts := bytes.Split(data, []byte("\n\n"))

	stacks, err := parseStacks(strings.Split(string(parts[0]), "\n"))
	if err != nil {
		return nil, nil, err
	}
	moves, err := parseMoves(strings.Split(string(parts[1]), "\n"))
	if err != nil {
		return nil, nil, err
	}

	return stacks, moves, nil
}

func parseStacks(lines []string) (Stacks, error) {
	// find max stack number
	line := strings.ReplaceAll(lines[len(lines)-1], "   ", ",")
	numbers, err := csv.ParseInts(strings.TrimSpace(line))
	if err != nil {
		return nil, err
	}
	numStacks := numbers[len(numbers)-1]

	// remove bottom line of "stack numbers"
	lines = lines[:len(lines)-1]

	stacks := make(Stacks, numStacks)
	for _, line := range lines {
		if line == "" {
			continue
		}

		for i := 0; 1+(4*i) < len(line); i++ {
			crate := string(line[1+(4*i)])
			if strings.TrimSpace(crate) == "" {
				continue
			}

			stacks[i] = append(stacks[i], crate)
		}
	}

	return stacks, nil
}

func parseMoves(lines []string) ([]Move, error) {
	var moves []Move

	for _, line := range lines {
		if line == "" {
			continue
		}
		line = strings.ReplaceAll(line, "move ", "")
		line = strings.ReplaceAll(line, " from ", ",")
		line = strings.ReplaceAll(line, " to ", ",")

		numbers, err := csv.ParseInts(line)
		if err != nil {
			return nil, err
		}

		moves = append(moves, Move{
			Quantity: numbers[0],
			// convert to 0 indexing
			From: numbers[1]-1,
			To: numbers[2]-1,
		})
	}
	return moves, nil
}

type Stacks [][]string

type CrateMover3000 struct {}

func (mover CrateMover3000) Rearrange(stacks Stacks, move Move) {
	crates := stacks[move.From][:move.Quantity]
	// remove off From
	stacks[move.From] = stacks[move.From][move.Quantity:]
	// put on To
	for _, crate := range crates {
		stacks[move.To] = append([]string{crate}, stacks[move.To]...)
	}
}

type CrateMover3001 struct {}

func (mover CrateMover3001) Rearrange(stacks Stacks, move Move) {
	crates := stacks[move.From][:move.Quantity]
	// remove off From
	stacks[move.From] = stacks[move.From][move.Quantity:]
	// put on To, retaining their order
	for i := len(crates)-1; i >= 0; i-- {
		stacks[move.To] = append([]string{crates[i]}, stacks[move.To]...)
	}
}

type Move struct {
	Quantity, From, To int
}
