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
	puzzle.NewSolution("Full of Hot Air", part1, part2).Run()
}

func part1() error {
	lines, err := file.ReadLines(InputFile)
	if err != nil {
		return err
	}

	var sum int
	for _, snafu := range lines {
		sum += SnafuToDecimal(snafu)
	}

	fmt.Println(DecimalToSnafu(float64(sum)))
	return nil
}

func part2() error {
	fmt.Println("no part 2 today, it's Christmas! 🎄")
	return nil
}

var SnafuRuneToDecimal = map[rune]int{
	'0': 0,
	'1': 1,
	'2': 2,
	'=': -2,
	'-': -1,
}

var DecimalToSnafuRune = map[int]rune{
	0: '0',
	1: '1',
	2: '2',
	-2: '=',
	-1: '-',
}

func SnafuToDecimal(snafu string) int {
	var decimal int

	for power := 0; power <= len(snafu)-1; power++ {
		unit := int(math.Pow(5, float64(power)))

		char := rune(snafu[len(snafu)-1-power])
		digit := SnafuRuneToDecimal[char]
		
		decimal += digit * unit
	}
	return decimal
}

func DecimalToSnafu(decimal float64) string {
	var snafu string

	start := int(Log5(decimal))
	for power := start; power >= 0; power-- {
		unit := math.Pow(5, float64(power))
		digit := math.Round(decimal / unit)

		snafu += string(DecimalToSnafuRune[int(digit)])

		decimal -= digit * unit
	}

	return snafu
}

func Log5(n float64) float64 {
	return math.Log(n) / math.Log(5.0)
}
