package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Wait For It", part1, part2).Run()
}

func part1() error {
	races, _, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	product := 1
	for _, race := range races {
		recordBreaks := 0
		for chargeTime := 0; chargeTime < race.Time; chargeTime++ {
			if race.BrokeRecord(chargeTime) {
				recordBreaks += 1
			}
		}
		product *= recordBreaks
	}

	fmt.Printf("the number of ways to break the record in each race: %d\n", product)
	return nil
}

func part2() error {
	_, race, err := parseInput(InputFile)
	if err != nil {
		return err
	}

	recordBreaks := 0
	for chargeTime := 0; chargeTime < race.Time; chargeTime++ {
		if race.BrokeRecord(chargeTime) {
			recordBreaks += 1
		}
	}
	fmt.Printf("the number of ways to break the race record: %d\n", recordBreaks)

	return nil
}

func parseInput(path string) ([]Race, Race, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, Race{}, err
	}

	races, err := parseMultipleRaces(lines)
	if err != nil {
		return nil, Race{}, err
	}
	race, err := parseSingleRace(lines)
	if err != nil {
		return nil, Race{}, err
	}

	return races, race, nil
}

func parseMultipleRaces(lines []string) ([]Race, error) {
	times, err := csv.ParseInts(lines[0][9:], " ")
	if err != nil {
		return nil, err
	}
	distances, err := csv.ParseInts(lines[1][9:], " ")
	if err != nil {
		return nil, err
	}

	var races []Race
	for i := range times {
		races = append(races, Race{
			Time:           times[i],
			DistanceRecord: distances[i],
		})
	}
	return races, nil
}

func parseSingleRace(lines []string) (Race, error) {
	time, err := strconv.Atoi(strings.ReplaceAll(lines[0][9:], " ", ""))
	if err != nil {
		return Race{}, err
	}
	distance, err := strconv.Atoi(strings.ReplaceAll(lines[1][9:], " ", ""))
	if err != nil {
		return Race{}, err
	}
	return Race{Time: time, DistanceRecord: distance}, nil
}

type Race struct {
	Time           int
	DistanceRecord int
}

func (r Race) BrokeRecord(chargeTime int) bool {
	return (r.Time-chargeTime)*chargeTime > r.DistanceRecord
}
