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
	puzzle.NewSolution("CathodeRay Tube", part1, part2).Run()
}

func part1() error {
	instructions, err := parse(InputFile)
	if err != nil {
		return err
	}

	cpu := &CPU{RegisterX: 1}
	cpu.LoadInstructions(instructions)

	strengths := map[int]int{}

	for i := 0; i < 220; i++ {
		cycle := i + 1
		strengths[cycle] = cpu.RegisterX * cycle

		cpu.Tick()
	}

	sum := strengths[20] + strengths[60] + strengths[100] + strengths[140] + strengths[180] + strengths[220]
	fmt.Println(sum)
	return nil
}

func part2() error {
	instructions, err := parse(InputFile)
	if err != nil {
		return err
	}

	cpu := &CPU{RegisterX: 1}
	cpu.LoadInstructions(instructions)

	crt := NewCRT(40, 6)

	for i := 0; i < 240; i++ {
		spriteMiddle := cpu.RegisterX
		crtXPosition := i % crt.Width

		spriteVisible := false
		if crtXPosition >= spriteMiddle-1 && crtXPosition <= spriteMiddle+1 {
			spriteVisible = true
		}

		crt.WritePixel(spriteVisible)
		cpu.Tick()
	}

	fmt.Println(crt.Render())
	return nil
}

func parse(path string) ([]Instruction, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var instructions []Instruction

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		switch fields[0] {
		case "noop":
			instructions = append(instructions, Instruction{})
		case "addx":
			add, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}

			instructions = append(instructions, Instruction{})
			instructions = append(instructions, Instruction{AddX: add})
		}
	}
	return instructions, nil
}

type CPU struct {
	RegisterX    int
	Instructions []Instruction
}

func (cpu *CPU) LoadInstructions(instructions []Instruction) {
	cpu.Instructions = append(cpu.Instructions, instructions...)
}

func (cpu *CPU) Tick() {
	instruction := cpu.Instructions[0]
	cpu.Instructions = cpu.Instructions[1:]

	cpu.RegisterX += instruction.AddX
}

type Instruction struct {
	AddX int
}

type CRT struct {
	Width  int
	Screen [][]string
	Index  int
}

func NewCRT(width, height int) *CRT {
	crt := &CRT{Width: width}
	for i := 0; i < height; i++ {
		crt.Screen = append(crt.Screen, make([]string, width))
	}

	return crt
}

func (crt *CRT) WritePixel(on bool) {
	pixel := " "
	if on {
		pixel = "#"
	}

	row := crt.Index / crt.Width
	col := crt.Index % crt.Width

	crt.Screen[row][col] = pixel
	crt.Index += 1
}

func (crt *CRT) Render() string {
	var lines []string

	for _, row := range crt.Screen {
		lines = append(lines, strings.Join(row, ""))
	}

	return strings.Join(lines, "\n")
}
