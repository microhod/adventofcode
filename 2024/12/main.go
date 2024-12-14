package main

import (
	"fmt"
	"strings"

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
	puzzle.NewSolution("Garden Groups", part1, part2).Run()
}

func part1() error {
	m, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Println(m.PerimeterCost())
	return nil
}

func part2() error {
	m, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Println(m.SidePrice())
	return nil
}

func parse(path string) (Map, error) {
	m, _, err := file.ReadVectorsFunc(path, func(b byte) bool {
		return true
	})
	return m, err
}

type Map map[byte][]plane.Vector

func (m Map) PerimeterCost() int {
	var cost int
	for _, vectors := range m {
		for _, region := range FindAllRegions(set.NewSet(vectors...)) {
			cost += region.area * region.perimeter
		}
	}
	return cost
}

func (m Map) SidePrice() int {
	var cost int
	for _, vectors := range m {
		for _, region := range FindAllRegions(set.NewSet(vectors...)) {
			cost += region.area * region.sides
		}
	}
	return cost
}

type Region struct {
	vectors   set.Set[plane.Vector]
	area      int
	perimeter int
	sides     int
}

func FindAllRegions(vectors set.Set[plane.Vector]) []Region {
	var regions []Region
	for {
		start, ok := vectors.Pop()
		if !ok {
			break
		}

		region := NewRegion(vectors, start)
		regions = append(regions, region)

		vectors = set.Diff(vectors, region.vectors)
	}
	return regions
}

func NewRegion(vectors set.Set[plane.Vector], start plane.Vector) Region {
	region := Region{
		vectors: set.NewSet[plane.Vector](),
	}
	visit := []plane.Vector{start}
	for len(visit) > 0 {
		v := visit[0]
		visit = visit[1:]

		if region.vectors.Contains(v) {
			continue
		}
		region.vectors.Add(v)
		region.area += 1
		region.perimeter += 4

		// calculate sides
		var neighbours []plane.Vector
		for _, n := range v.Neighbours() {
			if !vectors.Contains(n) {
				continue
			}
			if !region.vectors.Contains(n) {
				continue
			}
			neighbours = append(neighbours, n.Minus(v))
		}
		withV := set.NewSet(neighbours...)
		withV.Add(plane.Vector{X: 0, Y: 0})
		region.sides += numSides(withV) - numSides(set.NewSet(neighbours...))

		for _, n := range v.OrthogonalNeighbours() {
			if !vectors.Contains(n) {
				continue
			}
			if region.vectors.Contains(n) {
				region.perimeter -= 2
				continue
			}
			visit = append(visit, n)
		}
	}

	return region
}

func numSides(shape set.Set[plane.Vector]) int {
	shape = set.NewSet(shape.ToSlice()...)

	// this is the only case which has an inner gap
	// which won't work with our algorithm
	// ###
	// #.#
	// ###
	if len(shape) == 8 {
		return 8
	}

	var total int
	for len(shape) > 0 {
		sides, region := numConnectedRegionSides(shape)
		total += sides
		shape = set.Diff(shape, region)
	}
	return total
}

// pick a start and "walk" around the shape counting sides by counting each time we have
// to turn on our walk
func numConnectedRegionSides(shape set.Set[plane.Vector]) (int, set.Set[plane.Vector]) {
	start := pickFirstNeightbour(shape)
	startDir := plane.East

	dir := startDir
	v := start.AddDirection(dir)

	region := set.NewSet[plane.Vector]()
	seen := set.NewSet[[3]int]()

	var sides int
	for {
		state := [3]int{v.X, v.Y, int(dir)}
		if seen.Contains(state) {
			break
		}
		seen.Add(state)

		if shape.Contains(v) {
			region.Add(v)

			// if we turn anti-clockwise and see a part of the shape, we've gone inside
			if shape.Contains(v.AddDirection(dir.Turn(-90))) {
				dir = dir.Turn(-90)
				v = v.AddDirection(dir)
				sides += 1
				continue
			}

			v = v.AddDirection(dir)
			continue
		}

		v = v.MinusDirection(dir)
		dir = dir.Turn(90)
		sides += 1
	}
	return sides, region
}

func pickFirstNeightbour(shape set.Set[plane.Vector]) plane.Vector {
	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			if shape.Contains(plane.Vector{X: x, Y: y}) {
				return plane.Vector{X: x, Y: y}
			}
		}
	}
	panic("no first found")
} 

func drawNeighbourhood(square set.Set[plane.Vector]) string {
	lines := make([]string, 3)
	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			if square.Contains(plane.Vector{X: x, Y: y}) {
				lines[y+1] += "#"
				continue
			}
			lines[y+1] += "."
		}
	}
	return strings.Join(lines, "\n")
}
