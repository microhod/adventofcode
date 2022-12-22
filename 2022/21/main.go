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
	puzzle.NewSolution("Monkey Math", part1, part2).Run()
}

func part1() error {
	monkeys, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Println(monkeys["root"].ShoutNumber())
	return nil
}

func part2() error {
	monkeys, err := parse(InputFile)
	if err != nil {
		return err
	}

	// change root so that its dependencies minus to get zero (i.e. are equal)
	monkeys["root"].operation = "-"
	humanNumber, err := GetRequiredHumanNumber(monkeys["root"], 0)
	if err != nil {
		return err
	}
	fmt.Println(humanNumber)

	return nil
}

func parse(path string) (map[string]*Monkey, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	monkeys := map[string]*Monkey{}
	for _, line := range lines {
		if line == "" {
			continue
		}

		name := line[:4]
		monkeys[name] = &Monkey{name: name}
	}

	for _, line := range lines {
		if line == "" {
			continue
		}

		name := line[:4]
		if number, err := strconv.Atoi(line[6:]); err == nil {
			monkeys[name].number = int(number)
			continue
		}

		equation := strings.Fields(line[6:])
		monkeys[name].operation = equation[1]
		monkeys[name].dependencies = []*Monkey{
			monkeys[equation[0]],
			monkeys[equation[2]],
		}
	}
	return monkeys, nil
}

const Human = "humn"

type Monkey struct {
	name   string
	number int

	dependencies []*Monkey
	operation    string
}

func (m *Monkey) ShoutNumber() int {
	if len(m.dependencies) == 0 {
		return m.number
	}

	var nums []int
	for _, monkey := range m.dependencies {
		nums = append(nums, monkey.ShoutNumber())
	}

	return Operations[m.operation](nums[0], nums[1])
}

func (m *Monkey) HasDependency(dependency string) bool {
	if m.name == dependency {
		return true
	}
	for _, monkey := range m.dependencies {
		if monkey.HasDependency(dependency) {
			return true
		}
	}
	return false
}

var Operations = map[string]func(int, int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

func GetRequiredHumanNumber(current *Monkey, target int) (int, error) {
	if current.name == Human {
		return target, nil
	}

	leftContainsHuman := current.dependencies[0].HasDependency(Human)
	rightContainsHuman := current.dependencies[1].HasDependency(Human)
	if leftContainsHuman && rightContainsHuman {
		return 0, fmt.Errorf("both paths depend on the human number")
	}
	if !(leftContainsHuman || rightContainsHuman) {
		return 0, fmt.Errorf("neither path depends on the human number")
	}

	humanIndex := 0
	if rightContainsHuman {
		humanIndex = 1
	}
	humanPath := current.dependencies[humanIndex]
	monkeyPath := current.dependencies[(humanIndex+1)%2]

	// do reverse operation to target
	target = applyReverseOperation(current, monkeyPath.ShoutNumber(), humanIndex, target)

	return GetRequiredHumanNumber(humanPath, target)
}

func applyReverseOperation(monkey *Monkey, monkeyNumber, humanIndex, result int) int {
	switch monkey.operation {
	case "+":
		// monkey + human = result
		// human = result - monkey
		return Operations["-"](result, monkeyNumber)
	case "-":
		if humanIndex == 1 {
			// monkey - human = result
			// human = monkey - result
			return Operations["-"](monkeyNumber, result)
		}
		// human - monkey = result
		// human = result + monkey
		return Operations["+"](result, monkeyNumber)
	case "*":
		// monkey * human = result
		// human = result / monkey
		return Operations["/"](result, monkeyNumber)
	case "/":
		if humanIndex == 1 {
			// monkey / human = result
			// human = monkey / result
			return Operations["/"](monkeyNumber, result)
		}
		// human / monkey = result
		// human = result * monkey
		return Operations["*"](result, monkeyNumber)
	default:
		panic("unsupported operation")
	}

}
