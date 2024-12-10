package main

import (
	"fmt"
	"slices"
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
	puzzle.NewSolution("Disk Fragmenter", part1, part2).Run()
}

func part1() error {
	disk, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Printf("block-compacted filesystem checksum: %d\n", disk.BlockCompact().Checksum())
	return nil
}

func part2() error {
	disk, err := parse(InputFile)
	if err != nil {
		return err
	}

	fmt.Printf("file-compacted filesystem checksum: %d\n", disk.FileCompact().Checksum())
	return nil
}

func parse(path string) (Disk, error) {
	input, err := file.Read(path)
	if err != nil {
		return nil, err
	}
	input = strings.TrimSpace(input)

	var disk Disk
	for i := 0; i < len(input); i += 2 {
		length, err := strconv.Atoi(string(input[i]))
		if err != nil {
			return nil, err
		}
		var spaceAfter int
		if i+1 < len(input) {
			spaceAfter, err = strconv.Atoi(string(input[i+1]))
			if err != nil {
				return nil, err
			}
		}
		disk = append(disk, File{ID: i / 2, Length: length, SpaceAfter: spaceAfter})
	}
	return disk, nil
}

type Disk []File

func (d Disk) BlockCompact() Disk {
	left := 0
	right := len(d) - 1

	for left < right {
		if d[left].SpaceAfter == 0 {
			left++
			continue
		}
		if d[right].Length == 0 {
			d = d[:right]
			right--
			continue
		}

		move := min(d[left].SpaceAfter, d[right].Length)

		d = slices.Insert(d, left+1, File{
			ID:         d[right].ID,
			Length:     move,
			SpaceAfter: d[left].SpaceAfter - move,
		})
		right += 1
		d[right].Length -= move
		d[left].SpaceAfter = 0
	}
	return d
}

func (d Disk) FileCompact() Disk {
	right := len(d) - 1
	for right >= 0 {
		var success bool
		d, success = d.tryFileCompact(right)
		if !success {
			right--
		}
	}
	return d
}

func (d Disk) tryFileCompact(i int) (Disk, bool) {
	f := d[i]
	for left := 0; left < i; left++ {
		if d[left].SpaceAfter < f.Length {
			continue
		}

		d = append(d[:i], d[i+1:]...)
		d = slices.Insert(d, left+1, File{
			ID:         f.ID,
			Length:     f.Length,
			SpaceAfter: d[left].SpaceAfter - f.Length,
		})
		d[i].SpaceAfter += f.Length + f.SpaceAfter
		d[left].SpaceAfter = 0
		return d, true
	}
	return d, false
}

func (d Disk) Checksum() int {
	var sum int
	var id int
	for _, f := range d {
		for range f.Length {
			sum += id * f.ID
			id++
		}
		id += f.SpaceAfter
	}
	return sum
}

func (d Disk) String() string {
	builder := new(strings.Builder)
	for i := range d {
		builder.WriteString(
			strings.Repeat(fmt.Sprint(d[i].ID), d[i].Length) + strings.Repeat(".", d[i].SpaceAfter),
		)
	}
	return builder.String()
}

type File struct {
	ID         int
	Length     int
	SpaceAfter int
}
