package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Perfectly Spherical Houses in a Vacuum", part1, part2).Run()
}

func part1() error {
	vectors, err := readVectors(InputFile)
	if err != nil {
		return err
	}

	pos := Vec(0,0)
	locations := map[Vector]int{
		pos: 1,
	}
	for _, vector := range vectors {
		pos = pos.Add(vector)
		locations[pos] += 1
	}

	fmt.Printf("number of locations visited: %d\n", len(locations))
	return nil
}

func part2() error {
	vectors, err := readVectors(InputFile)
	if err != nil {
		return err
	}

	pos := Vec(0,0)
	locations := map[Vector]int{
		pos: 1,
	}
	santa := pos
	robosanta := pos
	for i, vector := range vectors {
		if i %2 == 0 {
			santa = santa.Add(vector)
			locations[santa] += 1
			continue
		}

		robosanta = robosanta.Add(vector)
		locations[robosanta] += 1
	}

	fmt.Printf("number of locations visited: %d\n", len(locations))
	return nil
}

type Vector struct {
	X, Y int
}

func Vec(x, y int) Vector {
	return Vector{X: x, Y: y}
}

func (v Vector) Add(w Vector) Vector {
	return Vec(v.X + w.X, v.Y + w.Y)
}

func readVectors(path string) ([]Vector, error) {
	bytes, err := file.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	var vectors []Vector
	for _, b := range bytes {
		switch rune(b){
		case '^':
			vectors = append(vectors, Vec(0, 1))
		case '<':
			vectors = append(vectors, Vec(-1, 0))
		case '>':
			vectors = append(vectors, Vec(1, 0))
		case 'v':
			vectors = append(vectors, Vec(0, -1))
		default:
			return nil, fmt.Errorf("invalid character: '%s'", string(b))
		}
	}
	return vectors, nil
}
