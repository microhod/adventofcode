package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/geometry/plane"
	"github.com/microhod/adventofcode/internal/maths"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Beacon Exclusion Zone", part1, part2).Run()
}

func part1() error {
	sensors, err := parse(InputFile)
	if err != nil {
		return err
	}

	y := 2000000

	minX := math.MaxInt
	maxX := math.MinInt
	for _, sensor := range sensors {
		minX = maths.Min(minX, sensor.Origin.X-sensor.Radius)
		maxX = maths.Max(maxX, sensor.Origin.X+sensor.Radius)
	}

	var distressNotHere int
	for x := minX; x <= maxX; x++ {
		distress := plane.Vector{X: x, Y: y}

		for _, sensor := range sensors {
			if sensor.CanDetect(distress) && distress != sensor.NearestBeacon {
				distressNotHere += 1
				break
			}
		}
	}

	fmt.Println(distressNotHere)
	return nil
}

func part2() error {
	sensors, err := parse(InputFile)
	if err != nil {
		return err
	}

	limit := 4000000

	for y := 0; y <= limit; y++ {
		possibleDistressBeacons := []maths.Range{{Left: 0, Right: limit}}

		for _, sensor := range sensors {
			sensorDetection := sensor.DetectionAreaCrossSection(y)
			if sensorDetection == nil {
				continue
			}

			// remove any ranges the sensors would have already detected
			var outOfDetection []maths.Range
			for _, r := range possibleDistressBeacons {
				outOfDetection = append(outOfDetection, r.Diff(*sensorDetection)...)
			}
			possibleDistressBeacons = outOfDetection
		}

		// if we have any ranges after minusing all the beacon detection ranges,
		// we have found the distress beacon!
		if len(possibleDistressBeacons) > 0 {
			distressSignal := plane.Vector{
				X: possibleDistressBeacons[0].Left,
				Y: y,
			}
			fmt.Println(distressSignal.X*4000000 + distressSignal.Y)
			return nil
		}
	}

	return nil
}

func parse(path string) ([]Sensor, error) {
	b, err := file.ReadBytes(path)
	if err != nil {
		return nil, err
	}
	input := string(b)

	input = strings.ReplaceAll(input, "Sensor at ", "")
	input = strings.ReplaceAll(input, ": closest beacon is at ", ", ")
	input = strings.ReplaceAll(input, "x=", "")
	input = strings.ReplaceAll(input, "y=", "")

	var sensors []Sensor
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		nums, err := csv.ParseInts(line, ", ")
		if err != nil {
			return nil, err
		}

		sensorOrigin := plane.Vector{X: nums[0], Y: nums[1]}
		nearestBeacon := plane.Vector{X: nums[2], Y: nums[3]}

		sensors = append(sensors, Sensor{
			Origin:        sensorOrigin,
			Radius:        plane.ManhattanMetric(sensorOrigin, nearestBeacon),
			NearestBeacon: nearestBeacon,
		})
	}

	return sensors, nil
}

type Sensor struct {
	Origin        plane.Vector
	Radius        int
	NearestBeacon plane.Vector
}

func (s Sensor) CanDetect(v plane.Vector) bool {
	return plane.ManhattanMetric(s.Origin, v) <= s.Radius
}

func (s Sensor) DetectionAreaCrossSection(y int) *maths.Range {
	if y > s.Origin.Y+s.Radius || y < s.Origin.Y-s.Radius {
		// outside this sensor's detection area
		return nil
	}

	offset := maths.Abs(s.Origin.Y - y)

	return &maths.Range{
		Left:  s.Origin.X - (s.Radius - offset),
		Right: s.Origin.X + (s.Radius - offset),
	}
}
