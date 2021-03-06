package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	CourseFile = "course.txt"
	TestFile   = "test-course.txt"
)

type Command struct {
	Direction string
	Units     int
}

type Position struct {
	Horizontal int
	Depth      int
}

func main() {
	puzzle.NewSolution("Dive!", part1, part2).Run()
}

func part1() error {
	commands, err := readCommands(CourseFile)
	if err != nil {
		return err
	}

	position := Position{Horizontal: 0, Depth: 0}
	position = followCommands(commands, position)

	fmt.Printf("Position: %+v\n", position)
	fmt.Printf("Horizontal * Depth: %d\n", position.Horizontal * position.Depth)

	return nil
}

func part2() error {
	commands, err := readCommands(CourseFile)
	if err != nil {
		return err
	}

	position := Position{Horizontal: 0, Depth: 0}
	position = followCommandsWithAim(commands, position)

	fmt.Printf("Position: %+v\n", position)
	fmt.Printf("Horizontal * Depth: %d\n", position.Horizontal * position.Depth)

	return nil
}

func followCommands(commands []Command, initial Position) Position {
	position := initial
	for _, command := range(commands) {
		switch command.Direction {
		case "forward":
			position.Horizontal += command.Units
		case "up":
			position.Depth -= command.Units
		case "down":
			position.Depth += command.Units
		}
	}

	return position
}

func followCommandsWithAim(commands []Command, initial Position) Position {
	position := initial
	aim := 0
	for _, command := range(commands) {
		switch command.Direction {
		case "forward":
			position.Horizontal += command.Units
			position.Depth += aim * command.Units
		case "up":
			aim -= command.Units
		case "down":
			aim += command.Units
		}
	}

	return position
}

func readCommands(path string) ([]Command, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	commands := []Command{}
	for _, line := range(lines) {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected command to split into two space separated parts, but got %d parts", len(parts))
		}

		direction, unitsStr := parts[0], parts[1]

		units, err := strconv.Atoi(unitsStr)
		if err != nil {
			return nil, err
		}

		commands = append(commands, Command{
			Direction: direction,
			Units: units,
		})
	}

	return commands, nil
}
