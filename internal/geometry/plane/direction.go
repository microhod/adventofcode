package plane

import (
	"fmt"

	"github.com/microhod/adventofcode/internal/maths"
)

type Direction int

const (
	North     Direction = 0
	NorthEast Direction = 1
	East      Direction = 2
	SouthEast Direction = 3
	South     Direction = 4
	SouthWest Direction = 5
	West      Direction = 6
	NorthWest Direction = 7
)

func (d Direction) Orthogonal() bool {
	return int(d) % 2 == 0
}

func (d Direction) Opposite() Direction {
	return Direction(maths.Mod(int(d)+4, 8))
}

// Turn expects degrees to be multiple of 45. +ve means clockwise, -ve anti-clockwise
func(d Direction) Turn(degrees int) Direction {
	if degrees % 45 != 0 {
		panic(fmt.Sprintf("turn: degress must be multiple of 45, got: %d", degrees))
	}

	return Direction(maths.Mod(int(d) + (degrees / 45), 8))
}

var directionNames = []string{"North", "NorthEast", "East", "SouthEast", "South", "SouthWest", "West", "NorthWest"}

func (d Direction) String() string {
	return directionNames[d]
}

// 7 0 1
// 6   2
// 5 4 3
var DirectionToVector = map[Direction]Vector{
	North:     {0, -1},
	NorthEast: {1, -1},
	East:      {1, 0},
	SouthEast: {1, 1},
	South:     {0, 1},
	SouthWest: {-1, 1},
	West:      {-1, 0},
	NorthWest: {-1, -1},
}
