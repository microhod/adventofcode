// Package csv is a utility package for parsing csv encodings
package csv

import (
	"strconv"
	"strings"
)

func ParseInts(str string) ([]int, error) {
	nums := []int{}
	for _, s := range strings.Split(str, ",") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}

		nums = append(nums, n)
	}

	return nums, nil
}
