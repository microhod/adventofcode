package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Cube Conundrum", part1, part2).Run()
}

func part1() error {
	games, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	bag := map[Colour]int{
		Red: 12,
		Green: 13,
		Blue: 14,
	}

	var possible int
	for _, game := range games {
		if game.Possible(bag) {
			possible += game.id 
		}
	}

	fmt.Printf("sum of possible game IDs: %d\n", possible)
	return nil
}

func part2() error {
	games, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	var powers int
	for _, game := range games {
		powers += game.MinPossible().Power()
	}

	fmt.Printf("sum of minimum cube-set powers: %d\n", powers)
	return nil
}

func parseInput(path string) ([]Game, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var games []Game
	for _, line := range lines {
		game, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

func parseLine(line string) (Game, error) {
	parts := strings.Split(line, ": ")
	id, err := strconv.Atoi(strings.Fields(parts[0])[1])
	if err != nil {
		return Game{}, err
	}

	game := Game{id: id}
	rounds := strings.Split(parts[1], ";")
	for _, r := range rounds {
		round := make(Cubes)
		for _, cubes := range strings.Split(r, ",") {
			cubes = strings.TrimSpace(cubes)
			parts = strings.Fields(cubes)

			count, err := strconv.Atoi(parts[0])
			if err != nil {
				return Game{}, err
			}
			round[Colour(parts[1])] = count
		}
		game.rounds = append(game.rounds, round)
	}
	return game, nil
}

type Game struct {
	id     int
	rounds []Cubes
}

func (g Game) Possible(cubes Cubes) bool {
	for _, round := range g.rounds {
		for colour, count := range round {
			if count > cubes[colour] {
				return false
			}
		}
	}
	return true
}

func (g Game) MinPossible() Cubes {
	min := make(Cubes)
	for _, round := range g.rounds {
		for colour, count := range round {
			min[colour] = maths.Max(min[colour], count)
		}
	}
	return min
}

type Cubes map[Colour]int

func (c Cubes) Power() int {
	power := 1
	for _, colour := range []Colour{Red, Green, Blue} {
		power *= c[colour]
	}
	return power
}

type Colour string

const (
	Red   Colour = "red"
	Green Colour = "green"
	Blue  Colour = "blue"
)
