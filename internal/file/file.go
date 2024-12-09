// Package file contains utility methods for processing files
package file

import (
	"bufio"
	"os"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/geometry/plane"
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

func ReadCsvInts(path string, separator ...string) ([]int, error) {
	lines, err := ReadLines(path)
	if err != nil {
		return nil, err
	}

	if len(lines) < 1 {
		return []int{}, nil
	}

	return csv.ParseInts(lines[0], separator...)
}

func ReadAllCsvInts(path string, separator ...string) ([][]int, error) {
	lines, err := ReadLines(path)
	if err != nil {
		return nil, err
	}

	var nums [][]int
	for _, line := range lines {
		row, err := csv.ParseInts(line, separator...)
		if err != nil {
			return nil, err
		}
		nums = append(nums, row)
	}

	return nums, nil
}

func ReadVectors(path string, ch byte) ([]plane.Vector, plane.Vector, error) {
	lines, err := ReadLines(path)
	if err != nil {
		return nil, plane.Vector{}, err
	}

	var (
		vectors []plane.Vector
		limit   plane.Vector
	)
	for y := range lines {
		limit.Y = max(limit.Y, y)
		for x := range lines[y] {
			limit.X = max(limit.X, x)
			if lines[y][x] == ch {
				vectors = append(vectors, plane.Vector{X: x, Y: y})
			}
		}
	}
	return vectors, limit, nil
}

func ReadVectorsFunc(path string, f func(byte) bool) (map[byte][]plane.Vector, plane.Vector, error) {
	lines, err := ReadLines(path)
	if err != nil {
		return nil, plane.Vector{}, err
	}

	vectors := make(map[byte][]plane.Vector)
	var limit plane.Vector
	for y := range lines {
		limit.Y = max(limit.Y, y)
		for x := range lines[y] {
			limit.X = max(limit.X, x)
			if f(lines[y][x]) {
				vectors[lines[y][x]] = append(vectors[lines[y][x]], plane.Vector{X: x, Y: y})
			}
		}
	}
	return vectors, limit, nil
}
