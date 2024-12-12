package main

import (
	"fmt"
	"math"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Plutonian Pebbles", part1, part2).Run()
}

func part1() error {
	stones, err := file.ReadCsvInts(InputFile, " ")
	if err != nil {
		return err
	}
	var total int
	for _, s := range stones {
		total += Blink(s, 25)
	}
	fmt.Printf("# stones: %d\n", total)
	return nil
}

func part2() error {
	stones, err := file.ReadCsvInts(InputFile, " ")
	if err != nil {
		return err
	}
	
	var total int
	for _, s := range stones {
		total += Blink(s, 75)
	}
	fmt.Printf("# stones: %d\n", total)
	return nil
}

var memory = map[[2]int]int{}

func Blink(stone int, blinks int) int {
	if result, ok := memory[[2]int{stone, blinks}]; ok {
		return result
	}

	result := blink(stone, blinks)
	memory[[2]int{stone, blinks}] = result
	return result
}

func blink(stone int, blinks int) int {
	if blinks == 0 {
		return 1
	}
	if stone == 0 {
		return Blink(1, blinks-1)
	}
	if a, b, ok := splitByDigits(stone); ok {
		return Blink(a, blinks-1) + Blink(b, blinks-1)
	}
	return Blink(stone * 2024, blinks-1)
}

func splitByDigits(n int) (int, int, bool) {
	digits := numDigits(n)
	if digits%2 != 0 {
		return -1, -1, false
	}

	half := int(math.Pow10(digits / 2))
	return n / half, n % half, true
}

func numDigits(n int) int {
	return int(math.Log10(float64(n))) + 1
}
