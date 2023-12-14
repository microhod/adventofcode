package main

import (
	"fmt"
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
	puzzle.NewSolution("Parabolic Reflector Dish", part1, part2).Run()
}

func part1() error {
	platform, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	platform.TiltNorth()

	fmt.Printf("the load on the north support is: %d\n", platform.LoadOnNorthSupport())
	return nil
}

func part2() error {
	platform, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	cycles := 1000000000

	// find the periodicity of the spins so we can skip most iterations
	seen := make(map[string]int)
	remaining := 0
	for i := 1; i <= cycles; i++ {
		platform.Spin()

		if index, exists := seen[platform.String()]; exists {
			periodicity := i - index
			// take the remaining section and remove all muliples of the periodicity
			remaining = maths.Mod(cycles-i, periodicity)
			break
		}
		seen[platform.String()] = i
	}
	// do the remaining spins
	for j := 1; j <= remaining; j++ {
		platform.Spin()
	}
	fmt.Printf("the load on the north support is: %d\n", platform.LoadOnNorthSupport())
	return nil
}

func parseInput(path string) (*Platform, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	platform := new(Platform)
	for _, line := range lines {
		platform.Rocks = append(platform.Rocks, []Rock(line))
	}
	return platform, nil
}

type Platform struct {
	Rocks [][]Rock
}

func (p *Platform) Spin() {
	p.TiltNorth()
	p.TiltWest()
	p.TiltSouth()
	p.TiltEast()
}

func (p *Platform) TiltNorth() {
	for r := range p.Rocks {
		for c := range p.Rocks[r] {
			if p.Rocks[r][c] == Round {
				row := r
				for row > 0 {
					if p.Rocks[row-1][c] != Empty {
						break
					}

					p.Rocks[row-1][c] = p.Rocks[row][c]
					p.Rocks[row][c] = Empty
					row -= 1
				}
			}
		}
	}
}

func (p *Platform) TiltWest() {
	for c := range p.Rocks[0] {
		for r := range p.Rocks {
			if p.Rocks[r][c] == Round {
				col := c
				for col > 0 {
					if p.Rocks[r][col-1] != Empty {
						break
					}

					p.Rocks[r][col-1] = p.Rocks[r][col]
					p.Rocks[r][col] = Empty
					col -= 1
				}
			}
		}
	}
}

func (p *Platform) TiltSouth() {
	for r := len(p.Rocks) - 1; r >= 0; r-- {
		for c := len(p.Rocks[r]) - 1; c >= 0; c-- {
			if p.Rocks[r][c] == Round {
				row := r
				for row < len(p.Rocks)-1 {
					if p.Rocks[row+1][c] != Empty {
						break
					}

					p.Rocks[row+1][c] = p.Rocks[row][c]
					p.Rocks[row][c] = Empty
					row += 1
				}
			}
		}
	}
}

func (p *Platform) TiltEast() {
	for c := len(p.Rocks[0]) - 1; c >= 0; c-- {
		for r := len(p.Rocks) - 1; r >= 0; r-- {
			if p.Rocks[r][c] == Round {
				col := c
				for col < len(p.Rocks[0])-1 {
					if p.Rocks[r][col+1] != Empty {
						break
					}

					p.Rocks[r][col+1] = p.Rocks[r][col]
					p.Rocks[r][col] = Empty
					col += 1
				}
			}
		}
	}
}

func (p *Platform) LoadOnNorthSupport() int {
	var load int
	for row := range p.Rocks {
		for col := range p.Rocks[row] {
			if p.Rocks[row][col] == Round {
				load += len(p.Rocks) - row
			}
		}
	}
	return load
}

func (p *Platform) String() string {
	var lines []string
	for row := range p.Rocks {
		var line string
		for _, r := range p.Rocks[row] {
			line += string(r)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n") + "\n"
}

type Rock rune

const (
	Empty Rock = '.'
	Round Rock = 'O'
	Cube  Rock = '#'
)
