package main

import (
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
		panic("need year and day arguments")
	}
	year, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	day, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	// create client with token
	bytes, err := os.ReadFile(tokenFile)
	if err != nil {
		panic(err)
	}
	token := string(bytes)

	client := puzzle.NewClient(token)
	p, err := client.Get(year, day)
	if err != nil {
		panic(err)
	}

	// make folders
	err = os.MkdirAll(folder(year, day), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// README.md
	readme, err := os.Create(fmt.Sprintf("%s/%s", folder(year, day), readmeFile))
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(readme, p.Readme)

	// input.txt
	input, err := os.Create(fmt.Sprintf("%s/%s", folder(year, day), inputFile))
	if err != nil {
		panic(err)
	}
	fmt.Fprint(input, p.Input)

	// test.txt
	test, err := os.Create(fmt.Sprintf("%s/%s", folder(year, day), testFile))
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(test, p.TestInput)

	// main.go
	solutionPath := fmt.Sprintf("%s/%s", folder(year, day), solutionFile)
	// only create main.go if it doesn't already exist
    if _, err := os.Stat(solutionPath); os.IsNotExist(err) {
		main, err := os.Create(solutionPath)
		if err != nil {
			panic(err)
		}
		solution, err := puzzle.InitialSolutionFile(p)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(main, solution)
    }
}

func folder(year, day int) string {
	if day < 10 {
		return fmt.Sprintf("%d/0%d", year, day)
	}
	return fmt.Sprintf("%d/%d", year, day)
}
