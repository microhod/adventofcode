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
	puzzle.NewSolution("No Space Left On Device", part1, part2).Run()
}

func part1() error {
	root, err := parse(InputFile)
	if err != nil {
		return err
	}

	var total int
	for _, dir := range root.GetAllDirectories() {
		if size := dir.Size(); size <= 100000 {
			total += size
		}
	}

	fmt.Println(total)
	return nil
}

const (
	DiskSize = 70000000
	RequiredSpace = 30000000
)

func part2() error {
	root, err := parse(InputFile)
	if err != nil {
		return err
	}

	unusedSpace := DiskSize - root.Size()
	minDirToDeleteSize := RequiredSpace - unusedSpace

	var dirToDeleteSize int
	for _, dir := range root.GetAllDirectories() {
		size := dir.Size()
		if size >= minDirToDeleteSize && (size < dirToDeleteSize || dirToDeleteSize == 0) {
			dirToDeleteSize = size
		}
	}

	fmt.Println(dirToDeleteSize)
	return nil
}

func parse(path string) (*Dir, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var cmds []string
	var outputs [][]string

	var cmdOutput []string
	for _, line := range lines {
		if strings.HasPrefix(line, "$") {
			if len(cmds) > 0 {
				outputs = append(outputs, cmdOutput)
				cmdOutput = []string{}
			}

			cmds = append(cmds, line)
			continue
		}

		cmdOutput = append(cmdOutput, line)
	}
	outputs = append(outputs, cmdOutput)

	return runCommands(cmds, outputs)
}

func runCommands(cmds []string, outputs [][]string) (*Dir, error) {
	root := &Dir{
		Name:  "/",
		Dirs:  map[string]*Dir{},
		Files: map[string]*File{},
	}
	// skip the initial "$ cd /" command
	cmds = cmds[1:]
	outputs = outputs[1:]

	dir := root
	for i := range cmds {
		cmd := cmds[i]
		output := outputs[i]

		parts := strings.Split(cmd, " ")
		command := parts[1]

		switch command {
		case "cd":
			arg := parts[2]
			if arg == ".." {
				dir = dir.Parent
				continue
			}
			dir = dir.Dirs[arg]
		case "ls":
			parseLsOutput(output, dir)
		}
	}

	return root, nil
}

func parseLsOutput(output []string, parent *Dir) error {
	for _, line := range output {
		if strings.HasPrefix(line, "dir ") {
			name := strings.Split(line, " ")[1]
			dir := &Dir{
				Name:   name,
				Parent: parent,
				Files:  map[string]*File{},
				Dirs:   map[string]*Dir{},
			}
			parent.Dirs[name] = dir
			continue
		}

		parts := strings.Split(line, " ")
		s, name := parts[0], parts[1]
		size, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		file := &File{
			Name: name,
			Size: size,
		}
		parent.Files[name] = file
	}

	return nil
}

type Dir struct {
	Name   string
	Parent *Dir
	Files  map[string]*File
	Dirs   map[string]*Dir
}

type File struct {
	Name string
	Size int
}

func (d *Dir) GetAllDirectories() []*Dir {
	dirs := []*Dir{d}
	for _, child := range d.Dirs {
		dirs = append(dirs, child.GetAllDirectories()...)
	}
	return dirs
}

func (d *Dir) Size() int {
	var size int
	for _, file := range d.Files {
		size += file.Size
	}

	for _, dir := range d.Dirs {
		size += dir.Size()
	}
	return size
}
