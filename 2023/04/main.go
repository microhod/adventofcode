package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Scratchcards", part1, part2).Run()
}

func part1() error {
	cards, err := parseCards(InputFile)
	if err != nil {
		return err
	}

	var points int
	for _, card := range cards {
		points += card.Points()
	}

	fmt.Printf("the total point value of all the cards is: %d\n", points)
	return nil
}

func part2() error {
	cards, err := parseCards(InputFile)
	if err != nil {
		return err
	}

	cards.Process()

	fmt.Printf("the total cards we end up with is: %d\n", cards.Count())
	return nil
}

func parseCards(path string) (Cards, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var cards []*Card
	for _, line := range lines {
		card, err := parseCard(line)
		if err != nil {
			return nil, err
		}
		cards = append(cards, &card)
	}
	return cards, nil
}

func parseCard(line string) (Card, error) {
	line = strings.Split(line, ": ")[1]
	parts := strings.Split(line, " | ")

	winning, err := csv.ParseInts(parts[0], " ")
	if err != nil {
		return Card{}, err
	}
	have, err := csv.ParseInts(parts[1], " ")
	if err != nil {
		return Card{}, err
	}
	return Card{
		Replicas: 1,
		Winning:  set.NewSet(winning...),
		Have:     set.NewSet(have...),
	}, nil
}

type Cards []*Card

func (c Cards) Process() {
	for i := range c {
		next := i+1
		matches := c[i].Matches()

		for j := next; j < next+matches; j++ {
			c[j].Replicas += c[i].Replicas
		}
	}
}

func (c Cards) Count() int {
	var total int
	for _, card := range c {
		total += card.Replicas
	}
	return total
}

type Card struct {
	Replicas int
	Winning  set.Set[int]
	Have     set.Set[int]
}

func (c Card) Matches() int {
	var matches int
	for w := range c.Winning {
		if c.Have.Contains(w) {
			matches += 1
		}
	}
	return matches
}

func (c Card) Points() int {
	matches := c.Matches()
	if matches == 0 {
		return 0
	}
	return 1 << (matches - 1)
}
