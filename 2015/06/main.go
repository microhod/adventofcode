package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Probably a Fire Hazard", part1, part2).Run()
}

func part1() error {
	instructions, err := parseInstructions(InputFile)
	if err != nil {
		return err
	}

	lights := NewLightGrid()
	for _, i := range instructions {
		lights.Apply(i)
	}

	fmt.Printf("number of lights on: %d\n", lights.CountOn())
	return nil
}

func part2() error {
	instructions, err := parseInstructions(InputFile)
	if err != nil {
		return err
	}

	lights := NewBrightnessLightGrid()
	for _, i := range instructions {
		lights.Apply(i)
	}

	fmt.Printf("total brightness: %d\n", lights.Brightness())
	return nil
}

func parseInstructions(path string) ([]Instruction, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var instructions []Instruction
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		line = strings.ReplaceAll(line, "turn on ", "0,")
		line = strings.ReplaceAll(line, "turn off ", "1,")
		line = strings.ReplaceAll(line, "toggle ", "2,")

		line = strings.ReplaceAll(line, " through ", ",")

		nums, err := csv.ParseInts(line)
		if err != nil {
			return nil, err
		}
		if len(nums) != 5 {
			return nil, fmt.Errorf("invalid csv length: %d", len(nums))
		}

		instructions = append(instructions, Instruction{
			Command: Command(nums[0]),
			From: Vec(nums[1], nums[2]),
			To: Vec(nums[3], nums[4]),
		})
	}

	return instructions, nil
}

type Vector struct {
	X, Y int
}

func Vec(x, y int) Vector {
	return Vector{X: x, Y: y}
}

type Command int

const (
	TurnOn Command = iota
	TurnOff
	Toggle
)

type Instruction struct {
	Command  Command
	From, To Vector
}

type LightGrid struct {
	lights map[Vector]bool
}

func NewLightGrid() LightGrid {
	return LightGrid{lights: map[Vector]bool{}}
}

func (grid LightGrid) Apply(i Instruction) error {
	for x := i.From.X; x <= i.To.X; x++ {
		for y := i.From.Y; y <= i.To.Y; y++ {
			v := Vec(x, y)

			switch i.Command {
			case TurnOn:
				grid.lights[v] = true
			case TurnOff:
				grid.lights[v] = false
			case Toggle:
				grid.lights[v] = !grid.lights[v]
			default:
				return fmt.Errorf("invalid command: %d", i.Command)
			}
		}
	}
	return nil
}

func (grid LightGrid) CountOn() int {
	var count int
	for _, on := range grid.lights {
		if on {
			count += 1
		}
	}
	return count
}

type BrightnessLightGrid struct {
	lights map[Vector]int
}

func NewBrightnessLightGrid() BrightnessLightGrid {
	return BrightnessLightGrid{lights: map[Vector]int{}}
}

func (grid BrightnessLightGrid) Apply(i Instruction) error {
	for x := i.From.X; x <= i.To.X; x++ {
		for y := i.From.Y; y <= i.To.Y; y++ {
			v := Vec(x, y)

			switch i.Command {
			case TurnOn:
				grid.lights[v] += 1
			case TurnOff:
				grid.lights[v] -= 1
				if grid.lights[v] < 0 {
					grid.lights[v] = 0
				}
			case Toggle:
				grid.lights[v] += 2
			default:
				return fmt.Errorf("invalid command: %d", i.Command)
			}
		}
	}
	return nil
}

func (grid BrightnessLightGrid) Brightness() int {
	var total int
	for _, brightness := range grid.lights {
		total += brightness
	}
	return total
}
