package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Rope Bridge", part1, part2).Run()
}

func part1() error {
	moves, err := parse(InputFile)
	if err != nil {
		return err
	}

	rope := &Rope{
		Head:  Vector{},
		Tails: make([]Vector, 1),
	}
	visited := set.NewSet[Vector]()

	for _, move := range moves {
		rope.MoveHead(move)
		visited.Add(rope.Tails[0])
	}

	fmt.Println(len(visited))
	return nil
}

func part2() error {
	moves, err := parse(InputFile)
	if err != nil {
		return err
	}

	rope := &Rope{
		Head:  Vector{},
		Tails: make([]Vector, 9),
	}
	visited := set.NewSet[Vector]()

	for _, move := range moves {
		rope.MoveHead(move)
		visited.Add(rope.Tails[8])
	}

	fmt.Println(len(visited))
	return nil
}

func parse(path string) ([]Vector, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	vectors := map[string]Vector{
		"L": {-1, 0},
		"R": {1, 0},
		"U": {0, 1},
		"D": {0, -1},
	}

	var moves []Vector
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		direction, countStr := parts[0], parts[1]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return nil, err
		}

		for i := 0; i < count; i++ {
			moves = append(moves, vectors[direction])
		}
	}

	return moves, nil
}

type Rope struct {
	Head  Vector
	Tails []Vector
}

func (r *Rope) MoveHead(move Vector) {
	r.Head = r.Head.Add(move)

	prev := r.Head
	for i := range r.Tails {
		diff := prev.Minus(r.Tails[i])

		var tailMove Vector
		if abs(diff.X) > 1 || abs(diff.Y) > 1 {
			tailMove.X = diff.X
			tailMove.Y = diff.Y

			// ensure we only move by a maximum of 1 in either direction
			if abs(tailMove.X) > 1 {
				tailMove.X /= abs(tailMove.X)
			}
			if abs(tailMove.Y) > 1 {
				tailMove.Y /= abs(tailMove.Y)
			}
		}

		r.Tails[i] = r.Tails[i].Add(tailMove)
		prev = r.Tails[i]
	}
}

type Vector struct {
	X, Y int
}

func (v Vector) Add(w Vector) Vector {
	return Vector{v.X + w.X, v.Y + w.Y}
}

func (v Vector) Minus(w Vector) Vector {
	return Vector{v.X - w.X, v.Y - w.Y}
}

func abs(n int) int {
	if n < 0 {
		return -1 * n
	}
	return n
}
