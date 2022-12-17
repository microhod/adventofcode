// Package file contains utility methods for processing files
package file

import (
	"bufio"
	"os"

	"github.com/microhod/adventofcode/internal/encoding/csv"
)

func Read(path string) (string, error) {
	b, err := ReadBytes(path)
	return string(b), err
}

func ReadBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func ReadCsvInts(path string) ([]int, error) {
	lines, err := ReadLines(path)
	if err != nil {
		return nil, err
	}

	if len(lines) < 1 {
		return []int{}, nil
	}

	return csv.ParseInts(lines[0])
}
