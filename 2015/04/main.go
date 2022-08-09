package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("The Ideal Stocking Stuffer", part1, part2).Run()
}

func part1() error {
	key, err := readSecretKey(InputFile)
	if err != nil {
		return err
	}
	
	var num int
	var hash string

	for len(hash) < 5 || hash[:5] != "00000" {
		num += 1
		hash = md5Hex(key + fmt.Sprint(num))
	}

	fmt.Printf("AdventCoin number: %d\n", num)
	return nil
}

func part2() error {
	key, err := readSecretKey(InputFile)
	if err != nil {
		return err
	}
	
	var num int
	var hash string

	for len(hash) < 6 || hash[:6] != "000000" {
		num += 1
		hash = md5Hex(key + fmt.Sprint(num))
	}

	fmt.Printf("AdventCoin number: %d\n", num)
	return nil
}

func readSecretKey(path string) (string, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return "", err
	}

	if len(lines) < 1 {
		return "", fmt.Errorf("empty file")
	}
	return lines[0], nil
}

func md5Hex(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
