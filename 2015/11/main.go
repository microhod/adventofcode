package main

import (
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
	puzzle.NewSolution("Corporate Policy", part1, part2).Run()
}

var valid = func(p Password) bool {
	return !(contains(p, int(rune('i') - rune('a'))) || 
		contains(p, int(rune('o') - rune('a'))) || 
		contains(p, int(rune('i') - rune('l')))) && 
		len(pairs(p)) >= 2 &&
		hasStraight(p, 3)
}
var password Password

func part1() error {
	var err error

	password, err = parse(InputFile)
	if err != nil {
		return err
	}

	for !valid(password) {
		password.increment()
	}

	fmt.Println(password.String())

	return nil
}

func part2() error {
	// increment to ensure it's not the same as the previous valid password
	password.increment()

	for !valid(password) {
		password.increment()
	}

	fmt.Println(password.String())
	return nil
}

func parse(path string) (Password, error) {
	b, err := file.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	var password Password
	for _, r := range strings.TrimSpace(string(b)) {
		password = append(password, int(r - rune('a')))
	}

	return password, nil
}

type Password []int

func (p Password) String() string {
	var str string
	for _, num := range p {
		str += string(rune(num + int(rune('a'))))
	}
	return str
}

func (p Password) increment() {
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = (p[i] + 1) % 26

		if p[i] != 0 {
			return
		}
	}
}

func contains(nums []int, num int) bool {
	for _, n := range nums {
		if n == num {
			return true
		}
	}
	return false
}

func pairs(nums []int) []int {
	var p []int
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] == nums[i+1] {
			p = append(p, nums[i])
			i++
		}
	}
	return p
}

func hasStraight(nums []int, length int) bool {
	for i := 0; i <= len(nums)-length; i++ {
		if isStraight(nums[i : i+length]) {
			return true
		}
	}
	return false
}

func isStraight(nums []int) bool {
	for i := 0; i < len(nums)-1; i++ {
		if nums[i]+1 != nums[i+1] {
			return false
		}
	}
	return true
}
