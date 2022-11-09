package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Matchsticks", part1, part2).Run()
}

func part1() error {
	lines, err := readLines(InputFile)
	if err != nil {
		return err
	}

	fmt.Println(sum(lines) - sum(decodeLines(lines)))

	return nil
}

func part2() error {
	lines, err := readLines(InputFile)
	if err != nil {
		return err
	}

	fmt.Println(sum(encodeLines(lines)) - sum(lines))

	return nil
}

func readLines(path string) ([]string, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var filtered []string
	for _, line := range lines {
		if line == "" {
			continue
		}

		filtered = append(filtered, line)
	}

	return filtered, nil
}

func decodeLines(lines []string) []string {
	var decoded []string
	for _, line := range lines {
		decoded = append(decoded, decode(line))
	}
	return decoded
}

func decode(str string) string {
	// trim quotes
	str = str[1 : len(str)-1]
	str = strings.ReplaceAll(str, `\\`, `\`)
	str = strings.ReplaceAll(str, `\"`, `"`)

	decoded := str

	for i := 0; i < len(str)-3; i++ {
		s := str[i : i+4]
		if !strings.HasPrefix(s, `\x`) {
			continue
		}
		b, err := hex.DecodeString(s[2:])
		if err == nil {
			decoded = strings.ReplaceAll(decoded, s, string(b))
		}
	}

	return decoded
}

func encodeLines(lines []string) []string {
	var encoded []string
	for _, line := range lines {
		encoded = append(encoded, encode(line))
	}
	return encoded
}

func encode(str string) string {
	str = strings.ReplaceAll(str, `\`, `\\`)
	str = strings.ReplaceAll(str, `"`, `\"`)
	return `"` + str + `"`
}

func sum(strs []string) int {
	var total int
	for _, str := range strs {
		total += len(str)
	}
	return total
}
