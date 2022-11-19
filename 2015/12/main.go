package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("JSAbacusFramework.io", part1, part2).Run()
}

func part1() error {
	in, err := file.ReadBytes(InputFile)
	if err != nil {
		return err
	}
	input := string(in)

	nums := regexp.MustCompile(`-?\d+`).FindAllStringIndex(input, -1)

	var total int
	for _, num := range nums {
		n, err := strconv.Atoi(input[num[0]:num[1]])
		if err != nil {
			return err
		}
		total += n
	}

	fmt.Println(total)
	return nil
}

func part2() error {
	in, err := file.ReadBytes(InputFile)
	if err != nil {
		return err
	}
	input := string(in)

	reds := regexp.MustCompile(`"red"`).FindAllStringIndex(input, -1)

	var objectReds [][]int
	for _, red := range reds {
		if isInArray(input, red[0]) {
			continue
		}

		objectReds = append(objectReds, []int{
			findObjectStart(input, red[0]),
			findObjectEnd(input, red[1])+1,
		})
	}

	nums := regexp.MustCompile(`-?\d+`).FindAllStringIndex(input, -1)

	var total int
	for _, num := range nums {
		if inRange(objectReds, num) {
			continue
		}

		n, err := strconv.Atoi(input[num[0]:num[1]])
		if err != nil {
			return err
		}
		total += n
	}

	fmt.Println(total)

	return nil
}

func isInArray(str string, index int) bool {
	return findArrayStart(str, index) > findObjectStart(str, index)
}

func findArrayStart(str string, index int) int {
	brakets := -1
	for i := index; i >= 0; i-- {
		switch str[i] {
		case '[':
			brakets += 1
		case ']':
			brakets -= 1
		}
		if brakets == 0 {
			return i
		}
	}
	return -1
}

func findObjectStart(str string, index int) int {
	brakets := -1
	for i := index; i >= 0; i-- {
		switch str[i] {
		case '{':
			brakets += 1
		case '}':
			brakets -= 1
		}
		if brakets == 0 {
			return i
		}
	}
	return -1
}

func findObjectEnd(str string, index int) int {
	brakets := 1
	for i := index; i < len(str); i++ {
		switch str[i] {
		case '{':
			brakets += 1
		case '}':
			brakets -= 1
		}
		if brakets == 0 {
			return i
		}
	}
	return -1
}

func inRange(ranges [][]int, rng []int) bool {
	for _, r := range ranges {
		if rng[0] >= r[0] && rng[1] <= r[1] {
			return true
		}
	}
	return false
}
