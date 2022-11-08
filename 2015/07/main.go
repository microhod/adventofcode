package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Some Assembly Required", part1, part2).Run()
}

var wires = Wires{}
var operations Operations

func part1() error {
	var err error

	operations, err = parseOperations(InputFile)
	if err != nil {
		return err
	}

	err = wires.run(operations)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("wire[a] = %d\n", wires["a"])

	return nil
}

func part2() error {
	// reuse parsed operations & the wire state from part 1
	wires = Wires{
		"b": wires["a"],
	}

	err := wires.run(operations)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("wire[a] = %d\n", wires["a"])

	return nil
}

func parseOperations(path string) (Operations, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	operations := Operations{}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, " -> ")
		if len(parts) < 2 {
			return nil, fmt.Errorf("expected 'something -> something' but got '%s'", line)
		}
		target := parts[1]

		var operation Operation

		parts = strings.Fields(parts[0])
		switch len(parts) {
		case 1:
			operation.Gate = "SET"
			operation.Inputs = []string{parts[0]}
		case 2:
			operation.Gate = parts[0]
			operation.Inputs = []string{parts[1]}
		case 3:
			operation.Gate = parts[1]
			operation.Inputs = []string{parts[0], parts[2]}
		}

		if _, ok := operations[target]; ok {
			return nil, fmt.Errorf("duplicate target: '%s'", target)
		}
		operations[target] = operation
	}

	return operations, nil
}

type Operation struct {
	Inputs []string
	Gate   string
}

type Operations map[string]Operation

type Wires map[string]int

// run applies all operations concurrently
func (w Wires) run(operations Operations) error {
	mu := &sync.RWMutex{}
	eg := &errgroup.Group{}

	for t := range operations {
		target, op := t, operations[t]
		eg.Go(func() error {
			return w.apply(op, target, mu)
		})
	}

	return eg.Wait()
}

func (w Wires) apply(o Operation, target string, mu *sync.RWMutex) error {
	// wait for dependencies
	pending := append([]string{}, o.Inputs...)
	var inputs []int

	for len(pending) > 0 {
		input := pending[0]
		pending = pending[1:]

		// add to inputs if it's just a number
		num, err := strconv.Atoi(input)
		if err == nil {
			inputs = append(inputs, num)
			continue
		}

		// otherwise it's a wire
		// check if the wire has been set yet
		mu.RLock()
		num, ok := w[input]
		mu.RUnlock()
		if ok {
			inputs = append(inputs, num)
			continue
		}

		// add input back to pending and wait
		pending = append([]string{input}, pending...)
		time.Sleep(1 * time.Millisecond)
	}

	result, err := o.run(inputs)
	if err != nil {
		return err
	}

	mu.Lock()
	w[target] = result
	mu.Unlock()
	return nil
}

func (o Operation) run(inputs []int) (int, error) {
	switch o.Gate {
	case "SET":
		return inputs[0], nil
	case "NOT":
		return ^inputs[0], nil
	case "AND":
		return inputs[0] & inputs[1], nil
	case "OR":
		return inputs[0] | inputs[1], nil
	case "LSHIFT":
		return inputs[0] << inputs[1], nil
	case "RSHIFT":
		return inputs[0] >> inputs[1], nil
	default:
		return 0, fmt.Errorf("unsupported gate: '%s'", o.Gate)
	}
}
