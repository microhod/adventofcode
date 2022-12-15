package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/microhod/adventofcode/internal/encoding/csv"
	"github.com/microhod/adventofcode/internal/file"
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
		distress := Vector{x, y}

		for _, sensor := range sensors {
			if sensor.CanDetect(distress) && !distress.Equals(sensor.NearestBeacon) {
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
		possibleDistressBeacons := []Range{{Left: 0, Right: limit}}

		for _, sensor := range sensors {
			sensorDetection := sensor.DetectionAreaCrossSection(y)
			if sensorDetection == nil {
				continue
			}

			// remove any ranges the sensors would have already detected
			var outOfDetection []Range
			for _, r := range possibleDistressBeacons {
				outOfDetection = append(outOfDetection, r.Diff(*sensorDetection)...)
			}
			possibleDistressBeacons = outOfDetection
		}

		// if we have any ranges after minusing all the beacon detection ranges,
		// we have found the distress beacon!
		if len(possibleDistressBeacons) > 0 {
			distressSignal := Vector{
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

		sensorOrigin := Vector{X: nums[0], Y: nums[1]}
		nearestBeacon := Vector{X: nums[2], Y: nums[3]}

		sensors = append(sensors, Sensor{
			Origin:        sensorOrigin,
			Radius:        ManhattanDistance(sensorOrigin, nearestBeacon),
			NearestBeacon: nearestBeacon,
		})
	}

	return sensors, nil
}

type Sensor struct {
	Origin        Vector
	Radius        int
	NearestBeacon Vector
}

func (s Sensor) CanDetect(v Vector) bool {
	return ManhattanDistance(s.Origin, v) <= s.Radius
}

func (s Sensor) DetectionAreaCrossSection(y int) *Range {
	if y > s.Origin.Y+s.Radius || y < s.Origin.Y-s.Radius {
		// outside this sensor's detection area
		return nil
	}

	offset := maths.Abs(s.Origin.Y - y)

	return &Range{
		Left:  s.Origin.X - (s.Radius - offset),
		Right: s.Origin.X + (s.Radius - offset),
	}
}

type Vector struct {
	X, Y int
}

func (v Vector) Equals(u Vector) bool {
	return v.X == u.X && v.Y == u.Y
}

func ManhattanDistance(u Vector, v Vector) int {
	return maths.Abs(u.X-v.X) + maths.Abs(u.Y-v.Y)
}

type Range struct {
	Left  int
	Right int
}

func (r Range) Diff(remove Range) []Range {
	// distinct
	// [ r ]
	//        [ remove ]
	if remove.Right < r.Left || remove.Left > r.Right {
		return []Range{r}
	}
	// remove contains r
	//   [-r-]
	// [ remove ]
	if remove.Left <= r.Left && remove.Right >= r.Right {
		return nil
	}
	// r contains remove
	// [  |-----r----|  ]
	//     [ remove ]
	if r.Left <= remove.Left && r.Right >= remove.Right {
		// split into two ranges
		var split []Range

		if r.Left < remove.Left {
			split = append(split, Range{
				Left:  r.Left,
				Right: remove.Left - 1,
			})
		}
		if r.Right > remove.Right {
			split = append(split, Range{
				Left:  remove.Right + 1,
				Right: r.Right,
			})
		}
		return split
	}
	// trim right
	// [  |---r---]
	//     [ remove ]
	if r.Left < remove.Left {
		return []Range{{
			Left:  r.Left,
			Right: remove.Left - 1,
		}}
	}
	// trim left
	//      [----| r   ]
	// [ remove ]
	return []Range{{
		Left:  remove.Right + 1,
		Right: r.Right,
	}}
}
