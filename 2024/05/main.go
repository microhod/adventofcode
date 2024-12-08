package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
	islices "github.com/microhod/adventofcode/internal/slices" 
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Print Queue", part1, part2).Run()
}

func part1() error {
	rules, updates, err := parse(InputFile)
	if err != nil {
		return err
	}

	var middles int
	for _, u := range updates {
		if rules.CorrectlyOrdered(u) {
			middles += u.Middle()
		}
	}
	fmt.Printf("middle page nums from correctly-ordered updates: %d\n", middles)
	return nil
}

func part2() error {
	rules, updates, err := parse(InputFile)
	if err != nil {
		return err
	}

	var middles int
	for _, u := range updates {
		if rules.CorrectlyOrdered(u) {
			continue
		}
		u = rules.FixOrdering(u)
		middles += u.Middle()
	}
	fmt.Printf("middle page nums from fixed incorrectly-ordered updates: %d\n", middles)
	return nil
}

func parse(path string) (Rules, []Update, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, nil, err
	}

	parts := strings.Split(input, "\n\n")

	rules := make(Rules)
	for _, line := range strings.Split(parts[0], "\n") {
		if line == "" {
			continue
		}
		nums, err := csv.ParseInts(line, "|")
		if err != nil {
			return nil, nil, err
		}

		if rules[nums[1]] == nil {
			rules[nums[1]] = set.NewSet[int]()
		}
		rules[nums[1]].Add(nums[0])
	}

	var updates []Update
	for _, line := range strings.Split(parts[1], "\n") {
		if line == "" {
			continue
		}
		nums, err := csv.ParseInts(line)
		if err != nil {
			return nil, nil, err
		}
		updates = append(updates, Update(nums))
	}
	return rules, updates, nil
}

type Rules map[int]set.Set[int]

func (r Rules) CorrectlyOrdered(u Update) bool {
	for i := range u {
		before := r[u[i]]
		if before == nil {
			before = set.NewSet[int]()
		}
		before = set.Intersect(before, set.NewSet(u...))

		for b := range before {
			if slices.Index(u, b) > i {
				return false
			}
		}
	}
	return true
}

func (r Rules) FixOrdering(u Update) Update {
	for i := len(u)-1; i >= 0; i-- {
		before := r[u[i]]
		if before == nil {
			before = set.NewSet[int]()
		}
		before = set.Intersect(before, set.NewSet(u...))

		for b := range before {
			if slices.Index(u, b) < i {
				before.Remove(b)
			}
		}
		if len(before) == 0 {
			continue
		}

		// move pages which need to be before i
		move := before.ToSlice()
		slices.SortFunc(move, func(a, b int) int {
			return slices.Index(u, a) - slices.Index(u, b)
		})
		u = slices.Insert(u, i, move...)
		// remove duplicates
		for j := len(u)-1; j > i + len(move); j-- {
			if before.Contains(u[j]) {
				u = append(u[:j], u[j+1:]...)
			}
		}
	}
	return u
}

type Update []int

func (u Update) Middle() int {
	if len(u) % 2 != 1 {
		panic(fmt.Sprintf("update has an even number of pages: %d", len(u)))
	}
	return u[(len(u)-1)/2]
}

func (u Update) ShiftDown(i, j int) Update {
	if j >= i {
		panic(fmt.Sprintf("shiftdown: %d >= %d", j, i))
	}
	return islices.Appends(u[:j], []int{u[i]}, u[j:i], u[i+1:])
}
