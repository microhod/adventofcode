package main

import (
	"errors"
	"fmt"
	"slices"

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
	puzzle.NewSolution("Guard Gallivant", part1, part2).Run()
}

func part1() error {
	guard, m, err := parse(InputFile)
	if err != nil {
		return err
	}
	path, err := guard.Walk(m)
	if err != nil {
		return err
	}
	pos := set.NewSet[plane.Vector]()
	for _, p := range path {
		pos.Add(p.pos)
	}

	fmt.Printf("total guard positions: %d\n", len(pos))
	return nil
}

func part2() error {
	guard, m, err := parse(InputFile)
	if err != nil {
		return err
	}
	blocks, err := guard.Block(m)
	if err != nil {
		return err
	}

	fmt.Printf("total blocking positions: %d\n", len(blocks))
	return nil
}

func parse(path string) (*Guard, *Map, error) {
	guard, limit, err := file.ReadVectors(path, '^')
	if err != nil {
		return nil, nil, err
	}

	obstacles, _, err := file.ReadVectors(path, '#')
	if err != nil {
		return nil, nil, err
	}

	return &Guard{
			dir: plane.North,
			pos: guard[0],
		}, &Map{
			obstacles: set.NewSet(obstacles...),
			limit:     limit,
		}, nil
}

type Map struct {
	obstacles set.Set[plane.Vector]
	limit     plane.Vector
}

func (m Map) Block(v plane.Vector) *Map {
	obstacles := set.NewSet(m.obstacles.ToSlice()...)
	obstacles.Add(v)

	return &Map{obstacles: obstacles, limit: m.limit}
}

func (m Map) Draw(g *Guard, path []plane.Vector) string {
	dir := map[plane.Direction]byte{
		plane.North: '^',
		plane.East:  '>',
		plane.West:  '<',
		plane.South: 'v',
	}[g.dir]

	return plane.Draw(map[byte][]plane.Vector{
		'#': m.obstacles.ToSlice(),
		dir: {g.pos},
		'X': path,
	}, m.limit)
}

type Guard struct {
	dir plane.Direction
	pos plane.Vector
}

var ErrCycle = errors.New("cycle in path")

func (g *Guard) Walk(m *Map) ([]Move, error) {
	var path []Move

	out := func(p plane.Vector) bool {
		return p.X < 0 || p.Y < 0 || p.X > m.limit.X || p.Y > m.limit.Y
	}

	for !out(g.pos) {
		move := Move{pos: g.pos, dir: g.dir}
		if slices.Contains(path, move) {
			path = append(path, move)
			return path, ErrCycle
		}
		path = append(path, move)

		// move
		next := g.pos.Add(plane.DirectionToVector[g.dir])
		if m.obstacles.Contains(next) {
			g.dir = g.dir.Turn(90)
			continue
		}
		g.pos = next
	}

	return path, nil
}

func (g *Guard) Block(m *Map) ([]plane.Vector, error) {
	startPos := g.pos
	startDir := g.dir

	path, err := g.Walk(m)
	if err != nil {
		return nil, err
	}

	blocks := set.NewSet[plane.Vector]()
	for _, move := range path {
		if move.pos == startPos {
			continue
		}
		editied := m.Block(move.pos)
		// reset guard
		g.pos = startPos
		g.dir = startDir

		_, err := g.Walk(editied)
		if errors.Is(err, ErrCycle) {
			blocks.Add(move.pos)
		}
	}

	return blocks.ToSlice(), nil
}

type Move struct {
	pos plane.Vector
	dir plane.Direction
}
