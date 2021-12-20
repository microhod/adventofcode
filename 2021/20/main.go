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
	puzzle.NewSolution("Trench Map", part1, part2).Run()
}

func part1() error {
	algorithm, image, err := readInput(InputFile)
	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		image = image.Enhance(algorithm, (i*algorithm[0])%2)
	}

	fmt.Printf("num pixels lit after 2 enhancements: %d\n", image.CountLit())

	return nil
}

func part2() error {
	algorithm, image, err := readInput(InputFile)
	if err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		image = image.Enhance(algorithm, (i*algorithm[0])%2)
	}

	fmt.Printf("num pixels lit after 2 enhancements: %d\n", image.CountLit())

	return nil
}

func readInput(path string) ([]int, Image, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, nil, err
	}

	algorithm := readInts(lines[0])
	image := Image{}
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		nums := readInts(line)
		image = append(image, nums)
	}

	return algorithm, image, nil
}

func readInts(str string) []int {
	bools := []int{}
	for _, char := range str {
		if char == rune('#') {
			bools = append(bools, 1)
		} else {
			bools = append(bools, 0)
		}
	}
	return bools
}

type Image [][]int

func (img Image) CountLit() int {
	count := 0
	for _, row := range img {
		for _, n := range row {
			count += n
		}
	}
	return count
}

func (img Image) AddPadding(value int) Image {
	image := [][]int{}

	padding := []int{}
	for i := 0; i < len(img)+4; i++ {
		padding = append(padding, value)
	}
	image = append(image, padding, padding)

	for _, r := range img {
		row := []int{value, value}
		row = append(row, r...)
		row = append(row, value, value)
		image = append(image, row)
	}

	return append(image, padding, padding)
}

func (img Image) Enhance(algorithm []int, void int) Image {
	image := img.AddPadding(void)
	enhanced := image.Copy()

	for i := 0; i < len(image); i++ {
		for j := 0; j < len(image[i]); j++ {
			nhood := image.Neighbourhood(i, j, void)
			str := ""
			for _, n := range nhood {
				str += fmt.Sprint(n)
			}
			index, err := strconv.ParseInt(str, 2, 64)
			if err != nil {
				panic(err)
			}

			enhanced[i][j] = algorithm[index]
		}
	}

	return enhanced
}

func (img Image) Neighbourhood(i, j, void int) []int {
	neighbourhood := []int{}
	for a := -1; a <= 1; a++ {
		for b := -1; b <= 1; b++ {
			if i+a < 0 || i+a > len(img)-1 || j+b < 0 || j+b > len(img[0])-1 {
				neighbourhood = append(neighbourhood, void)
				continue
			}
			neighbourhood = append(neighbourhood, img[i+a][j+b])
		}
	}
	return neighbourhood
}

func (img Image) Copy() Image {
	copy := Image{}
	for _, row := range img {
		c := append([]int{}, row...)
		copy = append(copy, c)
	}
	return copy
}

func (img Image) String() string {
	lines := []string{}
	for _, row := range img {
		line := ""
		for _, n := range row {
			char := "."
			if n == 1 {
				char = "#"
			}
			line += char
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
