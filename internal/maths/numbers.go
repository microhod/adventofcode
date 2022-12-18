package maths

import "golang.org/x/exp/constraints"

type RealNumber interface {
	constraints.Integer | constraints.Float
}

func Min[T RealNumber](nums ...T) T {
	var min T
	if len(nums) < 1 {
		return min
	}

	min = nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
	}
	return min
}

func Max[T RealNumber](nums ...T) T {
	var max T
	if len(nums) < 1 {
		return max
	}

	max = nums[0]
	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}
	return max
}

func Abs[T RealNumber](num T) T {
	if num < 0 {
		return -num
	}
	return num
}

func Mod(d, m int) int {
	var res int = d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

func Sum[T RealNumber](nums ...T) T {
	var sum T
	for _, n := range nums {
		sum += n
	}
	return sum
}
