package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/geometry/plane"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("The Floor Will Be Lava", part1, part2).Run()
}

func part1() error {
	contraption, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	start := Beam{
		Position:  plane.Vector{X: 0, Y: 0},
		Direction: plane.East,
	}
	energised := contraption.Energise(start)

	fmt.Printf("the number of tiles energised is: %d\n", energised)
	return nil
}

func part2() error {
	contraption, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	var max int
	for _, start := range contraption.AllStarts() {
		max = maths.Max(max, contraption.Energise(start))
	}

	fmt.Printf("the number of tiles energised is: %d\n", max)
	return nil
}

func parseInput(path string) (Contraption, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	contraption := make(Contraption)
	for y := range lines {
		for x, r := range lines[y] {
			contraption[plane.Vector{X: x, Y: y}] = Tile(r)
		}
	}
	return contraption, err
}

type Contraption map[plane.Vector]Tile

func (c Contraption) Energise(start Beam) int {
	beams := []Beam{start}
	seen := set.NewSet[Beam]()

	for len(beams) > 0 {
		beam := beams[0]
		beams = beams[1:]

		if seen.Contains(beam) {
			continue
		}
		seen.Add(beam)

		beams = append(beams, c.next(beam)...)
	}

	energised := set.NewSet[plane.Vector]()
	for seen := range seen {
		energised.Add(seen.Position)
	}

	return len(energised)
}

func (c Contraption) next(beam Beam) []Beam {
	nexts := set.NewSet[Beam]()

	switch c[beam.Position] {
	case SplitVertical:
		if beam.Direction == plane.North || beam.Direction == plane.South {
			break
		}
		nexts.Add(Beam{
			Position:  beam.Position.Add(plane.DirectionToVector[plane.North]),
			Direction: plane.North,
		}, Beam{
			Position:  beam.Position.Add(plane.DirectionToVector[plane.South]),
			Direction: plane.South,
		})
	case SplitHorizontal:
		if beam.Direction == plane.East || beam.Direction == plane.West {
			break
		}
		nexts.Add(Beam{
			Position:  beam.Position.Add(plane.DirectionToVector[plane.East]),
			Direction: plane.East,
		}, Beam{
			Position:  beam.Position.Add(plane.DirectionToVector[plane.West]),
			Direction: plane.West,
		})
	case SlashMirror:
		beam.Direction = map[plane.Direction]plane.Direction{
			plane.North: plane.East,
			plane.East:  plane.North,
			plane.South: plane.West,
			plane.West:  plane.South,
		}[beam.Direction]
		nexts.Add(beam.Move())
	case BackSlashMirror:
		beam.Direction = map[plane.Direction]plane.Direction{
			plane.North: plane.West,
			plane.West:  plane.North,
			plane.South: plane.East,
			plane.East:  plane.South,
		}[beam.Direction]
		nexts.Add(beam.Move())
	}

	if len(nexts) == 0 {
		nexts.Add(beam.Move())
	}

	var result []Beam
	for next := range nexts {
		if _, exists := c[next.Position]; !exists {
			continue
		}
		result = append(result, next)
	}
	return result
}

func (c Contraption) AllStarts() []Beam {
	maxX := 0
	maxY := 0
	for v := range c {
		maxX = maths.Max(maxX, v.X)
		maxY = maths.Max(maxY, v.Y)
	}

	var starts []Beam
	// top & bottom
	for x := 0; x <= maxX; x++ {
		starts = append(starts, Beam{
			Position:  plane.Vector{X: x, Y: 0},
			Direction: plane.South,
		}, Beam{
			Position:  plane.Vector{X: x, Y: maxY},
			Direction: plane.North,
		})
	}
	// left & right
	for y := 0; y <= maxY; y++ {
		starts = append(starts, Beam{
			Position:  plane.Vector{X: 0, Y: y},
			Direction: plane.East,
		}, Beam{
			Position:  plane.Vector{X: maxX, Y: y},
			Direction: plane.West,
		})
	}
	return starts
}

type Beam struct {
	Position  plane.Vector
	Direction plane.Direction
}

func (b Beam) Move() Beam {
	b.Position = b.Position.Add(plane.DirectionToVector[b.Direction])
	return b
}

type Tile rune

const (
	Empty           Tile = '.'
	SplitVertical   Tile = '|'
	SplitHorizontal Tile = '-'
	SlashMirror     Tile = '/'
	BackSlashMirror Tile = '\\'
)
