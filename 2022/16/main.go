package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/graph"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Proboscidea Volcanium", part1, part2).Run()
}

// part1 takes ~35 seconds
func part1() error {
	tunnels, err := parse(InputFile)
	if err != nil {
		return err
	}
	start := Valve{"AA", 0}

	// remove zeros (apart from start)
	tunnels.RemoveZeros(start.Label)
	// connect all valves
	tunnels.ConnectAll()

	max := tunnels.FindMaximumFlow(start, 30)
	fmt.Println(max)
	return nil
}

// part2 takes ~ 15 mins (not great, but it works!)
func part2() error {
	tunnels, err := parse(InputFile)
	if err != nil {
		return err
	}
	start := Valve{"AA", 0}

	// remove zeros (apart from start)
	tunnels.RemoveZeros(start.Label)
	// connect all valves
	tunnels.ConnectAll()

	var valvesToChoose []Valve
	for valve := range tunnels.Graph {
		if valve == start {
			continue
		}
		valvesToChoose = append(valvesToChoose, valve)
	}
	
	// create jobs for every possible pair of valves
	choices := GetAllChoices(valvesToChoose)
	jobs := make(chan set.Set[Valve], len(choices))
	for _, choice := range choices {
		jobs <- choice
	}
	close(jobs)

	// run jobs
	workers := 4
	maxPressure := 0
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()

			for choice := range jobs {
				me := tunnels.Copy()
				elephant := tunnels.Copy()

				for _, valve := range valvesToChoose {
					if choice.Contains(valve) {
						me.RemoveValve(valve)
					} else {
						elephant.RemoveValve(valve)
					}
				}

				maxMe := me.FindMaximumFlow(start, 26)
				maxElephant := elephant.FindMaximumFlow(start, 26)

				mu.Lock()
				if maxMe+maxElephant > maxPressure {
					maxPressure = maxMe + maxElephant
				}
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	fmt.Println(maxPressure)
	return nil
}

func parse(path string) (Tunnels, error) {
	input, err := file.Read(path)
	if err != nil {
		return Tunnels{}, err
	}

	input = strings.ReplaceAll(input, "Valve ", "")
	input = strings.ReplaceAll(input, " has flow rate=", ", ")
	input = strings.ReplaceAll(input, "; tunnels lead to valves ", ", ")
	input = strings.ReplaceAll(input, "; tunnel leads to valve ", ", ")

	valves := map[string]Valve{}
	neighbours := map[string][]string{}

	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ", ")

		valve := parts[0]
		flowRate, err := strconv.Atoi(parts[1])
		if err != nil {
			return Tunnels{}, err
		}
		valves[valve] = Valve{
			Label:    valve,
			FlowRate: flowRate,
		}
		neighbours[valve] = parts[2:]
	}

	tunnels := Tunnels{Graph: graph.NewGraph[Valve]()}
	for _, valve := range valves {
		tunnels.Graph[valve] = map[Valve]int{}

		for _, neighbour := range neighbours[valve.Label] {
			tunnels.Graph[valve][valves[neighbour]] = 1
		}
	}

	return tunnels, nil
}

type Valve struct {
	Label    string
	FlowRate int
}

func (v Valve) String() string {
	return v.Label + "_" + fmt.Sprint(v.FlowRate)
}

type Flow struct {
	tunnels  Tunnels
	current  Valve
	duration int
	pressure int
}

type Tunnels struct {
	graph.Graph[Valve]
}

func (tunnels Tunnels) FindMaximumFlow(start Valve, timeLimit int) int {
	flows := []Flow{
		{
			tunnels:  tunnels.Copy(),
			current:  start,
			duration: 0,
			pressure: 0,
		},
	}
	var maxPressure int

	for len(flows) > 0 {
		// fmt.Println()
		// for _, flow := range flows {
		// 	fmt.Println(flow.current, flow.duration, flow.pressure)
		// }

		flow := flows[0]
		flows = flows[1:]

		if flow.duration >= timeLimit {
			if flow.pressure > maxPressure {
				maxPressure = flow.pressure
			}
			continue
		}
		if len(flow.tunnels.Graph) == 1 {
			flow.duration += 1
			flow.pressure += (timeLimit - flow.duration) * flow.current.FlowRate

			if flow.pressure > maxPressure {
				maxPressure = flow.pressure
			}
			continue
		}

		tunnelsWithoutCurrent := flow.tunnels.Copy()
		tunnelsWithoutCurrent.RemoveValve(flow.current)

		// unlock current valve
		if flow.current.FlowRate > 0 {
			flow.duration += 1
		}
		flow.pressure += (timeLimit - flow.duration) * flow.current.FlowRate

		// try each neighbour
		for neighbour, weight := range flow.tunnels.Graph[flow.current] {
			// move to neighbour
			flows = append(flows, Flow{
				tunnels:  tunnelsWithoutCurrent.Copy(),
				current:  neighbour,
				duration: flow.duration + weight,
				pressure: flow.pressure,
			})
		}
	}
	return maxPressure
}

func (tunnels Tunnels) ConnectAll() {
	var valves []Valve
	for v := range tunnels.Graph {
		valves = append(valves, v)
	}

	for i := 1; i < len(valves); i++ {
		for j := 0; j < i; j++ {
			if _, exists := tunnels.Graph[valves[i]][valves[j]]; !exists {
				tunnels.Connect(valves[i], valves[j])
			}
		}
	}
}

func (tunnels Tunnels) Connect(start, target Valve) {
	_, cost, exists := tunnels.DijkstraShortestPath(start, target)
	if !exists {
		panic("no path exists")
	}

	tunnels.Graph[start][target] = cost
	tunnels.Graph[target][start] = cost
}

func (tunnels Tunnels) RemoveZeros(startLabel string) {
	var zeros []Valve
	for valve := range tunnels.Graph {
		if valve.FlowRate == 0 && valve.Label != startLabel {
			zeros = append(zeros, valve)
		}
	}

	for _, zero := range zeros {
		tunnels.RemoveValve(zero)
	}
}

func (tunnels Tunnels) RemoveValve(valve Valve) {
	var neighbours []Valve
	for n := range tunnels.Graph[valve] {
		neighbours = append(neighbours, n)

		// remove valve from neighbours
		delete(tunnels.Graph[n], valve)
	}

	// tie neighbours together
	for i := 1; i < len(neighbours); i++ {
		for j := 0; j < i; j++ {
			iToJ := tunnels.Graph[valve][neighbours[i]] + tunnels.Graph[valve][neighbours[j]]

			current, exists := tunnels.Graph[neighbours[i]][neighbours[j]]
			if !exists || iToJ < current {
				tunnels.Graph[neighbours[i]][neighbours[j]] = iToJ
				tunnels.Graph[neighbours[j]][neighbours[i]] = iToJ
			}
		}
	}

	// delete valve from tunnels completely
	delete(tunnels.Graph, valve)
}

func (tunnels Tunnels) Copy() Tunnels {
	copy := Tunnels{Graph: graph.NewGraph[Valve]()}
	for from := range tunnels.Graph {
		copy.Graph[from] = map[Valve]int{}

		for to, cost := range tunnels.Graph[from] {
			copy.Graph[from][to] = cost
		}
	}
	return copy
}

func GetAllChoices(valves []Valve) []set.Set[Valve] {
	var choices []set.Set[Valve]
	for length := 1; length <= len(valves)/2; length++ {
		choices = append(choices, GetChoices(valves, length)...)
	}

	return choices
}

func GetChoices(valves []Valve, length int) []set.Set[Valve] {
	var choices []set.Set[Valve]

	if length == 1 {
		for _, valve := range valves {
			choices = append(choices, set.NewSet(valve))
		}
		return choices
	}

	for i := 0; i < len(valves); i++ {
		for _, choice := range GetChoices(valves[i+1:], length-1) {
			choice.Add(valves[i])
			choices = append(choices, choice)
		}
	}

	return choices
}
