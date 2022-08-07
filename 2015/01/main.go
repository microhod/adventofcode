package main

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
	puzzle.NewSolution("Not Quite Lisp", part1, part2).Run()
}

func part1() error {
	instructions, err := file.ReadBytes(InputFile)
	if err != nil {
		return err
	}

	floor := 0
	for _, instruction := range instructions {
		// '(' == 40 && ')' == 41
		// so '(' ==> 1 && ')' ==> -1
		floor += int(2 * (40.5 - float32(instruction)))
	}

	fmt.Printf("Santa is on floor: %d\n", floor)
	return nil
}

func part2() error {
	instructions, err := file.ReadBytes(InputFile)
	if err != nil {
		return err
	}

	floor := 0
	for position, instruction := range instructions {
		// '(' == 40 && ')' == 41
		// so '(' ==> 1 && ')' ==> -1
		floor += int(2 * (40.5 - float32(instruction)))

		// if in basement
		if floor == -1 {
			fmt.Printf("Got to basement at position: %d\n", position+1)
			return nil
		}
	}

	fmt.Println("Never got to the basement...oh dear!")
	return nil
}
