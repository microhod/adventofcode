package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Like a GIF For Your Yard", part1, part2).Run()
}

func part1() error {
	lights, err := parse(InputFile)
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		lights = lights.Iterate()
	}

	on, _ := lights.Count()
	fmt.Printf("lights on = %d\n", on)

	return nil
}

func part2() error {
	lights, err := parse(InputFile)
	if err != nil {
		return err
	}

	lights.TurnOnCorners()
	for i := 0; i < 100; i++ {
		lights = lights.Iterate()
		lights.TurnOnCorners()
	}

	on, _ := lights.Count()
	fmt.Printf("lights on = %d\n", on)

	return nil
}

func parse(path string) (Lights, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	lights := Lights{}
	for i, line := range lines {
		if line == "" {
			continue
		}

		lights = append(lights, make([]bool, len(line)))
		for j, ch := range line {
			lights[i][j] = ch == '#'
		}
	}

	return lights, nil
}

type Lights [][]bool

func NewLights(size int) Lights {
	lights := make(Lights, size)
	for i := range lights {
		lights[i] = make([]bool, size)
	}
	return lights
}

func (l Lights) TurnOnCorners() {
	l[0][0] = true
	l[0][len(l)-1] = true
	l[len(l)-1][0] = true
	l[len(l)-1][len(l)-1] = true
}

func (l Lights) Iterate() Lights {
	next := NewLights(len(l))

	for i := range l {
		for j := range l[i] {
			next[i][j] = l[i][j]
			
			on, _ := OnOff(l.Neighbourhood(i, j))
			if next[i][j] && !(on == 2 || on == 3) {
				next[i][j] = false
			}
			if !next[i][j] && on == 3 {
				next[i][j] = true
			}
		}
	}

	return next
}

func (l Lights) Neighbourhood(row, col int) []bool {
	var neighbours []bool
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			// calculate the neighbour row & col
			nrow, ncol := row+i, col+j

			// assume off if neighbours are out of bounds
			if nrow < 0 || nrow >= len(l) || ncol < 0 || ncol >= len(l) {
				neighbours = append(neighbours, false)
				continue
			}

			neighbours = append(neighbours, l[nrow][ncol])
		}
	}
	return neighbours
}

func (l Lights) String() string {
	var lines []string
	for i := range l {
		var line string
		for j := range l[i] {
			ch := '.'
			if l[i][j] {
				ch = '#'
			}
			line += string(ch)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (l Lights) Count() (int, int) {
	var lights []bool
	for i := range l {
		for j := range l[i] {
			lights = append(lights, l[i][j]) 
		}
	}
	return OnOff(lights)
}

func OnOff(lights []bool) (int, int) {
	var on, off int
	for _, isOn := range lights {
		if isOn {
			on++
			continue
		}
		off++
	}
	return on, off
}
