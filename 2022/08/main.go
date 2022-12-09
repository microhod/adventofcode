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
	puzzle.NewSolution("Treetop Tree House", part1, part2).Run()
}

func part1() error {
	heights, err := parse(InputFile)
	if err != nil {
		return err
	}

	var visible int
	for row := range heights {
		for col := range heights[row] {
			if any(heights.Visibility(Tree{row, col})) {
				visible += 1
			}
		}
	}

	fmt.Println(visible)
	return nil
}

func part2() error {
	heights, err := parse(InputFile)
	if err != nil {
		return err
	}

	var mostScenic int
	for row := range heights {
		for col := range heights[row] {

			score := product(heights.ViewingDistances(Tree{row, col}))
			if score > mostScenic {
				mostScenic = score
			}
		}
	}

	fmt.Println(mostScenic)
	return nil
}

type Tree struct {
	Row, Col int
}

func parse(path string) (Heights, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var heights Heights
	for _, line := range lines {
		row, err := parseInts(line)
		if err != nil {
			return Heights{}, err
		}

		heights = append(heights, row)
	}

	return heights, nil
}

func parseInts(str string) ([]int, error) {
	nums := []int{}
	for _, s := range strings.Split(str, "") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}

		nums = append(nums, n)
	}

	return nums, nil
}

type Heights [][]int

func (h Heights) Visibility(tree Tree) (bool, bool, bool, bool) {
	var (
		visibleLeft   = true
		visibleRight  = true
		visibleTop    = true
		visibleBottom = true
	)
	treeHeight := h[tree.Row][tree.Col]

	// left
	for col := 0; col < tree.Col; col++ {
		if h[tree.Row][col] >= treeHeight {
			visibleLeft = false
			break
		}
	}
	// right
	for col := tree.Col + 1; col < len(h[tree.Row]); col++ {
		if h[tree.Row][col] >= treeHeight {
			visibleRight = false
			break
		}
	}
	// top
	for row := 0; row < tree.Row; row++ {
		if h[row][tree.Col] >= treeHeight {
			visibleTop = false
			break
		}
	}
	// bottom
	for row := tree.Row + 1; row < len(h); row++ {
		if h[row][tree.Col] >= treeHeight {
			visibleBottom = false
			break
		}
	}

	return visibleLeft, visibleRight, visibleTop, visibleBottom
}

func (h Heights) ViewingDistances(tree Tree) (int, int, int, int) {
	var left, right, top, bottom int

	treeHeight := h[tree.Row][tree.Col]

	// left
	for col := tree.Col-1; col >= 0; col-- {
		left += 1
		if h[tree.Row][col] >= treeHeight {
			break
		}
	}
	// right
	for col := tree.Col + 1; col < len(h[tree.Row]); col++ {
		right += 1
		if h[tree.Row][col] >= treeHeight {
			break
		}
	}
	// top
	for row := tree.Row-1; row >= 0; row-- {
		top += 1
		if h[row][tree.Col] >= treeHeight {
			break
		}
	}
	// bottom
	for row := tree.Row + 1; row < len(h); row++ {
		bottom += 1
		if h[row][tree.Col] >= treeHeight {
			break
		}
	}

	return left, right, top, bottom
}

func any(conditions ...bool) bool {
	for _, c := range conditions {
		if c {
			return true
		}
	}
	return false
}

func product(nums ...int) int {
	var p = 1
	for _, n := range nums {
		p *= n
	}
	return p
}
