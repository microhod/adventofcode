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
		less, err := pair[0].LessThan(pair[1])
		if err != nil {
			return err
		}

		if less {
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

	dividerPackets := []Packet{"[[2]]", "[[6]]"}
	packets = append(packets, dividerPackets...)

	sort.Slice(packets, func(i, j int) bool {
		less, err := packets[i].LessThan(packets[j])
		if err != nil {
			// we kind of have to panic here
			panic(err)
		}
		return less
	})

	decoderKey := 1
	for index, packet := range packets {
		if packet == dividerPackets[0] || packet == dividerPackets[1] {
			decoderKey *= index + 1
		}
	}

	fmt.Println(decoderKey)
	return nil
}

func parse(path string) ([][2]Packet, []Packet, error) {
	data, err := file.ReadBytes(path)
	if err != nil {
		return nil, nil, err
	}

	var pairs [][2]Packet
	var packets []Packet

	for _, section := range strings.Split(string(data), "\n\n") {
		section := strings.TrimSpace(section)
		lines := strings.Split(section, "\n")

		pairs = append(pairs, [2]Packet{Packet(lines[0]), Packet(lines[1])})
		packets = append(packets, Packet(lines[0]), Packet(lines[1]))
	}

	return pairs, packets, nil
}

type Packet string

func (p Packet) LessThan(q Packet) (bool, error) {
	order, err := packetsInOrder(p, q)
	if err != nil {
		return false, fmt.Errorf("failed to check order: %w", err)
	}

	// the order should never be undecided
	if order == undecidedPacketOrder {
		panic("undecided order")
	}
	return order == rightPacketOrder, nil
}

type packetOrder string

const (
	undecidedPacketOrder packetOrder = "undecided"
	rightPacketOrder     packetOrder = "right"
	wrongPacketOrder     packetOrder = "wrong"
)

func packetsInOrder(left, right Packet) (packetOrder, error) {
	leftNum, err := strconv.Atoi(string(left))
	isLeftNum := err == nil
	rightNum, err := strconv.Atoi(string(right))
	isRightNum := err == nil

	if isLeftNum && isRightNum {
		if leftNum < rightNum {
			return rightPacketOrder, nil
		}
		if leftNum > rightNum {
			return wrongPacketOrder, nil
		}
		return undecidedPacketOrder, nil
	}
	if isLeftNum && !isRightNum {
		left = Packet(fmt.Sprintf("[%s]", left))
	}
	if !isLeftNum && isRightNum {
		right = Packet(fmt.Sprintf("[%s]", right))
	}

	leftPackets, err := left.split()
	if err != nil {
		return undecidedPacketOrder, err
	}
	rightPackets, err := right.split()
	if err != nil {
		return undecidedPacketOrder, err
	}

	for i := 0; i < max(len(leftPackets), len(rightPackets)); i++ {
		if i > len(leftPackets)-1 {
			return rightPacketOrder, nil
		}
		if i > len(rightPackets)-1 {
			return wrongPacketOrder, nil
		}

		order, err := packetsInOrder(leftPackets[i], rightPackets[i])
		if err != nil {
			return undecidedPacketOrder, err
		}
		if order != undecidedPacketOrder {
			return order, nil
		}
	}
	return undecidedPacketOrder, nil
}

func (p Packet) split() ([]Packet, error) {
	var parts []interface{}
	err := json.Unmarshal([]byte(p), &parts)
	if err != nil {
		return nil, fmt.Errorf("malformed list package: %w", err)
	}

	var list []Packet
	for _, element := range parts {
		marshalled, err := json.Marshal(element)
		if err != nil {
			return nil, fmt.Errorf("malformed list package element: %w", err)
		}

		list = append(list, Packet(marshalled))
	}

	return list, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
