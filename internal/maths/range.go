package maths

type Range struct {
	Left, Right int
}

func (r Range) ContainsValue(n int) bool {
	return n >= r.Left && n <= r.Right
}

func (r Range) Contains(s Range) bool {
	return s.Left >= r.Left && s.Right <= r.Right
}

func (r Range) Intersect(s Range) Range {
	return Range{Max(r.Left, s.Left), Min(r.Right, s.Right)}
}

func (r Range) Intersects(s Range) bool {
	sStartsInsideR := s.Left >= r.Left && s.Left <= r.Right
	rStartsInsideS := r.Left >= s.Left && r.Left <= s.Right

	return sStartsInsideR || rStartsInsideS
}

func (r Range) Valid() bool {
	return r.Left <= r.Right
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
