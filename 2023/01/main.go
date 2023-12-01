package main

import (
	"fmt"
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
	puzzle.NewSolution("Trebuchet?!", part1, part2).Run()
}

func part1() error {
	calibrations, err := parseInput(InputFile, false)
	if err != nil {
		return err
	}

	var sum int
	for _, value := range calibrations.Values() {
		sum += value
	}
	fmt.Printf("Sum of all calibration values: %d\n", sum)
	return nil
}

func part2() error {
	calibrations, err := parseInput(InputFile, true)
	if err != nil {
		return err
	}

	var sum int
	for _, value := range calibrations.Values() {
		sum += value
	}
	fmt.Printf("Sum of all calibration values: %d\n", sum)
	return nil
}

func parseInput(path string, includeTextNumbers bool) (Calibrations, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	calibrations := make(Calibrations, 0)
	for _, line := range lines {
		if includeTextNumbers {
			line = replaceTextNumbers(line)
		}

		var calibration []int
		for _, ch := range line {
			number, err := strconv.Atoi(string(ch))
			if err != nil {
				continue
			}
			calibration = append(calibration, number)
		}
		calibrations = append(calibrations, calibration)
	}

	return calibrations, nil
}

func replaceTextNumbers(txt string) string {
	replacements := map[string]string{
		"one":   "o1ne",
		"two":   "t2wo",
		"three": "t3hree",
		"four":  "f4our",
		"five":  "f5ive",
		"six":   "s6ix",
		"seven": "s7even",
		"eight": "e8ight",
		"nine":  "n9ine",
	}
	for t, n := range replacements {
		txt = strings.ReplaceAll(txt, t, t+n)
	}
	return txt
}

type Calibrations [][]int

func (c Calibrations) Values() []int {
	var values []int
	for _, calibration := range c {
		values = append(values, 10*calibration[0]+calibration[len(calibration)-1])
	}
	return values
}
