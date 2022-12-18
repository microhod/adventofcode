package main

import (
	"fmt"
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
	puzzle.NewSolution("Pyroclastic Flow", part1, part2).Run()
}

func part1() error {
	moves, err := parse(InputFile)
	if err != nil {
		return err
	}

	game := &Game{
		moves: moves,
		board: Board{
			// this is the ground
			{1, 1, 1, 1, 1, 1, 1},
		},
	}

	shapes := []Shape{minus, plus, corner, line, square}
	shapeIndex := 0
	for i := 0; i < 2022; i++ {
		shape := shapes[shapeIndex]
		shapeIndex = maths.Mod(shapeIndex+1, len(shapes))

		game.AddShape(shapeIndex, shape)
	}
	fmt.Println(game.height)
	return nil
}

func part2() error {
	moves, err := parse(InputFile)
	if err != nil {
		return err
	}

	game := &Game{
		moves: moves,
		board: Board{
			// this is the ground
			{1, 1, 1, 1, 1, 1, 1},
		},
	}
	shapes := []Shape{minus, plus, corner, line, square}

	/*
		The tower repeats every 1700 shapes (first is at shape 441 where height is 671), when shape is 1 and move is 5206
		height increase between repeats is 2623

		(1000000000000 - 441) / 1700 = 588235293 remainder 1459

		Therefore, height(1000000000000) = height(441) + 588235293 * repeat_height_increase + height_1_5206(1459)
										 = 671 + 588235293 * 2623 + height_1_5206(1459)
	*/

	shapeIndex := 1
	game.currentMoveIndex = 5206

	for i := 0; i < 1459; i++ {
		shape := shapes[shapeIndex]
		shapeIndex = maths.Mod(shapeIndex+1, len(shapes))

		game.AddShape(shapeIndex, shape)
	}

	fmt.Println(671 + (588235293 * 2623) + game.height)
	return nil
}

func parse(path string) ([]Vector, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	inputToMove := map[rune]Vector{
		'<': left,
		'>': right,
	}

	var moves []Vector
	for _, ch := range lines[0] {
		moves = append(moves, inputToMove[ch], down)
	}
	return moves, nil
}

var (
	minus = Shape{
		{1, 1, 1, 1},
	}
	plus = Shape{
		{0, 1, 0},
		{1, 1, 1},
		{0, 1, 0},
	}
	corner = Shape{
		{0, 0, 1},
		{0, 0, 1},
		{1, 1, 1},
	}
	line = Shape{
		{1},
		{1},
		{1},
		{1},
	}
	square = Shape{
		{1, 1},
		{1, 1},
	}

	left  = Vector{-1, 0}
	right = Vector{1, 0}
	down  = Vector{0, 1}
)

type Game struct {
	moves               []Vector
	currentMoveIndex    int
	board               Board
	heightSinceFullLine int
	height              int
}

type Board [][]int

type Shape [][]int

func (g *Game) AddShape(shapeIndex int, shape Shape) {
	// add empty lines
	g.board = g.board.AddEmptyLines(len(shape) + 3)

	// start position
	position := Vector{2, 0}

	// iterate until blocked going down
	for {
		move := g.moves[g.currentMoveIndex]
		g.currentMoveIndex = maths.Mod(g.currentMoveIndex+1, len(g.moves))

		newPosition := position.Add(move)

		outOfBounds := shape.OutOfBounds(g.board, newPosition)
		var overlaps bool
		if !outOfBounds {
			overlaps = shape.Overlaps(g.board, newPosition)
		}

		if !outOfBounds && !overlaps {
			position = newPosition
		}
		if (outOfBounds || overlaps) && move == down {
			heightBefore := len(g.board.Copy().RemoveEmptyLines()) - 1

			g.board = g.board.AddShape(shape, position)
			g.board = g.board.RemoveEmptyLines()

			g.height += len(g.board) - 1 - heightBefore

			/*
			truncate board to last full line:
			|.......|
			|.......|
			|.......|
			|..#....|
			|..#....|
			|#######| <- truncate to here
			|..###..|
			|...#...|
			|..####.|
			|#######|

			which gives (without empty lines):
			|..#....|
			|..#....|
			|#######|
			*/
			for i := 0; i < len(shape) && i < len(g.board); i++ {
				if maths.Sum(g.board[i]...) == 7 {
					g.board = g.board[:i+1]
					break
				}
			}

			break
		}
	}
}

func (b Board) AddShape(shape Shape, position Vector) Board {
	for y := 0; y < len(shape); y++ {
		for x := 0; x < len(shape[0]); x++ {
			if shape[y][x] == 1 {
				b[y+position.Y][x+position.X] = 1
			}
		}
	}
	return b
}

func (b Board) Copy() Board {
	copy := make(Board, len(b))
	for i := range b {
		copy[i] = append([]int{}, b[i]...)
	}
	return copy
}

func (b Board) AddEmptyLines(size int) Board {
	lines := make([][]int, size)
	for i := range lines {
		lines[i] = make([]int, len(b[0]))
	}
	return append(lines, b...)
}

func (b Board) RemoveEmptyLines() Board {
	firstNonEmptyLine := 0

	for maths.Sum(b[firstNonEmptyLine]...) == 0 {
		firstNonEmptyLine += 1
	}

	return append(Board{}, b[firstNonEmptyLine:]...)
}

func (b Board) String() string {
	var lines []string
	for i := 0; i < len(b); i++ {
		var line string
		for _, num := range b[i] {
			if num == 1 {
				line += "#"
			} else {
				line += "."
			}
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

type Vector struct {
	X, Y int
}

func (v Vector) Add(u Vector) Vector {
	return Vector{v.X + u.X, v.Y + u.Y}
}

func (s Shape) Overlaps(b Board, position Vector) bool {
	for row, boardRow := range b[position.Y : position.Y+len(s)] {
		for col, boardRowCol := range boardRow[position.X : position.X+len(s[0])] {
			overlap := s[row][col] & boardRowCol

			if overlap == 1 {
				return true
			}
		}
	}
	return false
}

func (s Shape) OutOfBounds(b Board, position Vector) bool {
	if position.X < 0 || position.Y < 0 {
		return true
	}
	if position.Y+len(s) > len(b) {
		return true
	}
	if position.X+len(s[0]) > len(b[0]) {
		return true
	}
	return false
}