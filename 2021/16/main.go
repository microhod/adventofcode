package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"

	PacketTypeLiteral     = 4
	PacketTypeSum         = 0
	PacketTypeProduct     = 1
	PacketTypeMin         = 2
	PacketTypeMax         = 3
	PacketTypeGreaterThan = 5
	PacketTypeLessThan    = 6
	PacketTypeEqual       = 7

	LengthTypeLength = 0
	LengthTypeCount  = 1
)

func main() {
	puzzle.NewSolution("Packet Decoder", part1, part2).Run()
}

func part1() error {
	binary, err := parseBinary(InputFile)
	if err != nil {
		return err
	}

	packet, _, err := parsePacket(binary)
	if err != nil {
		return err
	}

	fmt.Printf("sum of all packet versions: %d\n", packet.SumVersions())

	return nil
}

func part2() error {
	binary, err := parseBinary(InputFile)
	if err != nil {
		return err
	}

	packet, _, err := parsePacket(binary)
	if err != nil {
		return err
	}

	fmt.Printf("value of outer packet: %d\n", packet.GetValue())

	return nil
}

func parseBinary(path string) (*Binary, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, err
	}

	binary := ""
	for _, char := range lines[0] {
		switch char {
		case rune('0'):
			binary += "0000"
		case rune('1'):
			binary += "0001"
		case rune('2'):
			binary += "0010"
		case rune('3'):
			binary += "0011"
		case rune('4'):
			binary += "0100"
		case rune('5'):
			binary += "0101"
		case rune('6'):
			binary += "0110"
		case rune('7'):
			binary += "0111"
		case rune('8'):
			binary += "1000"
		case rune('9'):
			binary += "1001"
		case rune('A'):
			binary += "1010"
		case rune('B'):
			binary += "1011"
		case rune('C'):
			binary += "1100"
		case rune('D'):
			binary += "1101"
		case rune('E'):
			binary += "1110"
		case rune('F'):
			binary += "1111"
		}
	}

	return &Binary{binary}, nil
}

type Binary struct {
	str string
}

func (b *Binary) Pop(num int) string {
	tmp := b.str[:num]
	b.str = b.str[num:]

	return tmp
}

func (b *Binary) Length() int {
	return len(b.str)
}

func parsePacket(binary *Binary) (*Packet, int, error) {
	var err error
	packet := &Packet{}
	length := 0

	packet.Version, err = strconv.ParseInt(binary.Pop(3), 2, 64)
	if err != nil {
		return nil, 0, err
	}
	packet.TypeID, err = strconv.ParseInt(binary.Pop(3), 2, 64)
	if err != nil {
		return nil, 0, err
	}

	switch packet.TypeID {
	case PacketTypeLiteral:
		packet.Value, length, err = parseLiteralValue(binary)
		if err != nil {
			return nil, 0, err
		}
	default:
		packet.SubPackets, length, err = parseSubPackets(binary)
		if err != nil {
			return nil, 0, err
		}
	}

	return packet, length + 6, nil
}

func parseLiteralValue(binary *Binary) (int, int, error) {
	var err error

	valueBin := ""
	prefix := "1"
	length := 0
	for prefix != "0" {
		prefix = binary.Pop(1)
		valueBin += binary.Pop(4)
		length += 5
	}

	value, err := strconv.ParseInt(valueBin, 2, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(value), length, nil
}

func parseSubPackets(binary *Binary) ([]*Packet, int, error) {
	lengthTypeID, err := strconv.Atoi(binary.Pop(1))
	if err != nil {
		return nil, 0, err
	}

	switch lengthTypeID {
	case LengthTypeLength:
		packets, length, err := parseSubPacketsByLength(binary)
		if err != nil {
			return nil, 0, err
		}
		return packets, length + 1, nil
	case LengthTypeCount:
		packets, length, err := parseSubPacketsByCount(binary)
		if err != nil {
			return nil, 0, err
		}
		return packets, length + 1, nil
	default:
		return nil, 0, fmt.Errorf("invalid length type: %d", lengthTypeID)
	}
}

func parseSubPacketsByLength(binary *Binary) ([]*Packet, int, error) {
	expectedLength, err := strconv.ParseInt(binary.Pop(15), 2, 64)
	if err != nil {
		return nil, 0, err
	}

	packets := []*Packet{}
	length := 0
	for length < int(expectedLength) {
		p, l, err := parsePacket(binary)
		if err != nil {
			return nil, 0, err
		}

		length += l
		packets = append(packets, p)
	}

	return packets, length + 15, nil
}

func parseSubPacketsByCount(binary *Binary) ([]*Packet, int, error) {
	expectedCount, err := strconv.ParseInt(binary.Pop(11), 2, 64)
	if err != nil {
		return nil, 0, err
	}

	packets := []*Packet{}
	length := 0
	count := 0
	for count < int(expectedCount) {
		p, l, err := parsePacket(binary)
		if err != nil {
			return nil, 0, err
		}

		length += l
		packets = append(packets, p)
		count += 1
	}

	return packets, length + 11, nil
}

type Packet struct {
	Version    int64
	TypeID     int64
	Value      int
	SubPackets []*Packet
}

func (packet *Packet) SumVersions() int {
	sum := packet.Version
	toSum := []*Packet{}
	toSum = append(toSum, packet.SubPackets...)

	for len(toSum) > 0 {
		p := toSum[0]
		toSum = toSum[1:]

		sum += p.Version
		toSum = append(toSum, p.SubPackets...)
	}

	return int(sum)
}

func (packet *Packet) GetValue() int {
	switch packet.TypeID {
	case PacketTypeLiteral:
		return packet.Value
	case PacketTypeSum:
		return sum(packet.SubPackets)
	case PacketTypeProduct:
		return product(packet.SubPackets)
	case PacketTypeMin:
		return min(packet.SubPackets)
	case PacketTypeMax:
		return max(packet.SubPackets)
	case PacketTypeGreaterThan:
		return greaterThan(packet.SubPackets[0], packet.SubPackets[1])
	case PacketTypeLessThan:
		return lessThan(packet.SubPackets[0], packet.SubPackets[1])
	case PacketTypeEqual:
		return equal(packet.SubPackets[0], packet.SubPackets[1])
	}

	return 0
}

func sum(packets []*Packet) int {
	sum := 0
	for _, p := range packets {
		sum += p.GetValue()
	}
	return sum
}

func product(packets []*Packet) int {
	product := 1
	for _, p := range packets {
		product *= p.GetValue()
	}
	return product
}

func min(packets []*Packet) int {
	min := math.MaxInt
	for _, p := range packets {
		v := p.GetValue()
		if v < min {
			min = v
		}
	}
	return min
}

func max(packets []*Packet) int {
	max := 0
	for _, p := range packets {
		v := p.GetValue()
		if v > max {
			max = v
		}
	}
	return max
}

func greaterThan(p *Packet, q *Packet) int {
	if p.GetValue() > q.GetValue() {
		return 1
	}
	return 0
}

func lessThan(p *Packet, q *Packet) int {
	if p.GetValue() < q.GetValue() {
		return 1
	}
	return 0
}

func equal(p *Packet, q *Packet) int {
	if p.GetValue() == q.GetValue() {
		return 1
	}
	return 0
}
