package main

import (
	"encoding/json"
	"fmt"
	"sort"
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
	puzzle.NewSolution("Distress Signal", part1, part2).Run()
}

func part1() error {
	pairs, _, err := parse(InputFile)
	if err != nil {
		return err
	}

	var rightOrderIndexSum int
	for index, pair := range pairs {
		if InOrder(pair[0], pair[1]) == RightOrder {
			rightOrderIndexSum += index + 1
		}
	}

	fmt.Println(rightOrderIndexSum)
	return nil
}

func part2() error {
	_, packets, err := parse(InputFile)
	if err != nil {
		return err
	}

	dividerPackets := []string{"[[2]]","[[6]]"}
	packets = append(packets, dividerPackets...)

	sort.Slice(packets, func(i, j int) bool {
		return LessThan(packets[i], packets[j])
	})

	decoderKey := 1
	for index, packet := range packets {
		if packet == dividerPackets[0] || packet == dividerPackets[1] {
			decoderKey *= index+1
		}
	}

	fmt.Println(decoderKey)
	return nil
}

func parse(path string) ([][2]string, []string, error) {
	data, err := file.ReadBytes(path)
	if err != nil {
		return nil, nil, err
	}

	var pairs [][2]string
	var packets []string

	for _, section := range strings.Split(string(data), "\n\n") {
		section := strings.TrimSpace(section)
		lines := strings.Split(section, "\n")

		pairs = append(pairs, [2]string{lines[0], lines[1]})
		packets = append(packets, lines[0], lines[1])
	}

	return pairs, packets, nil
}

type OrderState string

const (
	UndecidedOrder OrderState = "undecided"
	RightOrder     OrderState = "right"
	WrongOrder     OrderState = "wrong"
)

func LessThan(packet1, packet2 string) bool {
	order := InOrder(packet1, packet2)

	// the order should never be undecided
	if order == UndecidedOrder {
		panic("undecided order")
	}
	return order == RightOrder
}

func InOrder(left, right string) OrderState {
	leftNum, err := strconv.Atoi(left)
	isLeftNum := err == nil
	rightNum, err := strconv.Atoi(right)
	isRightNum := err == nil

	if isLeftNum && isRightNum {
		if leftNum < rightNum {
			return RightOrder
		}
		if leftNum > rightNum {
			return WrongOrder
		}
		return UndecidedOrder
	}
	if isLeftNum && !isRightNum {
		left = fmt.Sprintf("[%s]", left)
	}
	if !isLeftNum && isRightNum {
		right = fmt.Sprintf("[%s]", right)
	}

	leftList := SplitList(left)
	rightList := SplitList(right)

	for i := 0; i < max(len(leftList), len(rightList)); i++ {
		if i > len(leftList)-1 {
			return RightOrder
		}
		if i > len(rightList)-1 {
			return WrongOrder
		}

		order := InOrder(leftList[i], rightList[i])
		if order != UndecidedOrder {
			return order
		}
	}
	return UndecidedOrder
}

func SplitList(str string) []string {
	var parts []interface{}
	err := json.Unmarshal([]byte(str), &parts)
	if err != nil {
		panic(err)
	}

	var list []string
	for _, element := range parts {
		marshalled, err := json.Marshal(element)
		if err != nil {
			panic(err)
		}

		list = append(list, string(marshalled))
	}

	return list
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
