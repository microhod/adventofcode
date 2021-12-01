package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	REPORT_FILE = "sonar-sweep-report.txt"
	TEST_FILE   = "test-report.txt"
)

func main() {
	report, err := readReport(REPORT_FILE)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	fmt.Println("--- Part 1 ---")
	fmt.Printf("# depth increases: %d\n", numIncreases(report))

	fmt.Println()

	fmt.Println("--- Part 2 ---")
	fmt.Printf("# depth increases (three-measurement sliding window): %d\n", numIncreases(slidingWindow(report, 3)))

	fmt.Println()
}

func numIncreases(report []int) int {
	if len(report) < 2 {
		return len(report)
	}

	increases := 0
	prev := report[0]
	for _, depth := range report[1:] {
		if depth > prev {
			increases += 1
		}
		prev = depth
	}

	return increases
}

func slidingWindow(report []int, size int) []int {
	if len(report) < size {
		return []int{}
	}

	slidingReport := []int{}
	for i := range report[:len(report)-(size-1)] {
		sum := 0
		for _, depth := range report[i:i+size] {
			sum += depth
		}
		slidingReport = append(slidingReport, sum)
	}

	return slidingReport
}

func readReport(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	report := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		depth, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}

		report = append(report, depth)
	}

	return report, nil
}
