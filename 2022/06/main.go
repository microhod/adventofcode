package main

import (
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
	"github.com/microhod/adventofcode/internal/set"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Tuning Trouble", part1, part2).Run()
}

func part1() error {
	stream, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Println(stream.PacketStart())
	return nil
}

func part2() error {
	stream, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Println(stream.MessageStart())
	return nil
}

func parse(path string) (DataStream, error) {
	data, err := file.ReadBytes(path)
	if err != nil {
		return "", err
	}

	return DataStream(strings.TrimSpace(string(data))), nil
}

type DataStream string

func (d DataStream) PacketStart() int {
	for i := 0; i <= len(d)-4; i++ {
		chars := set.NewSet([]byte(d[i:i+4])...)
		if len(chars) == 4 {
			return i+4
		}
	}
	return -1
}

func (d DataStream) MessageStart() int {
	for i := 0; i <= len(d)-14; i++ {
		chars := set.NewSet([]byte(d[i:i+14])...)
		if len(chars) == 14 {
			return i+14
		}
	}
	return -1
}
