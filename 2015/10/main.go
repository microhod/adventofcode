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

/*
Any number e.g. 111221 can be encoded in terms of deduplicated digits and counts:
	dedupe = 121
	counts = 321

This encodes the same amount of data, and these numbers will always be the same length.
Now looksay(111221) = 312211 which is exactly count[0]dedupe[0]count[1]dedupe[1]...

So we can save time by computing the next iteration in terms of the deduped numbers and their counts.
This means we never actually have to fully compute the number (or string).
*/

func main() {
	puzzle.NewSolution("Elves Look, Elves Say", part1, part2).Run()
}

func part1() error {
	nums, counts, err := parseLooksay(InputFile)
	if err != nil {
		return err
	}

	for i := 0; i < 40; i++ {
		nums, counts = looksay(nums, counts)
	}

	fmt.Println(sum(counts))
	return nil
}

func part2() error {
	nums, counts, err := parseLooksay(InputFile)
	if err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		nums, counts = looksay(nums, counts)
	}

	fmt.Println(sum(counts))
	return nil
}

func parseLooksay(path string) ([]int, []int, error) {
	b, err := file.ReadBytes(path)
	if err != nil {
		return nil, nil, err
	}

	var nums []int
	var counts []int

	for _, char := range strings.TrimSpace(string(b)) {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, nil, err
		}

		if len(nums) > 0 && digit == nums[len(nums)-1] {
			counts[len(counts)-1]++
			continue
		}
		nums = append(nums, digit)
		counts = append(counts, 1)
	}

	return nums, counts, nil
}

func looksay(nums []int, counts []int) ([]int, []int) {
	var newNums []int
	var newCounts []int

	for i := range nums {
		if len(newNums) > 0 && newNums[len(newNums)-1] == counts[i] {
			newCounts[len(newCounts)-1]++
		} else {
			newNums = append(newNums, counts[i])
			newCounts = append(newCounts, 1)
		}

		if len(newNums) > 0 && newNums[len(newNums)-1] == nums[i] {
			newCounts[len(newCounts)-1]++
		} else {
			newNums = append(newNums, nums[i])
			newCounts = append(newCounts, 1)
		}
	}

	return newNums, newCounts
}

func sum(counts []int) int {
	var total int
	for _, c := range counts {
		total += c
	}
	return total
}
