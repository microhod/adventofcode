package main

import (
	"fmt"
	"sort"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

var (
	round  = [2]rune{'(', ')'}
	square = [2]rune{'[', ']'}
	curly  = [2]rune{'{', '}'}
	angle  = [2]rune{'<', '>'}
)

func main() {
	puzzle.NewSolution("Syntax Scoring", part1, part2).Run()
}

func part1() error {
	lines, err := file.ReadLines(InputFile)
	if err != nil {
		return err
	}

	score := 0
	for _, line := range lines {
		illegal, index := findIllegalChar(line)
		if index >= 0 {
			score += errorScore(illegal)
		}
	}

	fmt.Printf("syntax error score: %d\n", score)

	return nil
}

func part2() error {
	lines, err := file.ReadLines(InputFile)
	if err != nil {
		return err
	}

	incomplete := []string{}
	for _, line := range lines {
		_, index := findIllegalChar(line)
		if index < 0 {
			incomplete = append(incomplete, line)
		}
	}

	scores := []int{}
	for _, line := range incomplete {
		completion, err := findCompletionString(line)
		if err != nil {
			return err
		}
		scores = append(scores, completionScore(completion))
	}

	sort.Ints(scores)
	fmt.Printf("middle completion score: %d\n", scores[len(scores)/2])

	return nil
}

func findIllegalChar(line string) (rune, int) {
	if len(line) < 1 {
		return rune(-1), -1
	}

	chunk := NewChunk(rune(line[0]), nil)
	if chunk == nil {
		return rune(line[0]), 0
	}

	for i, char := range line[1:] {
		if open(char) {
			chunk = NewChunk(char, chunk)
		} else if chunk.CanClose(char) {
			chunk = chunk.Parent
		} else {
			return char, i + 1
		}
	}
	return rune(-1), -1
}

func findCompletionString(line string) (string, error) {
	if len(line) < 1 {
		return "", fmt.Errorf("line is empty")
	}

	chunk := NewChunk(rune(line[0]), nil)
	if chunk == nil {
		return "", fmt.Errorf("first character is invalid")
	}

	for i, char := range line[1:] {
		if open(char) {
			chunk = NewChunk(char, chunk)
		} else if chunk.CanClose(char) {
			chunk = chunk.Parent
		} else {
			return "", fmt.Errorf("character '%s' at index '%d' is invalid", string(char), i+1)
		}
	}

	completion := ""
	for chunk != nil {
		completion += string(chunk.Type[1])
		chunk = chunk.Parent
	}

	return completion, nil
}

func errorScore(char rune) int {
	switch char {
	case round[1]:
		return 3
	case square[1]:
		return 57
	case curly[1]:
		return 1197
	case angle[1]:
		return 25137
	default:
		return 0
	}
}

func completionScore(completion string) int {
	score := 0
	for _, char := range completion {
		score = (score * 5) + completionCharScore(char)
	}
	return score
}

func completionCharScore(char rune) int {
	switch char {
	case round[1]:
		return 1
	case square[1]:
		return 2
	case curly[1]:
		return 3
	case angle[1]:
		return 4
	default:
		return 0
	}
}

func NewChunk(char rune, parent *Chunk) *Chunk {
	switch char {
	case round[0]:
		return &Chunk{Type: round, Parent: parent}
	case square[0]:
		return &Chunk{Type: square, Parent: parent}
	case curly[0]:
		return &Chunk{Type: curly, Parent: parent}
	case angle[0]:
		return &Chunk{Type: angle, Parent: parent}
	default:
		return nil
	}
}

type Chunk struct {
	Parent *Chunk
	Type   [2]rune
}

func (c *Chunk) CanClose(char rune) bool {
	return char == c.Type[1]
}

func open(char rune) bool {
	return contains([]rune{round[0], square[0], curly[0], angle[0]}, char)
}

func contains(chars []rune, char rune) bool {
	for _, c := range chars {
		if c == char {
			return true
		}
	}
	return false
}
