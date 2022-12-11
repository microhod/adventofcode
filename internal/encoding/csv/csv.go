// Package csv is a utility package for parsing csv encodings
package csv

import (
	"strconv"
	"strings"
)

func ParseInts(str string, separator ...string) ([]int, error) {
	if len(separator) < 1 {
		separator = []string{","}
	}

	nums := []int{}
	for _, s := range strings.Split(str, separator[0]) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}

		nums = append(nums, n)
	}

	return nums, nil
}
