package main

import (
	"fmt"
	"sort"
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
	puzzle.NewSolution("Camel Cards", part1, part2).Run()
}

func part1() error {
	rounds, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	fmt.Printf("total winnings: %d\n", rounds.Winnings())
	return nil
}

func part2() error {
	rounds, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	Cards['J'] = 0
	for _, round := range rounds {
		round.UseJokers = true
	}

	fmt.Printf("total winnings: %d\n", rounds.Winnings())
	return nil
}

func parseInput(path string) (Rounds, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var rounds Rounds
	for _, line := range lines {
		parts := strings.Fields(line)
		bid, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		rounds = append(rounds, &Round{
			Hand: Hand(parts[0]),
			Bid:  bid,
		})
	}
	return rounds, nil
}

type Rounds []*Round

func (r Rounds) Winnings() int {
	sort.Slice(r, func(i, j int) bool {
		return r[i].Less(r[j])
	})

	var total int
	for i, round := range r {
		total += (i + 1) * round.Bid
	}
	return total
}

type Round struct {
	Hand      Hand
	Type      HandType
	Bid       int
	UseJokers bool
}

func (r *Round) Less(s *Round) bool {
	rType, sType := r.GetType(), s.GetType()
	if rType == sType {
		return r.Hand.Less(s.Hand)
	}
	return rType < sType
}

func (r *Round) GetType() HandType {
	if r.Type == 0 {
		if r.UseJokers {
			r.Type = r.Hand.WithoutJokers().Type()
		} else {
			r.Type = r.Hand.Type()
		}
	}
	return r.Type
}

type Hand string

func (h Hand) Less(g Hand) bool {
	for i := range h {
		if Cards[rune(h[i])] > Cards[rune(g[i])] {
			return false
		}
		if Cards[rune(h[i])] < Cards[rune(g[i])] {
			return true
		}
		continue
	}
	return false
}

func (h Hand) WithoutJokers() Hand {
	counts := map[rune]int{}
	var max rune
	for _, ch := range h {
		if ch == 'J' {
			continue
		}

		counts[ch]++
		if counts[ch] > counts[max] {
			max = ch
		}
	}

	// replace the Jokers with the most common non-Joker card
	return Hand(strings.ReplaceAll(string(h), "J", string(max)))
}

func (h Hand) Type() HandType {
	counts := map[rune]int{}
	types := map[HandType]int{}

	for _, ch := range h {
		counts[ch] += 1
		switch counts[ch] {
		case 5:
			types[FiveOfAKind]++
			types[FourOfAKind]--
		case 4:
			types[FourOfAKind]++
			types[ThreeOfAKind]--
		case 3:
			types[ThreeOfAKind]++
			if types[TwoPair] > 0 {
				types[TwoPair]--
				types[OnePair]++
			} else {
				types[OnePair]--
			}
			// full house
			if types[OnePair] > 0 {
				types[FullHouse]++
				types[OnePair]--
				types[ThreeOfAKind]--
			}
		case 2:
			if types[OnePair] > 0 {
				types[TwoPair]++
				types[OnePair]--
			} else {
				types[OnePair]++
			}
			// full house
			if types[ThreeOfAKind] > 0 {
				types[FullHouse]++
				types[OnePair]--
				types[ThreeOfAKind]--
			}
		}
	}

	var max HandType
	for t := range types {
		max = maths.Max(max, t)
	}

	return maths.Max(HighCard, max)
}

type HandType int

const (
	None HandType = iota
	HighCard
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

var Cards = map[rune]int{
	'A': 13,
	'K': 12,
	'Q': 11,
	'J': 10,
	'T': 9,
	'9': 8,
	'8': 7,
	'7': 6,
	'6': 5,
	'5': 4,
	'4': 3,
	'3': 2,
	'2': 1,
}
