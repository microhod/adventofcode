package plane

import "github.com/microhod/adventofcode/internal/maths"

type Vector struct {
	X, Y int
}

func (v Vector) Add(u Vector) Vector {
	return Vector{X: v.X + u.X, Y: v.Y + u.Y}
}

func (v Vector) Minus(u Vector) Vector {
	return Vector{X: v.X - u.X, Y: v.Y - u.Y}
}

type Metric func(Vector, Vector) int

func ManhattanMetric(v, u Vector) int {
	diff := v.Minus(u)
	return maths.Abs(diff.X) + maths.Abs(diff.Y)
}

func (v Vector) Neighbours() map[Direction]Vector {
	neighbours := map[Direction]Vector{}
	for direction, diff := range DirectionToVector {
		neighbours[direction] = v.Add(diff)
	}
	return neighbours
}

func (v Vector) OrthogonalNeighbours() map[Direction]Vector {
	neighbours := map[Direction]Vector{}
	for direction, diff := range DirectionToVector {
		if !direction.Orthogonal() {
			continue
		}
		neighbours[direction] = v.Add(diff)
	}
	return neighbours
}
