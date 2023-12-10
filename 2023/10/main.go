package main

import (
	"fmt"
	"sort"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/geometry/plane"
	"github.com/microhod/adventofcode/internal/graph"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Pipe Maze", part1, part2).Run()
}

func part1() error {
	grid, start, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	fmt.Printf("the farthest point from the start is: %d\n", MaxDistance(grid.Loop(start), start))
	return nil
}

func part2() error {
	grid, start, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	// 1. find the loop
	loop := grid.Loop(start).Nodes()
	
	// 2. replace all non-loop pipes with ground
	grid.ReplaceNonLoopTiles(loop, '.')
	
	// 3. expand any possible 'squeezable' rows/columns
	grid, rows, cols := grid.Expand()

	// 4. do a 'flood fill' to find all tiles which are outside
	outside := grid.Outside(loop)
	for v := range outside {
		grid[v.Y][v.X] = 'O'
	}

	// 5. unexpand the previously expanded rows & columns
	grid = grid.Unexpand(rows, cols)

	// 6. count the remaining number of ground tiles	
	fmt.Printf("the number of inside tiles is: %d\n", grid.CountGround())
	return nil
}

func parseInput(path string) (Grid, plane.Vector, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, plane.Vector{}, err
	}

	grid := make(Grid, len(lines))
	for y, line := range lines {
		grid[y] = append(grid[y], '.')
		for _, r := range line {
			grid[y] = append(grid[y], Tile(r))
		}
		grid[y] = append(grid[y], '.')
	}
	
	ground := make([]Tile, len(grid[0]))
	for i := range ground {
		ground[i] = '.'
	}
	grid = append([][]Tile{ground}, grid...)
	grid = append(grid, ground)
	
	var start plane.Vector
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == 'S' {
				start = plane.Vector{X: x, Y: y}
			}
		}
	}

	return grid, start, nil
}

func MaxDistance(loop graph.Graph[plane.Vector], start plane.Vector) int {
	dist := map[plane.Vector]int{
		start: 0,
	}
	max := 0
	stack := []plane.Vector{start}
	for len(stack) > 0 {
		u := stack[0]
		stack = stack[1:]

		for v := range loop[u] {
			if _, exists := dist[v]; exists {
				continue
			}
			dist[v] = dist[u] + 1
			max = maths.Max(max, dist[v])
			stack = append(stack, v)
		}
	}
	return max
}

type Grid [][]Tile

func (g Grid) Loop(start plane.Vector) graph.Graph[plane.Vector] {
	loop := graph.NewGraph[plane.Vector]()

	stack := []plane.Vector{start}
	seen := set.NewSet[plane.Vector]()

	for len(stack) > 0 {
		v := stack[0]
		stack = stack[1:]
		tile := g[v.Y][v.X]

		if seen.Contains(v) {
			continue
		}
		seen.Add(v)

		neighbours := v.OrthogonalNeighbours()
		for direction := range tileDirections[tile] {
			neighbour := neighbours[direction]
			neighbourTile := g[neighbour.Y][neighbour.X]
			if neighbourTile.IsGround() {
				continue
			}
			if !tileDirections[neighbourTile].Contains(direction.Opposite()) {
				continue
			}

			loop.AddEdge(v, neighbour, 1)
			loop.AddEdge(neighbour, v, 1)

			stack = append(stack, neighbour)
		}
	}
	return loop
}

func (g Grid) CountGround() int {
	var ground int
	for y := 0; y < len(g); y++ {
		for x := 0; x < len(g[0]); x++ {
			if g[y][x].IsGround() {
				ground++
			}
		}
	}
	return ground
}

func (g Grid) ReplaceNonLoopTiles(loop set.Set[plane.Vector], tile Tile) {
	for y := 1; y < len(g)-1; y++ {
		for x := 1; x < len(g[0])-1; x++ {
			if loop.Contains(plane.Vector{X: x, Y: y}) {
				continue
			}
			g[y][x] = tile
		}
	}
}

func (g Grid) Expand() (Grid, []int, []int) {
	rows := g.FindSqueezableRows()
	cols := g.FindSqueezableColumns()

	for i := len(rows) - 1; i >= 0; i-- {
		g = g.ExpandRow(rows[i])
	}
	for i := len(cols) - 1; i >= 0; i-- {
		g = g.ExpandColumn(cols[i])
	}
	return g, rows, cols
}

func (g Grid) FindSqueezableRows() []int {
	rows := set.NewSet[int]()
	for y := 1; y < len(g)-2; y++ {
		for x := 1; x < len(g[y])-1; x++ {
			if g[y][x].CanSqueezeBetween(g[y+1][x], plane.South) {
				rows.Add(y)
				break
			}
		}
	}
	slice := rows.ToSlice()
	sort.Ints(slice)
	return slice
}

func (g Grid) FindSqueezableColumns() []int {
	columns := set.NewSet[int]()
	for x := 1; x < len(g[0])-2; x++ {
		for y := 1; y < len(g)-1; y++ {
			if g[y][x].CanSqueezeBetween(g[y][x+1], plane.East) {
				columns.Add(x)
				break
			}
		}
	}
	slice := columns.ToSlice()
	sort.Ints(slice)
	return slice
}

func (g Grid) ExpandRow(row int) Grid {
	middle := make([]Tile, len(g[row]))
	for x := range middle {
		middle[x] = g[row][x].Between(g[row+1][x], plane.South)
	}

	return append(g[:row+1], append([][]Tile{middle}, g[row+1:]...)...)
}

func (g Grid) ExpandColumn(col int) Grid {
	for y := 0; y < len(g); y++ {
		middle := g[y][col].Between(g[y][col+1], plane.East)
		g[y] = append(g[y][:col+1], append([]Tile{middle}, g[y][col+1:]...)...)
	}
	return g
}

func (g Grid) Unexpand(rows, cols []int) Grid {
	for _, row := range rows {
		g = append(g[:row+1], g[row+2:]...)
	}
	for _, col := range cols {
		for y := range g {
			g[y] = append(g[y][:col+1], g[y][col+2:]...)
		}
	}
	return g
}

func (g Grid) Outside(loop set.Set[plane.Vector]) set.Set[plane.Vector] {
	start := plane.Vector{X: 0, Y: 0}
	seen := set.NewSet[plane.Vector]()
	outside := set.NewSet[plane.Vector](start)
	stack := []plane.Vector{start}

	for len(stack) > 0 {
		v := stack[0]
		stack = stack[1:]

		if seen.Contains(v) {
			continue
		}
		seen.Add(v)

		for _, n := range v.Neighbours() {
			if !g.Valid(n) {
				continue
			}
			if !g[n.Y][n.X].IsGround() {
				continue
			}
			outside.Add(n)
			stack = append(stack, n)
		}
	}
	return outside
}

func (g Grid) Valid(v plane.Vector) bool {
	if v.X < 0 || v.X >= len(g[0]) {
		return false
	}
	if v.Y < 0 || v.Y >= len(g) {
		return false
	}
	return true
}

var tileDirections = map[Tile]set.Set[plane.Direction]{
	'|': set.NewSet(plane.North, plane.South),
	'-': set.NewSet(plane.East, plane.West),
	'L': set.NewSet(plane.North, plane.East),
	'J': set.NewSet(plane.North, plane.West),
	'7': set.NewSet(plane.South, plane.West),
	'F': set.NewSet(plane.South, plane.East),
	'.': set.NewSet[plane.Direction](),
	'S': set.NewSet(plane.North, plane.East, plane.South, plane.West),
}

type Tile rune

func (t Tile) ConnectedTo(s Tile, d plane.Direction) bool {
	return tileDirections[t].Contains(d) && tileDirections[s].Contains(d.Opposite())
}

func (t Tile) Between(s Tile, d plane.Direction) Tile {
	middle := Tile('.')
	if t.ConnectedTo(s, d) {
		switch d {
		case plane.East, plane.West:
			middle = '-'
		case plane.North, plane.South:
			middle = '|'
		}
	}
	return middle
}

func (t Tile) CanSqueezeBetween(s Tile, d plane.Direction) bool {
	if t.IsGround() || s.IsGround() {
		return false
	}
	if t.ConnectedTo(s, d) {
		return false
	}
	return true
}

func (t Tile) IsGround() bool {
	return t == '.'
}
