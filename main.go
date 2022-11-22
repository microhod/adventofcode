package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	tokenFile    = ".token"
	readmeFile   = "README.md"
	inputFile    = "input.txt"
	testFile     = "test.txt"
	solutionFile = "main.go"
)

// get the puzzle files for the year & day specified
func main() {
	// parse arguments
	if len(os.Args) < 3 {
		fail(fmt.Errorf("need year and day arguments"))
	}
	year, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fail(err)
	}
	day, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fail(err)
	}

	// create client with token
	bytes, err := os.ReadFile(tokenFile)
	if err != nil {
		fail(err)
	}
	token := string(bytes)

	client := puzzle.NewClient(token)
	p, err := client.Get(year, day)
	if err != nil {
		fail(err)
	}

	// make folders
	err = os.MkdirAll(folder(year, day), os.ModePerm)
	if err != nil {
		fail(err)
	}

	// README.md
	readme, err := os.Create(fmt.Sprintf("%s/%s", folder(year, day), readmeFile))
	if err != nil {
		fail(err)
	}
	fmt.Fprintln(readme, p.Readme)

	// input.txt
	input, err := os.Create(fmt.Sprintf("%s/%s", folder(year, day), inputFile))
	if err != nil {
		fail(err)
	}
	fmt.Fprint(input, p.Input)

	// test.txt
	testFilePath := fmt.Sprintf("%s/%s", folder(year, day), testFile)
	// only create test.txt if it doesn't already exist
	if !exists(testFilePath) {
		test, err := os.Create(testFilePath)
		if err != nil {
			fail(err)
		}
		fmt.Fprintln(test, p.TestInput)
	}

	// main.go
	solutionPath := fmt.Sprintf("%s/%s", folder(year, day), solutionFile)
	// only create main.go if it doesn't already exist
	if !exists(solutionPath) {
		main, err := os.Create(solutionPath)
		if err != nil {
			fail(err)
		}
		solution, err := puzzle.InitialSolutionFile(p)
		if err != nil {
			fail(err)
		}
		fmt.Fprint(main, solution)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func fail(err error) {
	fmt.Printf("ERROR: %s\n", err)
	os.Exit(1)
}

func folder(year, day int) string {
	if day < 10 {
		return fmt.Sprintf("%d/0%d", year, day)
	}
	return fmt.Sprintf("%d/%d", year, day)
}
