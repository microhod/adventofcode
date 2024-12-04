package main

import (
	"fmt"
	"slices"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Ceres Search", part1, part2).Run()
}

func part1() error {
	ws, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Printf("XMAS: %d\n", ws.Find("XMAS"))
	return nil
}

func part2() error {
	ws, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Printf("X-MAS: %d\n", ws.FindXMas())
	return nil
}

func parse(path string) (WordSearch, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var ws WordSearch
	for _, line := range lines {
		ws = append(ws, []byte(line))
	}
	return ws, nil
}

type WordSearch [][]byte

func (w WordSearch) FindXMas() int {
	var possible [][2]int
	for i := range w {
		for j := range w[i] {
			if w[i][j] == 'A' {
				possible = append(possible, [2]int{i, j})
			}
		}
	}

	var count int
	for _, p := range possible {
		i := p[0]
		j := p[1]

		if i < 1 || j < 1 || i >= len(w)-1 || j >= len(w[i])-1 {
			continue
		}

		tl := w[i-1][j-1]
		tr := w[i-1][j+1]
		bl := w[i+1][j-1]
		br := w[i+1][j+1]

		if !((tl == 'M' && br == 'S') || (tl == 'S' && br == 'M')) {
			continue
		}

		if (tr == 'M' && bl == 'S') || (tr == 'S' && bl == 'M') {
			count++
		}
	}
	return count
}

func (w WordSearch) Find(word string) int {
	letters := []byte(word)
	slices.Reverse(letters)
	backwards := string(letters)

	return w.find(word) + w.find(backwards)
}

func (w WordSearch) find(word string) int {
	var possible [][2]int
	for i := range w {
		for j := range w[i] {
			if w[i][j] == word[0] {
				possible = append(possible, [2]int{i, j})
			}
		}
	}

	var count int
	for _, p := range possible {
		for _, found := range w.getWords(p[0], p[1], len(word)) {
			if found == word {
				count++
			}
		}
	}
	return count
}

func (w WordSearch) getWords(i, j int, length int) []string {
	length -= 1

	var words []string
	// right
	if j+length < len(w[i]) {
		words = append(words, string(w[i][j:j+length+1]))
	}
	// down
	if i+length < len(w) {
		var down []byte
		for ii := i; ii <= i+length; ii++ {
			down = append(down, w[ii][j])
		}
		words = append(words, string(down))
	}
	// diag TL -> BR
	if i+length < len(w) && j+length < len(w[i]) {
		var diag []byte
		for d := 0; d <= length; d++ {
			diag = append(diag, w[i+d][j+d])
		}
		words = append(words, string(diag))
	}
	// diag BL -> TR
	if i-length >= 0 && j+length < len(w[i]) {
		var diag []byte
		for d := 0; d <= length; d++ {
			diag = append(diag, w[i-d][j+d])
		}
		words = append(words, string(diag))
	}

	return words
}
