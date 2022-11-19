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
	puzzle.NewSolution("Reindeer Olympics", part1, part2).Run()
}

func part1() error {
	deer, err := parse(InputFile)
	if err != nil {
		return err
	}

	winner, distance := DistanceRace(deer, 2503)

	fmt.Printf("And the winner is....%s with an incredible distance of %dkm!\n", winner.Name, distance)
	return nil
}

func part2() error {
	deer, err := parse(InputFile)
	if err != nil {
		return err
	}

	winner, distance := PointsRace(deer, 2503)

	fmt.Printf("And the winner is....%s with %d points!\n", winner, distance)
	return nil
}

func parse(path string) ([]Reindeer, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var deer []Reindeer
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Vixen can fly 19 km/s for 7 seconds, but then must rest for 124 seconds.
		fields := strings.Fields(line)
		name, s, flySecs, restSecs := fields[0], fields[3], fields[6], fields[13]

		speed, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		flySeconds, err := strconv.Atoi(flySecs)
		if err != nil {
			return nil, err
		}
		restSeconds, err := strconv.Atoi(restSecs)
		if err != nil {
			return nil, err
		}

		deer = append(deer, Reindeer{
			Name:        name,
			Speed:       speed,
			FlySeconds:  flySeconds,
			RestSeconds: restSeconds,
		})
	}

	return deer, nil
}

type Reindeer struct {
	Name        string
	Speed       int
	FlySeconds  int
	RestSeconds int
}

func (r Reindeer) NewFlight() *Flight {
	return &Flight{
		Deer:     r,
		Distance: 0,
		Flying:   true,
		FlyLeft:  r.FlySeconds,
		RestLeft: r.RestSeconds,
	}
}

type Flight struct {
	Deer     Reindeer
	Distance int
	Flying   bool
	FlyLeft  int
	RestLeft int
}

func (f *Flight) Tick() {
	if f.Flying {
		f.FlyLeft--
		f.Distance += f.Deer.Speed
		if f.FlyLeft == 0 {
			f.FlyLeft = f.Deer.FlySeconds
			f.Flying = false
		}
	} else {
		f.RestLeft--
		if f.RestLeft == 0 {
			f.RestLeft = f.Deer.RestSeconds
			f.Flying = true
		}
	}
}

func DistanceRace(deer []Reindeer, seconds int) (Reindeer, int) {
	var flights []*Flight
	for _, d := range deer {
		flights = append(flights, d.NewFlight())
	}

	for i := 0; i < seconds; i++ {
		for _, f := range flights {
			f.Tick()
		}
	}

	var winner Reindeer
	var distance int

	for _, f := range flights {
		if f.Distance > distance {
			winner = f.Deer
			distance = f.Distance
		}
	}

	return winner, distance
}

func PointsRace(deer []Reindeer, seconds int) (string, int) {
	var flights []*Flight
	for _, d := range deer {
		flights = append(flights, d.NewFlight())
	}

	points := map[string]int{}
	for i := 0; i < seconds; i++ {
		for _, f := range flights {
			f.Tick()
		}

		var winning []string
		var distance int
		for _, f := range flights {
			if f.Distance == distance {
				winning = append(winning, f.Deer.Name)
			}
			if f.Distance > distance {
				winning = []string{f.Deer.Name}
				distance = f.Distance
			}
		}

		for _, deer := range winning {
			points[deer]++
		}
	}

	var winner string
	var maxPoints int

	for deer, points := range points {
		if points > maxPoints {
			winner = deer
			maxPoints = points
		}
	}

	return winner, maxPoints
}
