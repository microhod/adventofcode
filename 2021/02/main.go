package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	COURSE_FILE = "course.txt"
	TEST_FILE   = "test-course.txt"
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
	commands, err := readCommands(COURSE_FILE)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	fmt.Println("--- Part 1 ---")

	position := Position{Horizontal: 0, Depth: 0}
	position = followCommands(commands, position)

	fmt.Printf("Position: %+v\n", position)
	fmt.Printf("Horizontal * Depth: %d\n", position.Horizontal * position.Depth)

	fmt.Println()

	fmt.Println("--- Part 2 ---")

	position = Position{Horizontal: 0, Depth: 0}
	position = followCommandsWithAim(commands, position)

	fmt.Printf("Position: %+v\n", position)
	fmt.Printf("Horizontal * Depth: %d\n", position.Horizontal * position.Depth)

	fmt.Println()
}

func followCommands(commands []Command, initial Position) Position {
	position := initial
	for _, command := range(commands) {
		switch command.Direction {
		case "forward":
			position.Horizontal += command.Units
			break
		case "up":
			position.Depth -= command.Units
			break
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
			break
		case "up":
			aim -= command.Units
			break
		case "down":
			aim += command.Units
		}
	}

	return position
}

func readCommands(path string) ([]Command, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	commands := []Command{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := scanner.Text()

		parts := strings.SplitN(command, " ", 2)
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
