package main

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Science for Hungry People", part1, part2).Run()
}

func part1() error {
	var max int
	for a := 0; a <= 100; a++ {
		for b := 0; b <= 100-a; b++ {
			for c := 0; c <= 100-a-b; c++ {
				recipe := []int{a, b, c, 100 - a - b - c}

				score := product(multiply(properties[:4], recipe))
				if score > max {
					max = score
				}
			}
		}
	}
	fmt.Println(max)

	return nil
}

func part2() error {
	var max int
	for a := 0; a <= 100; a++ {
		for b := 0; b <= 100-a; b++ {
			for c := 0; c <= 100-a-b; c++ {
				recipe := []int{a, b, c, 100 - a - b - c}

				calories := product(multiply(properties[4:], recipe))
				if calories != 500 {
					continue
				}

				score := product(multiply(properties[:4], recipe))
				if score > max {
					max = score
				}
			}
		}
	}
	fmt.Println(max)
	return nil
}

var properties = [][]int{
	{2, 0, 0, 0},   // capacity
	{0, 5, 0, -1},  // durability
	{-2, -3, 5, 0}, // flavour
	{0, 0, -1, 5},  // texture
	{3, 3, 8, 8},   // calories
}

func multiply(matrix [][]int, vector []int) []int {
	result := make([]int, len(matrix))
	for row := range matrix {
		for col := range matrix[row] {
			result[row] += matrix[row][col] * vector[col]
		}
	}
	return result
}

func product(vector []int) int {
	result := 1
	for _, element := range vector {
		// any negative elements are moved back to zero
		if element < 0 {
			element = 0
		}
		result *= element
	}
	return result
}
