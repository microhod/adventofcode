package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

// direction enum
const (
	X Direction = iota
	Y
)

func main() {
	puzzle.NewSolution("Transparent Origami", part1, part2).Run()
}

func part1() error {
	man, err := readManual(InputFile)
	if err != nil {
		return err
	}

	man.Fold(man.Folds[0])

	dots := 0
	for _, row := range man.Dots {
		for _, dot := range row {
			dots += dot
		}
	}
	fmt.Printf("# dots after 1 fold: %d\n", dots)

	return nil
}

func part2() error {
	man, err := readManual(InputFile)
	if err != nil {
		return err
	}

	for _, fold := range man.Folds {
		man.Fold(fold)
	}

	fmt.Println(man)

	return nil
}

func readManual(path string) (*Manual, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(bytes), "\n\n")
	if len(parts) < 2 {
		return nil, fmt.Errorf("expected 2 sections separated by a blank line only got %d", len(parts))
	}

	dots, err := readDots(strings.Split(parts[0], "\n"))
	if err != nil {
		return nil, err
	}

	folds, err := readFolds(strings.Split(parts[1], "\n"))
	if err != nil {
		return nil, err
	}

	return &Manual{
		Dots:  dots,
		Folds: folds,
	}, nil
}

func readDots(lines []string) ([][]int, error) {
	rows := 0
	cols := 0
	coords := [][]int{}
	for _, line := range lines {
		coord, err := csv.ParseInts(line)
		if err != nil {
			return nil, err
		}

		rows = max(rows, coord[1]+1)
		cols = max(cols, coord[0]+1)

		coords = append(coords, coord)
	}

	grid := [][]int{}
	for i := 0; i < rows; i++ {
		grid = append(grid, make([]int, cols))
	}

	for _, coord := range coords {
		grid[coord[1]][coord[0]] = 1
	}

	return grid, nil
}

func readFolds(lines []string) ([]Fold, error) {
	folds := []Fold{}
	for _, line := range lines {
		direction := X
		if strings.Contains(line, "y") {
			direction = Y
		}

		if parts := strings.Split(line, "="); len(parts) > 1 {
			index, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				return nil, err
			}

			folds = append(folds, Fold{
				Direction: direction,
				Index:     index,
			})
		}
	}
	return folds, nil
}

type Manual struct {
	Dots  [][]int
	Folds []Fold
}

func (m *Manual) Fold(fold Fold) {
	switch fold.Direction {
	case X:
		m.FoldX(fold.Index)
	case Y:
		m.FoldY(fold.Index)
	default:
		return
	}
}

func (m *Manual) FoldX(index int) {
	right := [][]int{}
	for _, row := range m.Dots {
		right = append(right, row[index+1:])
	}

	for i, row := range m.Dots {
		m.Dots[i] = row[:index]
	}

	for r, row := range right {
		for c, dot := range row {
			m.Dots[r][len(m.Dots[0])-1-c] = max(m.Dots[r][len(m.Dots[0])-1-c], dot)
		}
	}
}

func (m *Manual) FoldY(index int) {
	lower := m.Dots[index+1:]
	m.Dots = m.Dots[:index]

	for r, row := range lower {
		for c, dot := range row {
			m.Dots[len(m.Dots)-1-r][c] = max(m.Dots[len(m.Dots)-1-r][c], dot)
		}
	}
}

func (m Manual) String() string {
	str := ""
	for _, line := range m.Dots {
		for _, dot := range line {
			if dot == 1 {
				str += "#"
			} else {
				str += "."
			}
		}
		str += "\n"
	}

	for _, fold := range m.Folds {
		str += fmt.Sprintln(fold)
	}
	return str
}

type Fold struct {
	Direction Direction
	Index     int
}

func (f Fold) String() string {
	return fmt.Sprintf("%s=%d", f.Direction.String(), f.Index)
}

type Direction int

func (d Direction) String() string {
	switch d {
	case X:
		return "x"
	case Y:
		return "y"
	default:
		return fmt.Sprintf("%d", d)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
