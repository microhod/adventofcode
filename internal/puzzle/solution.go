// Package puzzle contains general formatting etc for any of the puzzles
package puzzle

import (
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/mgutz/ansi"
	"github.com/microhod/adventofcode/internal/christmas"
)

var (
	BoldRed   = ansi.ColorFunc("red+bh")
	BoldGreen = ansi.ColorFunc("green+bh")

	templateSolutionFile = 
`package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("{{.Name}}", part1, part2).Run()
}

func part1() error {
	return nil
}

func part2() error {
	return nil
}
`
)

func InitialSolutionFile(puzzle *Puzzle) (string, error) {
	tmpl, err := template.New("").Parse(templateSolutionFile)
	if err != nil {
		return "", err
	}
	
	builder := new(strings.Builder)
	err = tmpl.Execute(builder, puzzle)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

type Solution struct {
	Name  string
	Parts []func() error
}

func NewSolution(name string, parts ...func() error) *Solution {
	return &Solution{Name: name, Parts: parts}
}

func (s *Solution) Run() {
	// disable timstamps for logging
	log.SetFlags(0)

	// print christmas tree
	fmt.Println()
	fmt.Println(christmas.Tree())

	// print puzzle name
	fmt.Println(christmas.Lights())
	fmt.Println()
	fmt.Println(BoldGreen(fmt.Sprintf("Puzzle: %s", s.Name)))
	fmt.Println()
	fmt.Println(christmas.Lights())

	fmt.Println()

	for i, part := range s.Parts {
		// Print part number
		fmt.Println(BoldRed(fmt.Sprintf("Part %d", i+1)))
		fmt.Println()

		// run part
		start := time.Now()
		err := part()
		elapsed := time.Since(start)

		if err != nil {
			log.Fatalf("oh no! Christmas is cancelled 😱 => %s", err.Error())
		}

		fmt.Println()
		log.Printf("⏰ %s", elapsed)
		fmt.Println()
	}
}
