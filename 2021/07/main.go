package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	CrabsFile = "crabs.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("The Treachery of Whales", part1, part2).Run()
}

func part1() error {
	crabs, err := file.ReadCsvInts(CrabsFile)
	if err != nil {
		return err
	}

	pos, cost := getAlign(crabs, constantFuelCost)
	fmt.Printf("🦀 cost to align on position %d (constant fuel cost): %d\n", pos, cost)

	return nil
}

func part2() error {
	crabs, err := file.ReadCsvInts(CrabsFile)
	if err != nil {
		return err
	}

	pos, cost := getAlign(crabs, variableFuelCost)
	fmt.Printf("🦀 cost to align on position %d (variable fuel cost): %d\n", pos, cost)

	return nil
}

func getAlign(crabs []int, costFunc func([]int,int)int) (int, int) {
	// start at the average (should be 'close' to this)
	minPos := avg(crabs)
	minCost := costFunc(crabs, minPos)

	// default to direction being right
	direction := 1
	cost := costFunc(crabs, minPos+direction)
	// if cost on the left is less than the cost on the right,
	// flip direction to the left
	if leftCost := costFunc(crabs, minPos-1); leftCost < cost {
		cost = leftCost
		direction = -1
	}
	pos := minPos + direction

	for cost < minCost {
		minPos = pos
		minCost = cost

		pos += direction
		cost = costFunc(crabs, pos)

	}

	return minPos, minCost
}

func constantFuelCost(crabs []int, position int) int {
	cost := 0
	for _, crab := range crabs {
		cost += diff(crab, position)
	}
	return cost
}

func variableFuelCost(crabs []int, position int) int {
	cost := 0
	for _, crab := range crabs {
		distance := diff(crab, position)
		// sum of integers between 1 and distance
		cost += (distance * (distance + 1)) / 2
	}
	return cost
}

func diff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func avg(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}

	return sum / len(nums)
}
