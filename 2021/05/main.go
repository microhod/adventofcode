package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	LinesFile = "lines.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Hydrothermal Venture", part1, part2).Run()
}

func part1() error {
	lines, err := readLines(LinesFile)
	if err != nil {
		return err
	}

	lines = filterOutDiagonals(lines)

	fmt.Printf("# overlaps (ignoring diagonals): %d\n", len(findOverlaps(lines)))

	return nil
}

func part2() error {
	lines, err := readLines(LinesFile)
	if err != nil {
		return err
	}

	fmt.Printf("# overlaps: %d\n", len(findOverlaps(lines)))

	return nil
}

type Point struct {
	X, Y int
}

type Line struct {
	Start, End Point
}

func (line Line) Points() []Point {
	points := []Point{}
	if line.Start.X == line.End.X {
		start := min(line.Start.Y, line.End.Y)
		end := max(line.Start.Y, line.End.Y)

		for i := start; i < end+1; i++ {
			points = append(points, Point{
				X: line.Start.X,
				Y: i,
			})
		}
	} else if line.Start.Y == line.End.Y {
		start := min(line.Start.X, line.End.X)
		end := max(line.Start.X, line.End.X)

		for i := start; i < end+1; i++ {
			points = append(points, Point{
				X: i,
				Y: line.Start.Y,
			})
		}
	} else {
		// assume 45 degrees so difference between x and y is the same
		diff := diff(line.Start.X, line.End.X)

		for i := 0; i < diff+1; i++ {
			x := line.Start.X + i
			if line.Start.X > line.End.X {
				x = line.Start.X - i
			}

			y := line.Start.Y + i
			if line.Start.Y > line.End.Y {
				y = line.Start.Y - i
			}

			points = append(points, Point{
				X: x,
				Y: y,
			})
		}
	}

	return points
}

func readLines(path string) ([]Line, error) {
	fileLines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	lines := []Line{}
	for _, fileLine := range fileLines {
		// e.g. "0,9 -> 5,9" -> ["0,9", "->"", "5,9"]
		parts := strings.Fields(fileLine)
		if len(parts) != 3 {
			return nil, fmt.Errorf("expected line of format '0,9 -> 5,9' but got '%s'", fileLine)
		}

		start, err := readPoint(parts[0])
		if err != nil {
			return nil, err
		}

		end, err := readPoint(parts[2])
		if err != nil {
			return nil, err
		}

		lines = append(lines, Line{Start: start, End: end})
	}

	return lines, nil
}

func readPoint(str string) (Point, error) {
	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return Point{}, fmt.Errorf("expected point of format '0,1' but got '%s'", str)
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return Point{}, err
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return Point{}, err
	}

	return Point{X: x, Y: y}, nil
}

func filterOutDiagonals(lines []Line) []Line {
	filtered := []Line{}
	for _, line := range lines {
		if line.Start.X == line.End.X || line.Start.Y == line.End.Y {
			filtered = append(filtered, line)
		}
	}

	return filtered
}

func findOverlaps(lines []Line) []Point {
	// (x, y) => [x: [y: 1]]
	positions := map[int]map[int]int{}
	overlaps := []Point{}

	for _, line := range lines {
		for _, point := range line.Points() {
			if positions[point.X] == nil {
				positions[point.X] = map[int]int{}
			}

			positions[point.X][point.Y] += 1

			if positions[point.X][point.Y] == 2 {
				overlaps = append(overlaps, point)
			}
		}
	}

	return overlaps
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func diff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}
