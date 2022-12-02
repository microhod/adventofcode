package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Rock Paper Scissors", part1, part2).Run()
}

func part1() error {
	game, err := parsePart1(InputFile, part1Round)
	if err != nil {
		return err
	}

	game.Play()

	fmt.Printf("The scores are in:\nOpponent: %d\nMe:       %d\n", game.Scores[0], game.Scores[1])
	return nil
}

func part2() error {
	game, err := parsePart1(InputFile, part2Round)
	if err != nil {
		return err
	}

	game.Play()

	fmt.Printf("The scores are in:\nOpponent: %d\nMe:       %d\n", game.Scores[0], game.Scores[1])
	return nil
}

func parsePart1(path string, parseRound func([]string) Round) (*RpsGame, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	game := &RpsGame{
		Rounds: []Round{},
		Scores: [2]int{},
	}

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		game.Rounds = append(game.Rounds, parseRound(fields))
	}

	return game, nil
}

func part1Round(fields []string) Round {
	shapes := map[string]Shape{
		"A": Rock,
		"B": Paper,
		"C": Scissors,
		"X": Rock,
		"Y": Paper,
		"Z": Scissors,
	}
	return Round{
		Shapes: [2]Shape{
			shapes[fields[0]],
			shapes[fields[1]],
		},
	}
}

func part2Round(fields []string) Round {
	shapes := map[string]Shape{
		"A": Rock,
		"B": Paper,
		"C": Scissors,
	}
	opponentShape := shapes[fields[0]]

	var myShape Shape
	switch fields[1] {
	case "X": // lose
		myShape = Shape(mod(int(opponentShape)-1, 3))
	case "Y": // draw
		myShape = opponentShape
	case "Z": // win
		myShape = Shape(mod(int(opponentShape)+1, 3))
	}

	if myShape < 0 {
		panic(strings.Join([]string{strings.Join(fields, " "), fmt.Sprint(opponentShape), fmt.Sprint(myShape %3)}, " "))
	}

	return Round{
		Shapes: [2]Shape{
			opponentShape,
			myShape,
		},
	}
}

type RpsGame struct {
	Rounds []Round
	Scores [2]int
}

func (r *RpsGame) Play() {
	for _, round := range r.Rounds {
		score := round.Score()
		r.Scores[0] += score[0]
		r.Scores[1] += score[1]
	}
}

type Round struct {
	Shapes [2]Shape
}

func (r Round) Winner() int {
	if r.Shapes[0].Beats(r.Shapes[1]) {
		return 0
	}
	if r.Shapes[1].Beats(r.Shapes[0]) {
		return 1
	}
	// -1 means draw
	return -1
}

func (r Round) Score() [2]int {
	scores := [2]int{
		r.Shapes[0].Score(),
		r.Shapes[1].Score(),
	}

	winner := r.Winner()
	if winner < 0 {
		scores[0] += 3
		scores[1] += 3
		return scores
	}

	scores[winner] += 6
	return scores
}

const (
	Rock Shape = iota
	Paper
	Scissors
)

type Shape int

func (s Shape) Beats(t Shape) bool {
	return s == (t+1) % 3
}

func (s Shape) Score() int {
	return int(s+1)
}

func mod(a, b int) int {
	m := a % b
	if m < 0 {
		m += b
	}
	return m
}
