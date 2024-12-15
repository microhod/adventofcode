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
	puzzle.NewSolution("Claw Contraption", part1, part2).Run()
}

func part1() error {
	machines, err := parse(InputFile)
	if err != nil {
		return err
	}

	var total int64
	for i := range machines {
		solution := machines[i].Solve()
		total += 3*solution[0] + solution[1]
	}
	fmt.Printf("fewest tokens to win all prizes: %d\n", total)
	return nil
}

func part2() error {
	machines, err := parse(InputFile)
	if err != nil {
		return err
	}

	var total int64
	for i := range machines {
		machines[i].Prize.X += 10000000000000
		machines[i].Prize.Y += 10000000000000

		solution := machines[i].Solve()
		total += 3*solution[0] + solution[1]
	}
	fmt.Printf("fewest tokens to win all prizes: %d\n", total)
	return nil
}

func parse(path string) ([]ClawMachine, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, err
	}

	claws := strings.Split(input, "\n\n")

	var machines []ClawMachine
	for _, claw := range claws {
		nums := regexp.MustCompile(`\d+`).FindAllString(claw, 6)
		parse := func(s string) int64 {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic(err)
			}
			return n
		}

		machines = append(machines, ClawMachine{
			A:     Vector64{X: parse(nums[0]), Y: parse(nums[1])},
			B:     Vector64{X: parse(nums[2]), Y: parse(nums[3])},
			Prize: Vector64{X: parse(nums[4]), Y: parse(nums[5])},
		})
	}
	return machines, nil
}

type ClawMachine struct {
	A, B, Prize Vector64
}

func (m ClawMachine) Solve() [2]int64 {
	// https://en.wikipedia.org/wiki/Cramer%27s_rule
	a := (m.Prize.X * m.B.Y - m.Prize.Y * m.B.X) /
	(m.A.X * m.B.Y - m.A.Y * m.B.X)
	b := (m.Prize.Y * m.A.X - m.Prize.X * m.A.Y) /
	(m.B.Y * m.A.X - m.B.X * m.A.Y)

	if m.run([2]int64{a, b}) != m.Prize {
		return [2]int64{}
	}

	return [2]int64{a, b}
}

func (m ClawMachine) run(ab [2]int64) Vector64 {
	return Vector64{
		X: ab[0]*m.A.X + ab[1]*m.B.X,
		Y: ab[0]*m.A.Y + ab[1]*m.B.Y,
	}
}

type Vector64 struct {
	X, Y int64
}
