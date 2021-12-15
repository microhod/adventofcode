package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Chiton", part1, part2).Run()
}

func part1() error {
	ints, err := readInts(InputFile)
	if err != nil {
		return err
	}
	graph := readGraph(ints)

	source := graph[0]
	target := graph[len(graph)-1]

	min := graph.Dijkstra(source, target)

	fmt.Printf("min risk path: %d\n", min)

	return nil
}

func part2() error {
	ints, err := readInts(InputFile)
	if err != nil {
		return err
	}
	ints = fullMap(ints, 5)
	graph := readGraph(ints)

	source := graph[0]
	target := graph[len(graph)-1]

	fmt.Printf("this is unfortunately very slow, it takes over 30 mins 😢\n")
	min := graph.Dijkstra(source, target)

	fmt.Printf("min risk path: %d\n", min)

	return nil
}

func readInts(path string) ([][]int, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	nums := [][]int{}
	for _, line := range lines {
		row := []int{}
		for _, char := range line {
			risk, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, err
			}
			row = append(row, risk)
		}
		nums = append(nums, row)
	}

	return nums, nil
}

func fullMap(tile [][]int, multiplier int) [][]int {
	fullHorizontal := [][]int{}
	for _, row := range tile {
		fullRow := []int{}
		for i := 0; i < multiplier; i++ {
			fullRow = append(fullRow, sum(row, i)...)
		}
		fullHorizontal = append(fullHorizontal, fullRow)
	}

	full := [][]int{}
	for i := 0; i < multiplier; i++ {
		for j := 0; j < len(fullHorizontal); j++ {
			full = append(full, sum(fullHorizontal[j], i))
		}
	}

	return full
}

func sum(nums []int, a int) []int {
	result := []int{}
	for _, n := range nums {
		result = append(result, (n-1+a) % 9 + 1)
	}
	return result
}

func readGraph(nums [][]int) Graph {
	nodes := [][]*Node{}
	for _, line := range nums {
		row := []*Node{}
		for _, risk := range line {
			row = append(row, &Node{Risk: risk})
		}
		nodes = append(nodes, row)
	}

	for i, row := range nodes {
		for j, node := range row {
			if i > 0 {
				node.Neighbours = append(node.Neighbours, nodes[i-1][j])
			}
			if i < len(nodes)-1 {
				node.Neighbours = append(node.Neighbours, nodes[i+1][j])
			}
			if j > 0 {
				node.Neighbours = append(node.Neighbours, nodes[i][j-1])
			}
			if j < len(row)-1 {
				node.Neighbours = append(node.Neighbours, nodes[i][j+1])
			}
		}
	}

	graph := Graph{}
	for _, row := range nodes {
		graph = append(graph, row...)
	}

	return graph
}

type Graph []*Node

type Node struct {
	Risk       int
	Neighbours []*Node
}

type PartialPath struct {
	Path []*Node
	Next *Node
}

func (p *PartialPath) Risk() int {
	risk := 0
	for _, n := range p.Path {
		risk += n.Risk
	}
	return risk + p.Next.Risk
}

func (graph Graph) Dijkstra(source *Node, target *Node) int {
	queue := []*Node{}
	dist := map[*Node]int{}
    for _, v := range graph {
        dist[v] = math.MaxInt
        queue = append(queue, v)
	}
	dist[source] = 0

	for len(queue) > 0 {
		index := minIndex(queue, dist)

		u := queue[index]
		queue = append(queue[:index], queue[index+1:]...)

		if u == target {
			return dist[target]
		}

		unvisited := filter(u.Neighbours, queue)

		for _, v := range unvisited {
			alt := dist[u] + v.Risk
			if alt < dist[v] {
				dist[v] = alt
			}
		}
	}

	return dist[target]
}

func minIndex(nodes []*Node, dist map[*Node]int) int {
	min := 0
	for i := range nodes[1:] {
		if dist[nodes[i]] < dist[nodes[min]] {
			min = i
		}
	}
	return min
}

func filter(nodes []*Node, filter []*Node) []*Node {
	filtered := []*Node{}
	for _, n := range nodes {
		if contains(filter, n) {
			filtered = append(filtered, n)
		}
	}
	return filtered
}

func contains(nodes []*Node, node *Node) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}
