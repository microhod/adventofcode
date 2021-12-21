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
	puzzle.NewSolution("Dirac Dice", part1, part2).Run()
}

func part1() error {
	players, err := readPlayers(InputFile)
	if err != nil {
		return err
	}

	die := NewDeterministicDie(100)
	game := &Game{
		BoardSize:    10,
		Players:      players,
		Die:          die,
		WinningScore: 1000,
	}

	winner := game.Play()
	loser := players[(winner+1)%2]

	fmt.Printf("losing score * die value = %d * %d = %d\n", loser.Score, die.CountRolls(), loser.Score*die.CountRolls())

	return nil
}

func part2() error {
	return nil
}

func readPlayers(path string) ([]*Player, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	players := []*Player{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		position, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		players = append(players, &Player{
			Position: position - 1,
		})
	}

	return players, nil
}

type Player struct {
	Position int
	Score    int
}

type Game struct {
	BoardSize     int
	Players       []*Player
	Die           Die
	WinningScore  int
	currentPlayer int
}

func (g *Game) Play() int {
	for g.Winner() < 0 {
		g.PlayRound()
	}
	return g.Winner()
}

func (g *Game) Winner() int {
	for i, p := range g.Players {
		if p.Score >= g.WinningScore {
			return i
		}
	}
	return -1
}

func (g *Game) PlayRound() {
	player := g.Players[g.currentPlayer]

	roll := g.RollDie()
	player.Position = (player.Position + roll) % g.BoardSize
	player.Score += player.Position + 1

	g.currentPlayer = (g.currentPlayer + 1) % len(g.Players)
}

func (g *Game) RollDie() int {
	sum := 0
	for i := 0; i < 3; i++ {
		sum += g.Die.Roll()
	}
	return sum
}

type Die interface {
	Roll() int
}

type DeterministicDie struct {
	max   int
	value int
	count int
}

func NewDeterministicDie(max int) *DeterministicDie {
	return &DeterministicDie{max: max}
}

func (d *DeterministicDie) Value() int {
	return d.value
}

func (d *DeterministicDie) Roll() int {
	d.count += 1
	d.value += 1
	if d.value > d.max {
		d.value = 1
	}

	return d.value
}

func (d *DeterministicDie) CountRolls() int {
	return d.count
}
