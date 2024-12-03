package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Mull It Over", part1, part2).Run()
}

func part1() error {
	program, err := parse(InputFile, false)
	if err != nil {
		return err
	}

	fmt.Printf("multiplication result: %d\n", program.Run())
	return nil
}

func part2() error {
	program, err := parse(InputFile, true)
	if err != nil {
		return err
	}

	fmt.Printf("multiplication result: %d\n", program.Run())
	return nil
}

func parse(path string, conditionals bool) (*Program, error) {
	memory, err := file.Read(path)
	if err != nil {
		return nil, err
	}

	pattern := `mul\(\d+,\d+\)`
	if conditionals {
		pattern = `(mul\(\d+,\d+\))|(do\(\))|(don't\(\))`
	}

	var ops []Operation
	for _, op := range regexp.MustCompile(pattern).FindAllString(memory, -1) {
		if strings.HasPrefix(op, "don't") {
			ops = append(ops, func(_ bool) (int, bool) {
				return 0, false
			})
			continue
		}
		if strings.HasPrefix(op, "do") {
			ops = append(ops, func(_ bool) (int, bool) {
				return 0, true
			})
			continue
		}

		op = op[4 : len(op)-1]
		parts := strings.Split(op, ",")

		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		b, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		ops = append(ops, func(enabled bool) (int, bool) {
			if !enabled {
				return 0, enabled
			}
			return a * b, enabled
		})
	}
	return &Program{
		enabled: true,
		ops:     ops,
	}, nil
}

type Program struct {
	ops     []Operation
	enabled bool
	total   int
}

func (p Program) Run() int {
	for _, op := range p.ops {
		var result int
		result, p.enabled = op(p.enabled)
		p.total += result
	}
	return p.total
}

type Operation func(bool) (int, bool)
