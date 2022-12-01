package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Calorie Counting", part1, part2).Run()
}

func part1() error {
	elves, err := parse(InputFile)
	if err != nil {
		return err
	}

	var maxCalories int
	for _, elf := range elves {
		if total := elf.TotalCalories(); total > maxCalories {
			maxCalories = total
		}
	}

	fmt.Printf("Max calories carried by an elf: %d\n", maxCalories)
	return nil
}

func part2() error {
	elves, err := parse(InputFile)
	if err != nil {
		return err
	}

	totals := make([]int, len(elves))
	for i := range elves {
		totals[i] = elves[i].TotalCalories()
	}

	sort.Ints(totals)
	top3 := totals[len(totals)-1] + totals[len(totals)-2] + totals[len(totals)-3]

	fmt.Printf("Total calories carried by the top 3 elves: %d\n", top3)

	return nil
}

func parse(path string) ([]Elf, error) {
	data, err := file.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	var elves []Elf

	elfCalories := bytes.Split(data, []byte("\n\n"))
	for _, lines := range elfCalories {
		var e Elf

		for _, line := range bytes.Split(lines, []byte("\n")) {
			if len(line) == 0 {
				continue
			}

			calories, err := strconv.Atoi(string(line))
			if err != nil {
				return nil, err
			}

			e.Calories = append(e.Calories, calories)
		}

		elves = append(elves, e)
	}

	return elves, nil
}

type Elf struct {
	Calories []int
}

func (e Elf) TotalCalories() int {
	var total int
	for _, calories := range e.Calories {
		total += calories
	}
	return total
}
