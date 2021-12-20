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
	puzzle.NewSolution("Snailfish", part1, part2).Run()
}

func part1() error {
	trees, err := readTrees(InputFile)
	if err != nil {
		return err
	}

	tree := trees[0]
	for _, t := range trees[1:] {
		tree = Add(tree, t)
	}

	fmt.Printf("magnitude: %d\n", tree.Magnitude())

	return nil
}

func part2() error {
	trees, err := readTrees(InputFile)
	if err != nil {
		return err
	}

	max := 0
	for i := range trees {
		for j := range trees {
			if i == j {
				continue
			}
			// always read again as a quick & dirty fix for the trees changing
			// when adding (as everything is passed by reference)
			trees, err := readTrees(InputFile)
			if err != nil {
				return err
			}
			if m := Add(trees[i], trees[j]).Magnitude(); m > max {
				max = m
			}
		}
	}

	fmt.Printf("max add magnitude: %d\n", max)

	return nil
}

func readTrees(path string) ([]*Tree, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	trees := []*Tree{}
	for _, line := range lines {
		tree, err := readTree(line)
		if err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}

	return trees, nil
}

func readTree(str string) (*Tree, error) {
	var err error
	tree := &Tree{}
	for _, char := range strings.Split(str, "") {
		switch char {
		case "[":
			left := &Tree{Parent: tree}
			tree.Left = left
			tree = left
		case ",":
			right := &Tree{Parent: tree.Parent}
			tree.Parent.Right = right
			tree = right
		case "]":
			tree = tree.Parent
		default:
			tree.Value, err = strconv.Atoi(string(char))
			if err != nil {
				return nil, err
			}
		}
	}
	return tree, nil
}

func Add(t1, t2 *Tree) *Tree {
	tree := &Tree{}
	t1.Parent, t2.Parent = tree, tree
	tree.Left, tree.Right = t1, t2
	tree.Reduce()
	return tree
}

type Tree struct {
	Parent *Tree
	Left   *Tree
	Right  *Tree
	Value  int
}

func (tree *Tree) Reduce() {
	reduction := tree.reduction()
	for reduction != nil {
		reduction()
		reduction = tree.reduction()
	}
}

func (tree *Tree) reduction() func() {
	values := tree.ValuesFromLeft()
	for _, t := range values {
		if t.Parent.IsValuePair() && t.Parent.Depth() > 3 {
			return t.Parent.Explode(values)
		}
	}
	for _, t := range values {
		if t.Value > 9 {
			return t.Split
		}
	}
	return nil
}

func (tree *Tree) ValuesFromLeft() []*Tree {
	values := []*Tree{}
	t := tree.getLeftMost()
	for t != nil {
		if t.IsValue() {
			values = append(values, t)
		}
		t = t.getNext()
	}
	return values
}

func (tree *Tree) getLeftMost() *Tree {
	t := tree
	for t.Left != nil {
		t = t.Left
	}
	return t
}

func (tree *Tree) getNext() *Tree {
	if tree.Right != nil {
		return tree.Right.getLeftMost()
	} else {
		t := tree
		for t.Parent != nil && t == t.Parent.Right {
			t = t.Parent
		}
		return t.Parent
	}
}

func (tree *Tree) Magnitude() int {
	if tree.IsValue() {
		return tree.Value
	}
	return (3 * tree.Left.Magnitude()) + (2 * tree.Right.Magnitude())
}

func (tree *Tree) String() string {
	if tree.IsValue() {
		return fmt.Sprint(tree.Value)
	}
	return fmt.Sprintf("[%s,%s]", tree.Left.String(), tree.Right.String())
}

func (tree *Tree) CountValues() int {
	if tree.IsValue() {
		return 1
	}
	return tree.Left.CountValues() + tree.Right.CountValues()
}

func (tree *Tree) Depth() int {
	depth := 0
	t := tree
	for t.Parent != nil {
		depth += 1
		t = t.Parent
	}
	return depth
}

func (tree *Tree) IsValuePair() bool {
	return tree.Left.IsValue() && tree.Right.IsValue()
}

func (tree *Tree) IsValue() bool {
	return tree.Left == nil && tree.Right == nil
}

func (tree *Tree) Split() {
	if !tree.IsValue() {
		return
	}
	if tree.Value < 10 {
		return
	}
	tree.Left = &Tree{Parent: tree, Value: tree.Value / 2}
	// the addition on the end is to ensure it is Value / 2 if Value is even and (Value / 2) + 1 of odd
	tree.Right = &Tree{Parent: tree, Value: (tree.Value / 2) + (tree.Value - 2*(tree.Value/2))}
}

func (tree *Tree) Explode(values []*Tree) func() {
	if tree.IsValue() {
		return func() {}
	}

	return func() {
		if left := getLeft(values, tree.Left); left != nil {
			left.Value += tree.Left.Value
		}
		if right := getRight(values, tree.Right); right != nil {
			right.Value += tree.Right.Value
		}
		tree.Left, tree.Right, tree.Value = nil, nil, 0
	}
}

func getLeft(values []*Tree, tree *Tree) *Tree {
	var left *Tree
	for i, v := range values {
		if v == tree && i > 0 {
			left = values[i-1]
		}
	}
	return left
}

func getRight(values []*Tree, tree *Tree) *Tree {
	var right *Tree
	for i, v := range values {
		if v == tree && i < len(values)-1 {
			right = values[i+1]
		}
	}
	return right
}
