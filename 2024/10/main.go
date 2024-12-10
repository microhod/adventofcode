package main

import (
	"fmt"
	"strconv"

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
	puzzle.NewSolution("Hoof It", part1, part2).Run()
}

func part1() error {
	m, err := parse(InputFile)
	if err != nil {
		return err
	}

	var total int
	for _, score := range m.TrailHeadScores() {
		total += score
	}
	fmt.Printf("scores of all trailheads: %d\n", total)
	return nil
}

func part2() error {
	m, err := parse(InputFile)
	if err != nil {
		return err
	}

	var total int
	for _, score := range m.TrailHeadRatings() {
		total += score
	}
	fmt.Printf("ratings of all trailheads: %d\n", total)
	return nil
}

func parse(path string) (Map, error) {
	return  file.ReadVectorMapFunc(path, func(b byte) (int, error) {
		return strconv.Atoi(string(b))
	})
}

type Map map[plane.Vector]int

func (m Map) TrailHeadScores() map[plane.Vector]int {
	scores := make(map[plane.Vector]int)
	for _, head := range m.findTrailHeads() {
		summits := set.NewSet[plane.Vector]()
		m.findSummits(head, summits)
		scores[head] = len(summits)
	}
	return scores
}

func (m Map) findSummits(trail plane.Vector, summits set.Set[plane.Vector]) {
	if m[trail] == 9 {
		summits.Add(trail)
	}

	var next []plane.Vector
	for _, v := range trail.OrthogonalNeighbours() {
		if _, ok := m[v]; !ok {
			continue
		}
		if m[v] == m[trail]+1 {
			next = append(next, v)
		}
	}

	for _, n := range next {
		m.findSummits(n, summits)
	}
}

func (m Map) TrailHeadRatings() map[plane.Vector]int {
	ratings := make(map[plane.Vector]int)
	for _, head := range m.findTrailHeads() {
		ratings[head] += m.findTrails(head)
	}
	return ratings
}

func (m Map) findTrails(trail plane.Vector) int {
	if m[trail] == 9 {
		return 1
	}

	var next []plane.Vector
	for _, v := range trail.OrthogonalNeighbours() {
		if _, ok := m[v]; !ok {
			continue
		}
		if m[v] == m[trail]+1 {
			next = append(next, v)
		}
	}

	var rating int
	for _, n := range next {
		rating += m.findTrails(n)
	}
	return rating
}

func (m Map) findTrailHeads() []plane.Vector {
	var trailheads []plane.Vector
	for v, h := range m {
		if h == 0 {
			trailheads = append(trailheads, v)
		}
	}
	return trailheads
}
