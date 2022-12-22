package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Monkey Map", part1, part2).Run()
}

func part1() error {
	board, path, err := parse(InputFile, simpleWrap)
	if err != nil {
		return err
	}

	var start *Tile
	for _, tile := range board[0] {
		if tile.Contents == Empty {
			start = tile
			break
		}
	}

	end, direction, err := start.FollowPath(Right, path)
	if err != nil {
		return err
	}

	fmt.Println(1000*(end.Row+1) + 4*(end.Col+1) + int(direction))
	return nil
}

func part2() error {
	board, path, err := parse(InputFile, cubeWrap[InputFile])
	if err != nil {
		return err
	}

	var start *Tile
	for _, tile := range board[0] {
		if tile.Contents == Empty {
			start = tile
			break
		}
	}

	end, direction, err := start.FollowPath(Right, path)
	if err != nil {
		return err
	}

	fmt.Println(1000*(end.Row+1) + 4*(end.Col+1) + int(direction))
	return nil
}

func parse(filePath string, wrapBoundaries func(Board)) (Board, Path, error) {
	input, err := file.Read(filePath)
	if err != nil {
		return nil, Path{}, err
	}

	parts := strings.Split(input, "\n\n")

	board, err := parseBoard(parts[0], wrapBoundaries)
	if err != nil {
		return nil, Path{}, err
	}
	path, err := parsePath(parts[1])

	return board, path, err
}

func parseBoard(input string, wrapBoundaries func(Board)) (Board, error) {
	input = strings.TrimSuffix(input, "\n")
	lines := strings.Split(input, "\n")

	// find row length
	var maxRow int
	for _, row := range lines {
		if len(row) > maxRow {
			maxRow = len(row)
		}
	}

	// make empty board
	board := make(Board, len(lines))
	for row := range board {
		board[row] = make([]*Tile, maxRow)
	}

	// populate intial tiles
	for row := range board {
		for col := range board[row] {
			contents := OutOfBounds
			if col < len(lines[row]) {
				contents = TileContents(lines[row][col])
			}

			board[row][col] = &Tile{
				Contents: contents,
				Row:      row,
				Col:      col,
			}
		}
	}

	// connect tiles to their neighbours
	for row := range board {
		for col := range board[row] {
			if col < len(board[row])-1 && board[row][col+1].Contents != OutOfBounds {
				board[row][col].Neighbours[Right].tile = board[row][col+1]
			}
			if row < len(board)-1 && board[row+1][col].Contents != OutOfBounds {
				board[row][col].Neighbours[Down].tile = board[row+1][col]
			}
			if col > 0 && board[row][col-1].Contents != OutOfBounds {
				board[row][col].Neighbours[Left].tile = board[row][col-1]
			}
			if row > 0 && board[row-1][col].Contents != OutOfBounds {
				board[row][col].Neighbours[Up].tile = board[row-1][col]
			}
		}
	}

	wrapBoundaries(board)

	return board, nil
}

func simpleWrap(board Board) {
	// wrap around tiles
	rowStart := map[int]*Tile{}
	colStart := map[int]*Tile{}
	for row := range board {
		for col := range board[row] {
			tile := board[row][col]
			if tile.Contents == OutOfBounds {
				continue
			}

			if tile.Neighbours[Left].tile == nil && rowStart[row] == nil {
				rowStart[row] = tile
			}
			if tile.Neighbours[Up].tile == nil && colStart[col] == nil {
				colStart[col] = tile
			}

			if tile.Neighbours[Right].tile == nil {
				tile.Neighbours[Right].tile = rowStart[row]
				rowStart[row].Neighbours[Left].tile = tile
			}
			if tile.Neighbours[Down].tile == nil {
				tile.Neighbours[Down].tile = colStart[col]
				colStart[col].Neighbours[Up].tile = tile
			}
		}
	}
}

var cubeWrap = map[string]func(Board){
	TestFile:  testCubeWrap,
	InputFile: inputCubeWrap,
}

func testCubeWrap(board Board) {
	//   1
	// 234
	//   56
	for row := range board {
		for col := range board[row] {
			// left 1 <-> 3 top
			// 0,8 -> 4,4
			// 1,8 -> 4,5
			if row < 4 && col == 8 {
				tile := board[row][col]
				neighbour := board[4][4+row]

				tile.Neighbours[Left].tile = neighbour
				tile.Neighbours[Left].turn = AntiClockwise
				neighbour.Neighbours[Up].tile = tile
				neighbour.Neighbours[Up].turn = ClockWise
			}
			// top 1 <-> 2 top
			// 0,8 -> 4,3
			// 0,9 -> 4,2
			if row == 0 && col >= 8 && col <= 11 {
				tile := board[row][col]
				neighbour := board[4][11-col]

				tile.Neighbours[Up].tile = neighbour
				tile.Neighbours[Up].turn = 2 * AntiClockwise
				neighbour.Neighbours[Up].tile = tile
				neighbour.Neighbours[Up].turn = 2 * ClockWise
			}
			// right 1 <-> 6 right
			// 0,11 -> 11,15
			// 1,11 -> 10,15
			if row < 4 && col == 11 {
				tile := board[row][col]
				neighbour := board[11-row][15]

				tile.Neighbours[Right].tile = neighbour
				tile.Neighbours[Right].turn = 2 * ClockWise
				neighbour.Neighbours[Right].tile = tile
				neighbour.Neighbours[Right].turn = 2 * AntiClockwise
			}
			// right 4 <-> 6 top
			// 4,11 -> 8,15
			// 5,11 -> 8,14
			if row >= 4 && row <= 7 && col == 11 {
				tile := board[row][col]
				neighbour := board[8][19-row]

				tile.Neighbours[Right].tile = neighbour
				tile.Neighbours[Right].turn = ClockWise
				neighbour.Neighbours[Up].tile = tile
				neighbour.Neighbours[Up].turn = AntiClockwise
			}
			// left 2 <-> 6 bottom
			// 4,0 -> 11,15
			// 5,0 -> 11,14
			if row >= 4 && row <= 7 && col == 0 {
				tile := board[row][col]
				neighbour := board[11][19-row]

				tile.Neighbours[Left].tile = neighbour
				tile.Neighbours[Left].turn = 3 * AntiClockwise
				neighbour.Neighbours[Down].tile = tile
				neighbour.Neighbours[Down].turn = 3 * ClockWise
			}
			// bottom 2 <-> 5 bottom
			// 7,0 -> 11,11
			// 7,1 -> 11,10
			if row == 7 && col < 4 {
				tile := board[row][col]
				neighbour := board[11][11-col]

				tile.Neighbours[Down].tile = neighbour
				tile.Neighbours[Down].turn = 2 * AntiClockwise
				neighbour.Neighbours[Down].tile = tile
				neighbour.Neighbours[Down].turn = 2 * ClockWise
			}
			// bottom 3 <-> 5 left
			// 7,4 -> 11,8
			// 7,5 -> 10,8
			if row == 7 && col >= 4 && col <= 7 {
				tile := board[row][col]
				neighbour := board[15-col][8]

				tile.Neighbours[Down].tile = neighbour
				tile.Neighbours[Down].turn = AntiClockwise
				neighbour.Neighbours[Left].tile = tile
				neighbour.Neighbours[Left].turn = ClockWise
			}
		}
	}
}

func inputCubeWrap(board Board) {
	//  12
	//  3
	// 45
	// 6
	for row := range board {
		for col := range board[row] {
			// top 1 <-> 6 left
			// 0,50 -> 150,0
			// 0,51 -> 151,0
			if row == 0 && col >= 50 && col < 100 {
				tile := board[row][col]
				neighbour := board[100+col][0]

				tile.Neighbours[Up].tile = neighbour
				tile.Neighbours[Up].turn = 3 * AntiClockwise
				neighbour.Neighbours[Left].tile = tile
				neighbour.Neighbours[Left].turn = 3 * ClockWise
			}
			// left 1 <-> 4 left
			// 0,50 -> 149,0
			// 1,50 -> 148,0
			if row < 50 && col == 50 {
				tile := board[row][col]
				neighbour := board[149-row][0]

				tile.Neighbours[Left].tile = neighbour
				tile.Neighbours[Left].turn = 2 * AntiClockwise
				neighbour.Neighbours[Left].tile = tile
				neighbour.Neighbours[Left].turn = 2 * ClockWise
			}
			// top 2 <-> 6 bottom
			// 0,100 -> 199,0
			// 0,101 -> 199,1
			if row == 0 && col >= 100 && col < 150 {
				tile := board[row][col]
				neighbour := board[199][col-100]

				tile.Neighbours[Up].tile = neighbour
				neighbour.Neighbours[Down].tile = tile
			}
			// right 2 <-> right 5
			// 0,149 -> 149,99
			// 1,149 -> 148,99
			if row < 50 && col == 149 {
				tile := board[row][col]
				neighbour := board[149-row][99]

				tile.Neighbours[Right].tile = neighbour
				tile.Neighbours[Right].turn = 2 * ClockWise
				neighbour.Neighbours[Right].tile = tile
				neighbour.Neighbours[Right].turn = 2 * AntiClockwise
			}
			// bottom 2 <-> 3 right
			// 49,100 -> 50,99
			// 49,101 -> 51,99
			if row == 49 && col >= 100 && col < 150 {
				tile := board[row][col]
				neighbour := board[col-50][99]

				tile.Neighbours[Down].tile = neighbour
				tile.Neighbours[Down].turn = ClockWise
				neighbour.Neighbours[Right].tile = tile
				neighbour.Neighbours[Right].turn = AntiClockwise
			}
			// left 3 <-> 4 top
			// 50,50 -> 100,0
			// 51,50 -> 100,1
			if row >= 50 && row < 100 && col == 50 {
				tile := board[row][col]
				neighbour := board[100][row-50]

				tile.Neighbours[Left].tile = neighbour
				tile.Neighbours[Left].turn = AntiClockwise
				neighbour.Neighbours[Up].tile = tile
				neighbour.Neighbours[Up].turn = ClockWise
			}
			// bottom 5 <-> 6 right
			// 149,50 -> 150,49
			// 149,51 -> 151,49
			if row == 149 && col >= 50 && col < 100 {
				tile := board[row][col]
				neighbour := board[100+col][49]

				tile.Neighbours[Down].tile = neighbour
				tile.Neighbours[Down].turn = ClockWise
				neighbour.Neighbours[Right].tile = tile
				neighbour.Neighbours[Right].turn = AntiClockwise
			}
		}
	}
}

var letterToTurn = map[rune]int{
	'R': 1,
	'L': -1,
}

func parsePath(input string) (Path, error) {
	var path Path
	var moveBuffer string
	for _, char := range strings.TrimSpace(input) {
		if unicode.IsNumber(char) {
			moveBuffer += string(char)
			continue
		}

		path.Turns = append(path.Turns, letterToTurn[char])

		if len(moveBuffer) > 0 {
			move, err := strconv.Atoi(moveBuffer)
			if err != nil {
				return Path{}, err
			}
			path.Moves = append(path.Moves, move)

			moveBuffer = ""
		}
	}
	if len(moveBuffer) > 0 {
		move, err := strconv.Atoi(moveBuffer)
		if err != nil {
			return Path{}, err
		}
		path.Moves = append(path.Moves, move)
	}
	return path, nil
}

type Path struct {
	Moves []int
	Turns []int
}

type Board [][]*Tile

type Tile struct {
	Contents   TileContents
	Neighbours [4]struct {
		tile *Tile
		turn int
	}

	Row, Col int
}

func (tile *Tile) FollowPath(direction Direction, path Path) (*Tile, Direction, error) {
	if len(path.Moves) < 1 {
		return tile, direction, nil
	}

	move := path.Moves[0]
	path.Moves = path.Moves[1:]

	for i := 0; i < move; i++ {
		next := tile.Neighbours[direction]
		// don't move if we hit a wall
		if next.tile.Contents == Wall {
			break
		}
		// this should never happen if we parsed the board correctly
		if next.tile.Contents == OutOfBounds {
			return tile, 0, fmt.Errorf("went out of bounds")
		}

		tile = next.tile
		direction = direction.Turn(next.turn)
	}

	var turn int
	if len(path.Turns) > 0 {
		turn = path.Turns[0]
		path.Turns = path.Turns[1:]
	}
	direction = direction.Turn(turn)

	return tile.FollowPath(direction, path)
}

type TileContents string

const (
	Empty       TileContents = "."
	Wall        TileContents = "#"
	OutOfBounds TileContents = " "
)

type Direction int

const (
	Right Direction = iota
	Down
	Left
	Up
)

const (
	ClockWise     = 1
	AntiClockwise = -1
)

func (d Direction) Turn(amount int) Direction {
	return Direction(maths.Mod(int(d)+amount, 4))
}
