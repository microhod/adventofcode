package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	FishFile = "fish.txt"
	TestFile = "test.txt"
)

func main() {
	puzzle.NewSolution("Lanternfish", part1, part2).Run()
}

func part1() error {
	fish, err := readLanternFish(FishFile)
	if err != nil {
		return err
	}

	days := 80
	size := modelPopulationGrowth(fish, days)

	fmt.Printf("population size after %d days = %d\n", days, size)

	return nil
}

func part2() error {
	fish, err := readLanternFish(FishFile)
	if err != nil {
		return err
	}

	days := 256
	size := modelPopulationGrowth(fish, days)

	fmt.Printf("population size after %d days = %d\n", days, size)

	return nil
}

func modelPopulationGrowth(fish []*LanternFish, days int) int {
	sizes := getPopulationSizes(days)

	total := 0
	for _, f := range fish {
		// use the data based on sizes of a fish with initial state 0 to calculate
		// populations of each fish's children
		// (the populations are independent so we can just sum them to get the total)
		total += sizes[days-1-f.Timer]
	}

	return total
}

/*
This function gets the population size after the day provided
for the initial condition of one fish with initial state of 0.

The sequence is modelled as follows (where the 'days' parameter corresponds to n
and P_n is the population after n days):

  P_n = P_(n-7) + P_(n-9) (for n >= 7)

  with initial conditions P_(-1) = P_(0) = 1

(We calculate for a fish of initial state 0 as any other state is then just
and earlier element in the sequence, so we only have to compute this list once!)

Small explanation:

Taking one fish with an initial state of 6 gives:

6 =7=days=> 8,6 =7=days=> 1,6,8 =2=days=> 6,4,6,8

so taking P_n = population if the starting point was 6,4,6,8, then
6,8 is P_(n-7) and 4,6 is just 6,8 but two days later, hence P_(n-9).
*/
func getPopulationSizes(days int) []int {
	// initial values P_(-1) = P_(0) = 1
	sizes := []int{1, 1}
	for day := 1; day <= days; day++ {
		var size int

		switch {
		case day > 7:
			size = sizes[day+1-7] + sizes[day+1-9]
		default:
			size = 2
		}

		sizes = append(sizes, size)
	}

	// remove the initial conditions of P_(-1) and P_(0)
	return sizes[2:]
}

type LanternFish struct {
	Timer int
}

func readLanternFish(path string) ([]*LanternFish, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	if len(lines) < 1 {
		return []*LanternFish{}, nil
	}

	timers, err := parseCsvInts(lines[0])
	if err != nil {
		return nil, err
	}

	fish := []*LanternFish{}
	for _, timer := range timers {
		fish = append(fish, &LanternFish{Timer: timer})
	}

	return fish, nil
}

func parseCsvInts(str string) ([]int, error) {
	nums := []int{}
	for _, s := range strings.Split(str, ",") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}

		nums = append(nums, n)
	}

	return nums, nil
}
