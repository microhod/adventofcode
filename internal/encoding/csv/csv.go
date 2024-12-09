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
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}

		nums = append(nums, n)
	}

	return nums, nil
}

func ParseInt64s(str string, separator ...string) ([]int64, error) {
	if len(separator) < 1 {
		separator = []string{","}
	}

	var nums []int64
	for _, s := range strings.Split(str, separator[0]) {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}

		nums = append(nums, n)
	}

	return nums, nil
}
