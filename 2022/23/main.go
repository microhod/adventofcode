package main

import (
	"fmt"
	"math"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Unstable Diffusion", part1, part2).Run()
}

func part1() error {
	grove, err := parse(InputFile)
	if err != nil {
		return err
	}

	for round := 0; round < 10; round++ {
		grove.MoveEvles()
	}

	fmt.Println(grove.EmptySpaces())
	return nil
}

func part2() error {
	grove, err := parse(InputFile)
	if err != nil {
		return err
	}

	round := 1
	for grove.MoveEvles() {
		round += 1
	}

	fmt.Println(round)
	return nil
}

func parse(path string) (*Grove, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	grove := &Grove{
		elves:      map[Vector]*Elf{},
		directions: []Direction{North, South, West, East},
	}
	for y, line := range lines {
		if line == "" {
			continue
		}

		for x, char := range line {
			if char != '#' {
				continue
			}

			elf := Elf(len(grove.elves))
			grove.elves[Vector{x, y}] = &elf
		}
	}
	return grove, nil
}

type Grove struct {
	elves      map[Vector]*Elf
	directions []Direction
}

// MoveEvles runs a round of elf movement and returns a boolean
// indicating whether any elves moved or not
func (g *Grove) MoveEvles() bool {
	proposals := g.proposeNextPositions()
	if len(proposals) == 0 {
		return false
	}
	for next, elves := range proposals {
		if len(elves) != 1 {
			continue
		}
		elf := elves[0]

		// move elf to next position
		g.elves[next] = g.elves[elf]
		// delete old position
		delete(g.elves, elf)
	}

	// shift directions
	g.directions = append(g.directions[1:], g.directions[0])
	return true
}

func (g *Grove) proposeNextPositions() map[Vector][]Vector {
	proposals := map[Vector][]Vector{}
	for position := range g.elves {
		next := g.proposeNextPosition(position)
		if next != position {
			proposals[next] = append(proposals[next], position)
		}
	}
	return proposals
}

func (g *Grove) proposeNextPosition(current Vector) Vector {
	neighbours := current.Neighbours()

	var proposals []Vector

	for _, direction := range g.directions {
		previousDiagonal := Direction(maths.Mod(int(direction)-1, 8))
		nextDiagonal := Direction(maths.Mod(int(direction)+1, 8))

		noElves := g.elves[neighbours[direction]] == nil
		noElves = noElves && g.elves[neighbours[previousDiagonal]] == nil
		noElves = noElves && g.elves[neighbours[nextDiagonal]] == nil

		if noElves {
			proposals = append(proposals, current.Add(DirectionToVector[direction]))
		}
	}

	if len(proposals) == 0 || len(proposals) == len(g.directions) {
		return current
	}
	return proposals[0]
}

func (g *Grove) BoundingRectangle() (int, int, int, int) {
	minX, minY := math.MaxInt, math.MaxInt
	maxX, maxY := math.MinInt, math.MinInt
	for position := range g.elves {
		if position.X < minX {
			minX = position.X
		}
		if position.X > maxX {
			maxX = position.X
		}
		if position.Y < minY {
			minY = position.Y
		}
		if position.Y > maxY {
			maxY = position.Y
		}
	}
	return minX, maxX, minY, maxY
}

func (g *Grove) EmptySpaces() int {
	minX, maxX, minY, maxY := g.BoundingRectangle()
	totalSpaces := maths.Abs((1 + maxX - minX) * (1 + maxY - minY))

	return totalSpaces - len(g.elves)
}

type Elf int

type Vector struct {
	X, Y int
}

func (v Vector) Add(u Vector) Vector {
	return Vector{v.X + u.X, v.Y + u.Y}
}

type Direction int

const (
	North     Direction = 0
	NorthEast Direction = 1
	East      Direction = 2
	SouthEast Direction = 3
	South     Direction = 4
	SouthWest Direction = 5
	West      Direction = 6
	NorthWest Direction = 7
)

// 7 0 1
// 6   2
// 5 4 3
var DirectionToVector = map[Direction]Vector{
	North:     {0, -1},
	NorthEast: {1, -1},
	East:      {1, 0},
	SouthEast: {1, 1},
	South:     {0, 1},
	SouthWest: {-1, 1},
	West:      {-1, 0},
	NorthWest: {-1, -1},
}

func (v Vector) Neighbours() map[Direction]Vector {
	neighbours := map[Direction]Vector{}
	for direction, diff := range DirectionToVector {
		neighbours[direction] = v.Add(diff)
	}
	return neighbours
}
