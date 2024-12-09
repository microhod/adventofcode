package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/geometry/plane"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Resonant Collinearity", part1, part2).Run()
}

func part1() error {
	antenas, limit, err := parse(InputFile)
	if err != nil {
		return err
	}

	antinodes := antenas.Antinodes(limit, true)

	fmt.Printf("# antinodes: %d\n", len(set.NewSet(antinodes...)))
	return nil
}

func part2() error {
	antenas, limit, err := parse(InputFile)
	if err != nil {
		return err
	}

	antinodes := antenas.Antinodes(limit, false)

	fmt.Printf("# antinodes: %d\n", len(set.NewSet(antinodes...)))
	return nil
}

func parse(path string) (Antenas, plane.Vector, error) {
	return file.ReadVectorsFunc(path, func(b byte) bool {
		return b != '.'
	})
}

type Antenas map[byte][]plane.Vector

func (a Antenas) Antinodes(limit plane.Vector, single bool) []plane.Vector {
	var anti []plane.Vector
	for ch := range a {
		anti = append(anti, a.antinodes(a[ch], limit, single)...)
	}
	return anti
}

func (a Antenas) antinodes(nodes []plane.Vector, limit plane.Vector, single bool) []plane.Vector {
	var anti []plane.Vector
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if !single {
				anti = append(anti, nodes[i], nodes[j])
			}

			ij := nodes[j].Minus(nodes[i])
			ji := nodes[i].Minus(nodes[j])

			ia := nodes[i].Add(ji)
			for ia.Within(limit) {
				anti = append(anti, ia)
				if single {
					break
				}
				ia = ia.Add(ji)
			}

			ja := nodes[j].Add(ij)
			for ja.Within(limit) {
				anti = append(anti, ja)
				if single {
					break
				}
				ja = ja.Add(ij)
			}
		}
	}
	return anti
}
