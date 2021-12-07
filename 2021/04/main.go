package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	BingoFile = "bingo.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Giant Squid", part1, part2).Run()
}

func part1() error {
	bingo, err := readBingo(BingoFile)
	if err != nil {
		return err
	}

	winners := bingo.Play()
	if len(winners) < 1 {
		fmt.Println("no winners...😢")
	}

	winner := winners[0]
	fmt.Printf("winning card:\n\n%s\n", winner.Card.String())

	fmt.Println()

	unmarkedSum := winner.Card.UnmarkedSum()
	fmt.Printf("score = %d * %d = %d\n", unmarkedSum, winner.Number, unmarkedSum*winner.Number)

	return nil
}

func part2() error {
	bingo, err := readBingo(BingoFile)
	if err != nil {
		return err
	}

	winners := bingo.Play()
	if len(winners) < 1 {
		fmt.Println("no winners...😢")
	}

	lastWinner := winners[len(winners)-1]
	fmt.Printf("last winning card:\n\n%s\n", lastWinner.Card.String())

	fmt.Println()

	unmarkedSum := lastWinner.Card.UnmarkedSum()
	fmt.Printf("score = %d * %d = %d\n", unmarkedSum, lastWinner.Number, unmarkedSum*lastWinner.Number)

	return nil
}

func readBingo(path string) (*Bingo, error) {
	var err error

	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	bingo := &Bingo{
		Cards:   []*Card{},
		Numbers: []int{},
	}

	if len(lines) < 1 {
		return bingo, nil
	}

	bingo.Numbers, err = csv.ParseInts(lines[0])
	if err != nil {
		return nil, err
	}

	lines = lines[1:]

	for len(lines) > 5 {
		// 1:6 to skip the blank line inbetween every grid
		card, err := parseBingoCard(lines[1:6])
		if err != nil {
			return nil, err
		}
		bingo.Cards = append(bingo.Cards, card)

		lines = lines[6:]
	}

	return bingo, nil
}

func parseBingoCard(lines []string) (*Card, error) {
	if len(lines) != 5 {
		return nil, fmt.Errorf("expected 5 lines, but got %d", len(lines))
	}

	card := &Card{
		Numbers: *new([5][5]CardNumber),
	}
	for row, line := range lines {
		for column, part := range strings.Fields(line) {
			n, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}

			card.Numbers[row][column] = CardNumber{Value: n}
		}
	}

	return card, nil
}

type Bingo struct {
	Cards   []*Card
	Numbers []int
}

func (bingo *Bingo) String() string {
	str := fmt.Sprintf("%v\n", bingo.Numbers)

	for _, card := range bingo.Cards {
		str += fmt.Sprint("\n", card.String(), "\n")
	}

	return str
}

type BingoWin struct {
	Card   *Card
	Number int
}

func (bingo *Bingo) Play() []BingoWin {
	winners := []BingoWin{}
	for _, n := range bingo.Numbers {
		for _, c := range bingo.Cards {
			// stop marking cards once they have won
			if c.Won {
				continue
			}

			c.Mark(n)
			if c.Won {
				winners = append(winners, BingoWin{
					Card:   c,
					Number: n,
				})
			}
		}
	}

	return winners
}

type Card struct {
	Numbers [5][5]CardNumber
	Won     bool
}

type CardNumber struct {
	Value  int
	Marked bool
}

func (num CardNumber) String() string {
	format := "%d"
	if num.Value < 10 {
		format = " %d"
	}

	str := fmt.Sprintf(format, num.Value)
	if num.Marked {
		str = ansi.Color(str, "green+b")
	}
	return str
}

func (card *Card) String() string {
	lines := make([]string, 5)

	for i := 0; i < 5; i++ {
		line := make([]string, 5)
		for j := 0; j < 5; j++ {
			line[j] += card.Numbers[i][j].String()
		}
		lines[i] = strings.Join(line, " ")
	}

	return strings.Join(lines, "\n")
}

func (card *Card) Mark(num int) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if card.Numbers[i][j].Value == num {
				card.Numbers[i][j].Marked = true
				card.Won = card.Won || card.Bingo(uint(i), uint(j))
			}
		}
	}
}

func (card *Card) UnmarkedSum() int {
	sum := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if !card.Numbers[i][j].Marked {
				sum += card.Numbers[i][j].Value
			}
		}
	}

	return sum
}

func (card *Card) Bingo(row uint, column uint) bool {
	if row > 4 || column > 4 {
		panic("row and column have to be between 0 and 4")
	}

	bingoRow := true
	bingoColumn := true
	for i := 0; i < 5; i++ {
		bingoRow = bingoRow && card.Numbers[row][i].Marked
		bingoColumn = bingoColumn && card.Numbers[i][column].Marked
	}

	return bingoRow || bingoColumn
}
