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
	puzzle.NewSolution("Knights of the Dinner Table", part1, part2).Run()
}

func part1() error {
	guests, err := parse(InputFile)
	if err != nil {
		return err
	}

	var max int
	for _, permutation := range getPermutations(len(guests.Names)) {
		total := guests.TotalHappiness(permutation)
		if total > max {
			max = total
		}
	}
	fmt.Println(max)

	return nil
}

func part2() error {
	guests, err := parse(InputFile)
	if err != nil {
		return err
	}

	// add me
	guests.Happiness["Me"] = map[string]int{}
	for _, name := range guests.Names {
		guests.Happiness["Me"][name] = 0
		guests.Happiness[name]["Me"] = 0
	}
	guests.Names = append(guests.Names, "Me")

	var max int
	for _, permutation := range getPermutations(len(guests.Names)) {
		total := guests.TotalHappiness(permutation)
		if total > max {
			max = total
		}
	}
	fmt.Println(max)

	return nil
}

func parse(path string) (Guests, error) {
	g := Guests{
		Happiness: map[string]map[string]int{},
	}

	lines, err := file.ReadLines(path)
	if err != nil {
		return g, err
	}

	seen := map[string]bool{}
	for _, line := range lines {
		fields := strings.Fields(line)

		from, sign, h, to := fields[0], fields[2], fields[3], fields[10]

		// trim full stop
		to = to[:len(to)-1]
		if !seen[from] {
			g.Names = append(g.Names, from)
			seen[from] = true
		}
		if !seen[to] {
			g.Names = append(g.Names, to)
			seen[to] = true
		}

		happiness, err := strconv.Atoi(h)
		if err != nil {
			return g, err
		}
		if sign == "lose" {
			happiness *= -1
		}
		if g.Happiness[from] == nil {
			g.Happiness[from] = map[string]int{}
		}
		g.Happiness[from][to] = happiness
	}

	return g, nil
}

type Guests struct {
	Names     []string
	Happiness map[string]map[string]int
}

func (g Guests) TotalHappiness(order []int) int {
	var total int
	for i := range order {
		guest := g.Names[order[i]]

		// wrap back around to start if at the end
		if i == len(order)-1 {
			i = -1
		}
		next := g.Names[order[i+1]]

		total += g.Happiness[guest][next]
		total += g.Happiness[next][guest]
	}

	return total
}

func getPermutations(n int) [][]int {
	var permutations [][]int

	partials := [][]int{}
	for i := 0; i < n; i++ {
		partials = append(partials, []int{i})
	}

	for len(partials) > 0 {
		partial := partials[0]
		partials = partials[1:]

		if len(partial) == n {
			permutations = append(permutations, partial)
		}

		for i := 0; i < n; i++ {
			if contains(partial, i) {
				continue
			}
			next := append([]int{}, partial...)
			next = append(next, i)
			partials = append(partials, next)
		}
	}
	return permutations
}

func contains(nums []int, num int) bool {
	for _, n := range nums {
		if n == num {
			return true
		}
	}
	return false
}
