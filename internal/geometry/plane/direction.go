package plane

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

