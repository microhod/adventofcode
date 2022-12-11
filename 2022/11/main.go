package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Monkey in the Middle", part1, part2).Run()
}

func part1() error {
	monkeys, err := parse(InputFile)
	if err != nil {
		return err
	}

	for _, monkey := range monkeys {
		monkey.InspectRelief = 3
	}

	for round := 0; round < 20; round++ {
		for _, monkey := range monkeys {
			monkey.TakeTurn()
		}
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].InspectCounter > monkeys[j].InspectCounter
	})

	fmt.Println(monkeys[0].InspectCounter * monkeys[1].InspectCounter)
	return nil
}

func part2() error {
	monkeys, err := parse(InputFile)
	if err != nil {
		return err
	}

	for round := 0; round < 10000; round++ {
		for _, monkey := range monkeys {
			monkey.TakeTurn()
		}
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].InspectCounter > monkeys[j].InspectCounter
	})

	fmt.Println(monkeys[0].InspectCounter * monkeys[1].InspectCounter)
	return nil
}

func parse(path string) ([]*Monkey, error) {
	data, err := file.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	monkeyInputs := strings.Split(string(data), "\n\n")
	var monkeys []*Monkey
	for i := 0; i < len(monkeyInputs); i++ {
		monkeys = append(monkeys, &Monkey{})
	}

	inspectModulo := 1
	for i, input := range monkeyInputs {
		err := parseMonkey(input, monkeys, i)
		if err != nil {
			return nil, err
		}

		inspectModulo *= monkeys[i].TargetDivisor
	}

	for i := range monkeys {
		monkeys[i].InspectModulo = inspectModulo
	}

	return monkeys, nil
}

func parseMonkey(input string, monkeys []*Monkey, index int) error {
	lines := strings.Split(input, "\n")

	items, err := csv.ParseInts(lines[1][18:], ", ")
	if err != nil {
		return err
	}

	divisor, err := strconv.Atoi(lines[3][21:])
	if err != nil {
		return err
	}
	operation, err := parseOperation(lines[2][19:])
	if err != nil {
		return err
	}

	targets := map[bool]*Monkey{}

	trueIndex, err := strconv.Atoi(lines[4][29:])
	if err != nil {
		return err
	}
	falseIndex, err := strconv.Atoi(lines[5][30:])
	if err != nil {
		return err
	}

	targets[true] = monkeys[trueIndex]
	targets[false] = monkeys[falseIndex]

	monkeys[index].Items = items
	monkeys[index].InspectOperation = operation
	monkeys[index].TargetDivisor = divisor
	monkeys[index].TargetMonkeys = targets
	return nil
}

func parseOperation(input string) (func(int) int, error) {
	if input == "old * old" {
		return func(old int) int {
			return old * old
		}, nil
	}

	fields := strings.Fields(input)

	operation := fields[1]
	number, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, err
	}

	switch operation {
	case "+":
		return func(old int) int {
			return old + int(number)
		}, nil
	case "*":
		return func(old int) int {
			return old * int(number)
		}, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}

type Monkey struct {
	Items []int

	InspectOperation func(int) int
	InspectRelief    int
	InspectModulo    int
	InspectCounter   int

	TargetDivisor int
	TargetMonkeys map[bool]*Monkey
}

func (m *Monkey) TakeTurn() {
	for _, item := range m.Items {
		m.Throw(m.Inspect(item))
	}
	m.Items = make([]int, 0)
}

func (m *Monkey) Inspect(item int) int {
	m.InspectCounter += 1

	item = m.InspectOperation(item)
	if m.InspectRelief > 0 {
		item /= m.InspectRelief
	}

	// ensure the items stay small by taking the remainder on division by the
	// product of all the divisors of the monkeys 
	// (I think this doesn't affect the tests because all the divisors are prime??)
	item = item % m.InspectModulo

	return item
}

func (m *Monkey) Throw(item int) {
	divisible := (item % m.TargetDivisor) == 0
	target := m.TargetMonkeys[divisible]

	target.Recieve(item)
}

func (m *Monkey) Recieve(items ...int) {
	m.Items = append(m.Items, items...)
}
