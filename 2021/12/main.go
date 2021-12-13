package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"

	Start = "start"
	End   = "end"
)

func main() {
	puzzle.NewSolution("Passage Pathing", part1, part2).Run()
}

func part1() error {
	_, start, err := readGraph(InputFile)
	if err != nil {
		return err
	}

	paths := start.TraversePaths(End, "")

	fmt.Printf("# paths: %d\n", len(paths))

	return nil
}

func part2() error {
	nodes, start, err := readGraph(InputFile)
	if err != nil {
		return err
	}

	// initial paths (with no double small nodes allowed)
	paths := start.TraversePaths(End, "")

	for _, label := range getSmallLabels(nodes) {
		// get paths with double pass on label
		doublePassPaths := start.TraversePaths(End, label)
		paths = append(paths, doublePassPaths...)
	}

	fmt.Printf("# paths: %d\n", len(paths))

	return nil
}

func readGraph(path string) ([]*Node, *Node, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, nil, err
	}

	nodes := []*Node{}
	var start *Node = nil

	getOrCreate := func(label string) *Node {
		node := getNode(nodes, label)
		if node == nil {
			node = &Node{Label: label}
			nodes = append(nodes, node)
		}
		if node.Label == Start {
			start = node
		}
		return node
	}
	for i, line := range lines {
		labels := strings.Split(line, "-")
		if len(labels) < 2 {
			return nil, nil, fmt.Errorf("line %d: expected '<node>-<other>' but got '%s'", i+1, line)
		}

		node0 := getOrCreate(labels[0])
		node1 := getOrCreate(labels[1])

		node0.Neighbours = append(node0.Neighbours, node1)
		node1.Neighbours = append(node1.Neighbours, node0)
	}

	return nodes, start, nil
}

type Node struct {
	Label      string
	Neighbours []*Node
}

type PartialPath struct {
	Path []*Node
	Next *Node
}

func (node *Node) String() string {
	return node.Label
}

func (node *Node) TraversePaths(end string, doubleVisit string) [][]*Node {
	paths := [][]*Node{}

	stack := []*PartialPath{}

	for _, n := range node.Neighbours {
		stack = append(stack, &PartialPath{
			Path: []*Node{node},
			Next: n,
		})
	}

	for len(stack) > 0 {
		partial := stack[0]
		stack = stack[1:]

		path := copyAndAppend(partial.Path, partial.Next)

		for _, n := range partial.Next.Neighbours {
			// add path if it is at the end
			if n.Label == end {
				full := copyAndAppend(path, n)

				if doubleVisit == "" || countLabel(full, doubleVisit) == 2 {
					paths = append(paths, full)
				}
				continue
			}
			// don't go back to the start
			if n.Label == Start {
				continue
			}
			// don't go back to double visit cave more than twice
			if doubleVisit == n.Label && countLabel(path, doubleVisit) == 2 {
				continue
			}
			// don't go back to small caves we don't want to double visit
			if doubleVisit != n.Label && strings.ToLower(n.Label) == n.Label && getNode(path, n.Label) != nil {
				continue
			}

			stack = append(stack, &PartialPath{
				Path: path,
				Next: n,
			})
		}
	}

	return paths
}

func (node *Node) Clone() string {
	label := fmt.Sprintf("%sclone", node.Label)
	clone := &Node{
		Label:      label,
		Neighbours: node.Neighbours,
	}
	for _, n := range node.Neighbours {
		n.Neighbours = append(n.Neighbours, clone)
	}
	return label
}

func copyAndAppend(old []*Node, nodes ...*Node) []*Node {
	copy := append([]*Node{}, old...)
	return append(copy, nodes...)
}

func getNode(nodes []*Node, label string) *Node {
	for _, n := range nodes {
		if n.Label == label {
			return n
		}
	}
	return nil
}

func countLabel(nodes []*Node, label string) int {
	count := 0
	for _, n := range nodes {
		if n.Label == label {
			count += 1
		}
	}
	return count
}

func getSmallLabels(nodes []*Node) []string {
	small := []string{}
	for _, n := range nodes {
		if strings.ToLower(n.Label) == n.Label && n.Label != Start && n.Label != End {
			small = append(small, n.Label)
		}
	}

	return small
}
