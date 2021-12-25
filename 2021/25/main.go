package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Sea Cucumber", part1, part2).Run()
}

func part1() error {
	ocean, err := readOcean(InputFile)
	if err != nil {
		return err
	}

	moved := 1
	step := 0
	for moved > 0 {
		moved = ocean.MoveHerds()
		step += 1
	}

	fmt.Printf("sea cucumbers stop moving after %d steps\n", step)

	return nil
}

func part2() error {
	fmt.Println("no part 2 today, it's Christmas! 🎄")
	return nil
}

func readOcean(path string) (Ocean, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	ocean := Ocean{}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		row := []rune{}
		for _, r := range line {
			row = append(row, r)
		}
		ocean = append(ocean, row)
	}

	return ocean, nil
}

type Ocean [][]rune

type Coord struct {
	Row, Col int
}

type move struct {
	from, to Coord
}

func (o Ocean) MoveHerds() int {
	moved := 0
	moves := []move{}
	// east
	for i, row := range o {
		for j, r := range row {
			if r == rune('>') {
				east := Coord{i, j + 1}
				if o.CanMove(east) {
					moves = append(moves, move{
						from: Coord{i, j},
						to:   east,
					})
				}
			}
		}
	}
	for len(moves) > 0 {
		move := moves[0]
		moves = moves[1:]
		o.Move(move.from, move.to)
		moved += 1
	}
	//south
	for i, row := range o {
		for j, r := range row {
			if r == rune('v') {
				south := Coord{i + 1, j}
				if o.CanMove(south) {
					moves = append(moves, move{
						from: Coord{i, j},
						to:   south,
					})
				}
			}
		}
	}
	for len(moves) > 0 {
		move := moves[0]
		moves = moves[1:]
		o.Move(move.from, move.to)
		moved += 1
	}

	return moved
}

func (o Ocean) Move(from Coord, to Coord) {
	if !o.CanMove(to) {
		return
	}
	rows := len(o)
	cols := len(o[0])

	o[to.Row % rows][to.Col % cols] = o[from.Row % rows][from.Col % cols]
	o[from.Row % rows][from.Col % cols] = rune('.')
}

func (o Ocean) CanMove(to Coord) bool {
	rows := len(o)
	cols := len(o[0])
	return o[to.Row % rows][to.Col % cols] == rune('.')
}

func (o Ocean) String() string {
	lines := []string{}
	for _, row := range o {
		line := ""
		for _, r := range row {
			line += string(r)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
