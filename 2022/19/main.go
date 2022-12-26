package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/microhod/adventofcode/internal/copy"
	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Not Enough Minerals", part1, part2).Run()
}

// part1 takes ~40 secs
func part1() error {
	blueprints, err := parse(InputFile)
	if err != nil {
		return err
	}

	var sum int
	for _, blueprint := range blueprints {
		max := blueprint.MaxGeodes(24)
		sum += blueprint.number * max
	}
	fmt.Println(sum)
	return nil
}

// part2 takes ~7 mins
func part2() error {
	blueprints, err := parse(InputFile)
	if err != nil {
		return err
	}

	// elephants ate all blueprints apart from the first three
	if len(blueprints) > 3 {
		blueprints = blueprints[:3]
	}

	product := 1
	for _, blueprint := range blueprints {
		max := blueprint.MaxGeodes(32)
		product *= max
	}

	fmt.Println(product)
	return nil
}

func parse(path string) ([]*Blueprint, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, err
	}

	input = strings.ReplaceAll(input, "Blueprint ", "")
	input = strings.ReplaceAll(input, ": Each ore robot costs ", ",")
	input = strings.ReplaceAll(input, " ore. Each clay robot costs ", ",")
	input = strings.ReplaceAll(input, " ore. Each obsidian robot costs ", ",")
	input = strings.ReplaceAll(input, " ore and ", ",")
	input = strings.ReplaceAll(input, " clay. Each geode robot costs ", ",")
	input = strings.ReplaceAll(input, " obsidian.", "")

	var blueprints []*Blueprint
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		nums, err := csv.ParseInts(line)
		if err != nil {
			return nil, err
		}

		blueprints = append(blueprints, NewBlueprint(
			nums[0],
			map[Resource]map[Resource]int{
				Ore:      {Ore: nums[1]},
				Clay:     {Ore: nums[2]},
				Obsidian: {Ore: nums[3], Clay: nums[4]},
				Geode:    {Ore: nums[5], Obsidian: nums[6]},
			},
		))
	}
	return blueprints, nil
}

type Resource string

const (
	None     Resource = ""
	Ore      Resource = "ore"
	Clay     Resource = "clay"
	Obsidian Resource = "obsidian"
	Geode    Resource = "geode"
)

type Blueprint struct {
	number int
	costs  map[Resource]map[Resource]int
	limits map[Resource]int
}

func NewBlueprint(number int, costs map[Resource]map[Resource]int) *Blueprint {
	limits := map[Resource]int{
		Ore:      maths.Max(costs[Clay][Ore], costs[Obsidian][Ore], costs[Geode][Ore]),
		Clay:     costs[Obsidian][Clay],
		Obsidian: costs[Geode][Obsidian],
		Geode:    math.MaxInt,
	}
	return &Blueprint{
		number: number,
		costs:  costs,
		limits: limits,
	}
}

func (bp *Blueprint) MaxGeodes(timeLimit int) int {
	stack := []*Mine{{
		robots:    Robots{Ore: 1},
		inventory: Inventory{},
		time:      1,
	}}
	seen := set.NewSet[string]()

	var maxGeodes int
	for len(stack) > 0 {
		mine := stack[0]
		stack = stack[1:]

		// skip situations already seen
		hash := mine.Hash()
		if seen.Contains(hash) {
			continue
		}
		seen.Add(hash)

		// check next options
		nextRobots := bp.getOptions(mine, timeLimit)

		// run robots
		mine.Run()

		// stop if at time limit
		if mine.time == timeLimit {
			maxGeodes = maths.Max(maxGeodes, mine.inventory[Geode])
			continue
		}

		// add next options
		mine.time += 1
		for _, robot := range nextRobots {
			next := mine.Build(robot, bp.costs[robot])
			// skip any options we've already seen
			if seen.Contains(next.Hash()) {
				continue
			}
			stack = append(stack, next)
		}
	}

	return maxGeodes
}

func (bp *Blueprint) getOptions(mine *Mine, timeLimit int) []Resource {
	afford := bp.getCanAfford(mine.inventory)

	// remove anything over the limit
	for _, robot := range afford.ToSlice() {
		if robot != Geode && mine.robots[robot] >= bp.limits[robot] {
			afford.Remove(robot)
		}
	}

	// don't build anything in the last minute
	if mine.time == timeLimit {
		return []Resource{None}
	}
	// only build geode in the second or third to last minute
	// as we can't benefit from anything else at this stage
	if mine.time >= timeLimit-2 {
		if afford.Contains(Geode) {
			return []Resource{Geode}
		}
		return []Resource{None}
	}
	// don't build Clay after timeLimit-5 as timeLimit-3 is the last
	// time Obsidian could be built
	if mine.time >= timeLimit-4 {
		afford.Remove(Clay)
	}

	// always build Geode if we can
	if afford.Contains(Geode) {
		return []Resource{Geode}
	}

	// initialise to just 'do nothing'
	options := []Resource{None}
	// try everything we can afford
	for robot := range afford {
		options = append(options, robot)
	}

	return options
}

func (bp *Blueprint) getCanAfford(inventory Inventory) set.Set[Resource] {
	afford := set.NewSet[Resource]()
	for robot := range bp.costs {
		if bp.canAfford(robot, inventory) {
			afford.Add(robot)
		}
	}
	return afford
}

func (bp *Blueprint) canAfford(robot Resource, inventory Inventory) bool {
	for required, number := range bp.costs[robot] {
		if inventory[required] < number {
			return false
		}
	}
	return true
}

type Mine struct {
	robots    Robots
	inventory Inventory
	time      int
}

func (mine *Mine) Run() {
	for resource, count := range mine.robots {
		mine.inventory[resource] += count
	}
}

func (mine *Mine) Build(robot Resource, costs map[Resource]int) *Mine {
	if robot == "" {
		return mine
	}

	copy := &Mine{
		robots:    copy.Map(mine.robots),
		inventory: copy.Map(mine.inventory),
		time:      mine.time,
	}
	copy.robots[robot] += 1

	for required, cost := range costs {
		copy.inventory[required] -= cost
	}
	return copy
}

func (mine *Mine) Hash() string {
	hash := make([]byte, 9)
	for idx, resource := range []Resource{Ore, Clay, Obsidian, Geode} {
		hash[2*idx] = byte(mine.robots[resource])
		hash[2*idx+1] = byte(mine.inventory[resource])
	}
	hash[8] = byte(mine.time)
	return string(hash)
}

type Inventory map[Resource]int

type Robots map[Resource]int
