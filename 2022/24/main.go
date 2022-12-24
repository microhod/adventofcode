package main

import (
	"fmt"
	"math"

	"github.com/microhod/adventofcode/internal/copy"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/geometry/plane"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Blizzard Basin", part1, part2).Run()
}

// there is the minimum time for the expedition from start to end
// we track it here so we don't have to recompute this in part 2
// (because I was lazy and didn't optimise very well, the code takes ages)
var there int

// part1 takes ~5 mins
func part1() error {
	valley, err := parse(InputFile)
	if err != nil {
		return err
	}

	start := plane.Vector{X: 1, Y: 0}
	end := plane.Vector{Y: valley.maxY, X: valley.maxX - 1}

	there = MinimumPath(valley, start, end)
	fmt.Println(there)
	return nil
}

// part2 takes ~10 mins
func part2() error {
	valley, err := parse(InputFile)
	if err != nil {
		return err
	}

	start := plane.Vector{X: 1, Y: 0}
	end := plane.Vector{Y: valley.maxY, X: valley.maxX - 1}

	// get the state of the valley after getting to the end
	for i := 0; i < there; i++ {
		valley.MoveBlizzards()
	}

	// go back and update the valley with its state after going back to the start
	back := MinimumPath(valley, end, start)
	for i := 0; i < back; i++ {
		valley.MoveBlizzards()
	}

	// go back to the end one final time
	thereAgain := MinimumPath(valley, start, end)
	
	fmt.Println(there + back + thereAgain)
	return nil
}

var RuneToDirection = map[rune]plane.Direction{
	'^': plane.North,
	'>': plane.East,
	'v': plane.South,
	'<': plane.West,
}

func parse(path string) (*Valley, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	valley := &Valley{
		blizzards: map[plane.Vector]set.Set[*Blizzard]{},
		maxX:      len(lines[0]) - 1,
		maxY:      len(lines) - 1,
	}
	for y := range lines {
		for x := range lines[y] {
			direction, exists := RuneToDirection[rune(lines[y][x])]
			if exists {
				position := plane.Vector{X: x, Y: y}
				valley.blizzards[position] = set.NewSet(&Blizzard{
					Direction: direction,
				})
			}
		}
	}
	return valley, nil
}

type Valley struct {
	blizzards  map[plane.Vector]set.Set[*Blizzard]
	maxX, maxY int
}

func (v *Valley) MoveBlizzards() {
	nextBlizzards := map[plane.Vector]set.Set[*Blizzard]{}

	for position, blizzards := range v.blizzards {
		for blizzard := range blizzards {
			next := v.getNextPosition(position, blizzard)

			if nextBlizzards[next] == nil {
				nextBlizzards[next] = set.NewSet[*Blizzard]()
			}
			nextBlizzards[next].Add(blizzard)
		}
	}
	v.blizzards = nextBlizzards
}

func (v *Valley) getNextPosition(current plane.Vector, blizzard *Blizzard) plane.Vector {
	next := current.Add(plane.DirectionToVector[blizzard.Direction])

	// wrap around if we hit a wall
	if next.X < 1 {
		next.X = v.maxX - 1
	}
	if next.X > v.maxX-1 {
		next.X = 1
	}
	if next.Y < 1 {
		next.Y = v.maxY - 1
	}
	if next.Y > v.maxY-1 {
		next.Y = 1
	}

	return next
}

func (v *Valley) IsWall(position plane.Vector) bool {
	x, y := position.X, position.Y

	// entry & exit
	if (y == 0 && x == 1) || (y == v.maxY && x == v.maxX-1) {
		return false
	}
	// boundary
	if x < 1 || x > v.maxX-1 || y < 1 || y > v.maxY-1 {
		return true
	}
	return false
}

func (v *Valley) Copy() *Valley {
	return &Valley{
		blizzards: copy.MapDeep(v.blizzards, func(position plane.Vector) set.Set[*Blizzard] {
			return set.NewSet(v.blizzards[position].ToSlice()...)
		}),
		maxX: v.maxX,
		maxY: v.maxY,
	}
}

type Blizzard struct {
	Direction plane.Direction
}

type Expedition struct {
	current plane.Vector
	time    int
}

func MinimumPath(v *Valley, start plane.Vector, end plane.Vector) int {
	// use a set for this so we don't duplicate effort if mutliple paths
	// reach the same position at the same time
	stack := set.NewSet(Expedition{current: start, time: 0})
	// valleys stores the state of the valley at each time interval (denoted by the slice index)
	valleys := []*Valley{v}

	min := math.MaxInt
	for len(stack) > 0 {
		// hacky thing to pick something from the set
		var expedition Expedition
		for e := range stack {
			expedition = e
			break
		}

		stack.Remove(expedition)

		// update the minimum if we reach the end
		if expedition.current == end {
			min = maths.Min(min, expedition.time)
			continue
		}
		// give up with the current expedition if the time is already above the minimum
		if expedition.time >= min {
			continue
		}

		// compute the valley at the next time only once 
		if len(valleys) <= expedition.time+1 {
			nextValley := valleys[expedition.time].Copy()
			nextValley.MoveBlizzards()
			valleys = append(valleys, nextValley)
		}
		valley := valleys[expedition.time+1]


		// add next options to the stack
		for _, next := range getNextOptions(valley, expedition.current) {
			stack.Add(Expedition{
				current: next, time: expedition.time + 1,
			})
		}
	}

	return min
}

func getNextOptions(valley *Valley, current plane.Vector) []plane.Vector {
	var options []plane.Vector
	for _, neighbour := range current.OrthogonalNeighbours() {
		if valley.IsWall(neighbour) {
			continue
		}
		if len(valley.blizzards[neighbour]) > 0 {
			continue
		}
		options = append(options, neighbour)
	}

	// if we can stay where we are, also add that as an option
	if len(valley.blizzards[current]) == 0 {
		options = append(options, current)
	}

	return options
}
