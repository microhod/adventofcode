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
	puzzle.NewSolution("Point of Incidence", part1, part2).Run()
}

func part1() error {
	patterns, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	var notes int
	for _, p := range patterns {		
		row, col := p.FindMirrors(-1, -1)
		notes += (100*row) + col
	}
	fmt.Printf("summarisation of notes: %d\n", notes)
	return nil
}

func part2() error {
	patterns, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	var notes int
	for _, p := range patterns {		
		row, col := p.FindMirrors(-1, -1)
		row, col = p.FindSmudge(row, col)
		notes += (100*row) + col
	}
	fmt.Printf("summarisation of notes: %d\n", notes)
	return nil
}

func parseInput(path string) ([]Pattern, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, err
	}
	var patterns []Pattern
	for _, part := range strings.Split(input, "\n\n") {
		part = strings.TrimSpace(part)
		patterns = append(patterns, strings.Split(part, "\n"))
	}
	return patterns, nil
}

type Pattern []string

func (p Pattern) FindSmudge(row, col int) (int, int) {
	swap := func(ch byte) byte {
		if ch == '.' {
			return '#'
		}
		return '.'
	}

	for i := range p {
		for j := range p[i] {
			p[i] = p[i][:j] + string(swap(p[i][j])) + p[i][j+1:]

			newRow, newCol := p.FindMirrors(row, col)
			if newRow + newCol > 0 && (newRow != row || newCol != col) {
				return newRow, newCol
			}

			p[i] = p[i][:j] + string(swap(p[i][j])) + p[i][j+1:]
		}
	}
	return 0, 0
}

func (p Pattern) FindMirrors(skipRow, skipCol int) (int, int) {
	var rows, cols int
	for mirror := 1; mirror < len(p); mirror++ {
		if skipRow == mirror {
			continue
		}
		if p.IsMirrorRow(mirror) {
			rows += mirror
		}
	}
	for mirror := 1; mirror < len(p[0]); mirror++ {
		if skipCol == mirror {
			continue
		}
		if p.IsMirrorCol(mirror) {
			cols += mirror
		}
	}
	return rows, cols
}

func (p Pattern) IsMirrorRow(mirror int) bool {
	l := p[:mirror]
	r := p[mirror:]

	for row := 0; row < maths.Min(len(l), len(r)); row++ {
		lRow := l[len(l)-1-row]
		rRow := r[row]

		if lRow != rRow {
			return false
		}
	}
	return true
}

func (p Pattern) IsMirrorCol(mirror int) bool {
	return p.Transpose().IsMirrorRow(mirror)
}

func (p Pattern) Transpose() Pattern {
	var transposed Pattern
	for col := range p[0] {
		var tRow string
		for row := range p {
			tRow += string(p[row][col])
		}
		transposed = append(transposed, tRow)
	}
	return transposed
}
