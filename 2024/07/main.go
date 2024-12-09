package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/copy"
	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Bridge Repair", part1, part2).Run()
}

func part1() error {
	equations, err := parse(InputFile)
	if err != nil {
		return err
	}

	var total int64
	for _, eq := range equations {
		if eq.Possible(1) {
			total += eq.target
		}
	}

	fmt.Printf("total test values from possible equations: %d\n", total)
	return nil
}

func part2() error {
	equations, err := parse(InputFile)
	if err != nil {
		return err
	}

	var total int64
	for _, eq := range equations {
		if eq.Possible(2) {
			total += eq.target
		}
	}

	fmt.Printf("total test values from possible equations: %d\n", total)
	return nil
}

func parse(path string) ([]Equation, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var equations []Equation
	for _, line := range lines {
		parts := strings.Split(line, ": ")
		target, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}

		nums, err := csv.ParseInt64s(parts[1], " ")
		if err != nil {
			return nil, err
		}

		equations = append(equations, Equation{
			nums:   nums,
			target: target,
		})
	}
	return equations, nil
}

type Equation struct {
	nums   []int64
	target int64
}

func (e Equation) Possible(maxOp int) bool {
	return e.possible(
		make([]uint8, len(e.nums)-1),
		0,
		maxOp,
	)
}

func (e Equation) possible(ops []uint8, next int, maxOp int) bool {
	if next >= len(ops) {
		return false
	}

	for op := 0; op <= maxOp; op++ {
		cops := copy.Slice(ops)
		cops[next] = uint8(op)
		if e.evaluate(cops) == e.target {
			return true
		}
		if e.possible(cops, next+1, maxOp) {
			return true
		}
	}

	return false
}

func (e Equation) evaluate(ops []uint8) int64 {
	if len(ops) != len(e.nums)-1 {
		panic(fmt.Sprintf("incorrect num operations for %d nums: %d", len(e.nums), len(ops)))
	}

	result := e.nums[0]
	for i, o := range ops {
		switch o {
		case 0:
			result += e.nums[i+1]
		case 1:
			result *= e.nums[i+1]
		case 2:
			var err error
			result, err = strconv.ParseInt(fmt.Sprint(result)+fmt.Sprint(e.nums[i+1]), 10, 64)
			if err != nil {
				panic(err)
			}
		default:
			panic(fmt.Sprintf("unknown operation: %d", o))
		}
	}

	return result
}
