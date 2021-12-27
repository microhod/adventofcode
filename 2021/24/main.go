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
	puzzle.NewSolution("Arithmetic Logic Unit", part1, part2).Run()
}

func part1() error {
	monad, err := readMonad(InputFile)
	if err != nil {
		return err
	}

	modelNums := findAll(monad, findStartingMax(monad))
	fmt.Printf("maximum model number: %d\n", modelNums[len(modelNums)-1])

	return nil
}

func part2() error {
	monad, err := readMonad(InputFile)
	if err != nil {
		return err
	}

	modelNums := findAll(monad, findStartingMin(monad))
	fmt.Printf("minimum model number: %d\n", modelNums[0])

	return nil
}

// read input

func readMonad(path string) ([][]string, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	monad := [][]string{}
	commands := []string{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "inp") {
			monad = append(monad, commands)
			commands = []string{}
		}
		commands = append(commands, line)
	}

	monad = append(monad, commands)
	return monad[1:], nil
}

// algorithms

/*
All sections of the monad look like:
inp w
mul x 0
add x PREV (z)
mod x 26
div z DIV
add x ADDX
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y ADDY
mul y x
add z y

ADDY > 0
DIV = 26 if ADDX < 0
DIV = 1  if ADDX > 0

the sequence of execution can be seen as:

x = (PREV % 26) + ADDX
z = (PREV / DIV)

then check x = (PREV % 26) + ADDX == w
x = 0 if x == w
x = 1 if x != w

y = 1 if x == w
y = 26  if x != w

z =   (PREV / DIV) if x == w
z = 26(PREV / DIV) if x != w

y = 0        if x == w
y = w + ADDY if x != w

z =   (PREV / DIV)            if x == w
z = 26(PREV / DIV) + w + ADDY if x != w

so, we can know the value of z based on x and the digit (w)

let's compute the first few iterations on our input:

1:
ADDX = 11
ADDY = 16
DIV  = 1
PREV = 0

check => (PREV % 26) + ADDX = ADDX
      => x != w_1 as 11 > w_1 for all 1 <= w_1 <= 9
=> z = w_1 + 16 (<26)

2:
ADDX = 12
ADDY = 11
DIV  = 1
PREV = w_1 + 16

check => (PREV % 26) + ADDX = w_1 + 16 + 12
      => x != w_2
=> z = 26(w_1 + 16) + w_2 + 11

3:
ADDX = 13
ADDY = 12
DIV  = 1
PREV = 26(w_1 + 16) + w_2 + 11

=> z = 26(26(w_1 + 16) + w_2 + 11) + w_3 + 12

4:
ADDX = -5
ADDY = 12
DIV  = 26
PREV = 26(26(w_1 + 16) + w_2 + 11) + w_3 + 12

check => (PREV % 26) + ADDX = w_3 + 12 - 5 = w_3 + 7
      => x == w_4 if w_4 == w_3 + 7
	  => 1 <= w_3 <= 2 && 8 <= w_4 <= 9

=> z = 26(w_1 + 16) + w_2 + 11

So we can see that once we reach the first negative ADDX, it will supply bounds to
the character at its index, and the previous positive ADDX index.

This can then be used to reduce the size of the search for max / min model number.

Generally the bounds in the case of `w_neg = w_pos + ADDY_pos + ADDX_neg` are:

                      1 <= w_pos <= 9 - (ADDY_pos + ADDX_neg)
ADDY_pos + ADDX_neg + 1 <= w_neg <= 9

However, these are linked so we cannot simply take all possible pairs in this range.

For instance, the maximum case is found when taking the maximum value for `w_pos` (as this is the most significant digit):

	w_pos = 9 - (ADDY_pos + ADDX_neg)
	w_neg = 9

Similarly, the minimum takes the minimum for `w_pos`.
*/

func findStartingMax(monad [][]string) [14]int {
	num := [14]int{}
	addYStack := [][]int{}

	for index, commands := range monad {
		addX, _ := strconv.Atoi(strings.Split(commands[5], " ")[2])
		addY, _ := strconv.Atoi(strings.Split(commands[15], " ")[2])
		if addX > 0 {
			addYStack = append([][]int{{index, addY}}, addYStack...)
		} else {
			indexPos, addYPos := addYStack[0][0], addYStack[0][1]
			addYStack = addYStack[1:]

			if 0 <= addYPos+addX && 9 > addYPos+addX {
				num[indexPos] = 9 - (addYPos + addX)
				num[index] = 9
			}
		}
	}

	return num
}

func findStartingMin(monad [][]string) [14]int {
	num := [14]int{}
	addYStack := [][]int{}

	for index, commands := range monad {
		addX, _ := strconv.Atoi(strings.Split(commands[5], " ")[2])
		addY, _ := strconv.Atoi(strings.Split(commands[15], " ")[2])
		if addX > 0 {
			addYStack = append([][]int{{index, addY}}, addYStack...)
		} else {
			indexPos, addYPos := addYStack[0][0], addYStack[0][1]
			addYStack = addYStack[1:]

			if 0 <= addYPos+addX && 9 > addYPos+addX {
				num[indexPos] = 1
				num[index] = 1 + (addYPos + addX)
			}
		}
	}

	return num
}

func findAll(monad [][]string, start [14]int) []int {
	digits := start
	zeros := []int{}
	for i, n := range digits {
		if n == 0 {
			zeros = append(zeros, i)
		}
	}

	possibilities := getAllDigits(len(zeros))
	
	validNums := []int{}
	for _, p := range possibilities {
		for i, index := range zeros {
			digits[index] = p[i]
		}
		n := num(digits[:])
		valid, _ := validate(monad, n, 0)
		if valid {
			validNums = append(validNums, n)
		}
	}

	return validNums
}

func getAllDigits(size int) [][]int {
	all := [][]int{}
	for n := pow10(size-1); n < pow10(size); n++ {
		digits := digits(n)
		if contains(digits, 0) {
			continue
		}
		all = append(all, digits)
	}
	return all
}

func validate(monad [][]string, num int, z int) (bool, error) {
	digits := []int{}
	for _, ch := range strings.Split(fmt.Sprint(num), "") {
		d, _ := strconv.Atoi(ch)
		digits = append(digits, d)
	}

	validator := []string{}
	for i := len(monad) - len(digits); i < len(monad); i++ {
		validator = append(validator, monad[i]...)
	}

	alu := &ALU{
		Register: [4]int{0, 0, 0, z},
	}
	digit := &IntReader{buf: digits}

	if err := alu.Run(validator, digit); err != nil {
		return false, err
	}

	return alu.Register[3] == 0, nil
}

// ALU interpreter

type Reader interface {
	Read() (int, error)
}

type IntReader struct {
	buf []int
	pos int
}

func (r *IntReader) Read() (int, error) {
	if r.pos >= len(r.buf) {
		return 0, fmt.Errorf("tried to read out of buffer range: %d", r.pos)
	}
	r.pos += 1
	return r.buf[r.pos-1], nil
}

type ALU struct {
	Register [4]int
	input    Reader
}

func (alu *ALU) Run(program []string, input Reader) error {
	alu.input = input
	for _, command := range program {
		if err := alu.Execute(command); err != nil {
			return err
		}
	}
	return nil
}

func (alu *ALU) Execute(command string) error {
	reg := map[string]int{
		"w": 0,
		"x": 1,
		"y": 2,
		"z": 3,
	}
	binOperations := map[string]func(int, int, bool){
		"add": alu.Add,
		"mul": alu.Mul,
		"div": alu.Div,
		"mod": alu.Mod,
		"eql": alu.Eql,
	}

	parts := strings.Split(command, " ")
	if parts[0] == "inp" {
		return alu.Inp(reg[parts[1]])
	}

	op := binOperations[parts[0]]
	a := reg[parts[1]]
	b := reg[parts[2]]
	literal := false

	if num, err := strconv.Atoi(parts[2]); err == nil {
		b = num
		literal = true
	}

	op(a, b, literal)
	return nil
}

func (alu *ALU) Inp(a int) error {
	num, err := alu.input.Read()
	if err != nil {
		return err
	}
	alu.Register[0] = num
	return nil
}

func (alu *ALU) Add(a, b int, literal bool) {
	alu.Register[a] += alu.getVar(b, literal)
}

func (alu *ALU) Mul(a, b int, literal bool) {
	alu.Register[a] *= alu.getVar(b, literal)
}

func (alu *ALU) Div(a, b int, literal bool) {
	alu.Register[a] /= alu.getVar(b, literal)
}

func (alu *ALU) Mod(a, b int, literal bool) {
	alu.Register[a] %= alu.getVar(b, literal)
}

func (alu *ALU) Eql(a, b int, literal bool) {
	if alu.Register[a] == alu.getVar(b, literal) {
		alu.Register[a] = 1
	} else {
		alu.Register[a] = 0
	}
}

func (alu *ALU) getVar(a int, literal bool) int {
	if literal {
		return a
	}
	return alu.Register[a]
}

// util functions

func digits(num int) []int {
	digits := []int{}
	for _, ch := range strings.Split(fmt.Sprint(num), "") {
		d, _ := strconv.Atoi(ch)
		digits = append(digits, d)
	}
	return digits
}

func num(digits []int) int {
	n := 0
	for i := range digits {
		n += pow10(i) * digits[len(digits)-1-i]
	}
	return n
}

func pow10(exp int) int {
	pow := 1
	for e := exp; e > 0; e-- {
		pow *= 10
	}
	return pow
}

func contains(nums []int, num int) bool {
	for _, n := range nums {
		if n == num {
			return true
		}
	}
	return false
}
