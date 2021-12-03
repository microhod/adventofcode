package main

import (
	"fmt"
	"strconv"

	"github.com/microhod/adventofcode/internal/file"
	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	DiagnosticFile = "diagnostic.txt"
	TestFile       = "test.txt"
)

type DiagnosticReport struct {
	BitLength    int
	Measurements []int
}

func main() {
	puzzle.NewSolution("Binary Diagnostic", part1, part2).Run()
}

func part1() error {
	report, err := readDiagnosticReport(DiagnosticFile)
	if err != nil {
		return err
	}

	g := gamma(report)
	e := epsilon(report)

	fmt.Printf("gamma: %d\n", g)
	fmt.Printf("epsilon: %d\n", e)
	fmt.Printf("gamma * epsilon: %d\n", g*e)

	return nil
}

func part2() error {
	report, err := readDiagnosticReport(DiagnosticFile)
	if err != nil {
		return err
	}

	oxy := oxygenGeneratorRating(report)
	co2 := co2ScrubberRating(report)

	fmt.Printf("oxygen generator rating: %d\n", oxy)
	fmt.Printf("CO2 scrubber rating: %d\n", co2)
	fmt.Printf("oxy * co2: %d\n", oxy*co2)

	return nil
}

func oxygenGeneratorRating(report *DiagnosticReport) int {
	nums := report.Measurements
	for i := report.BitLength - 1; i >= 0; i-- {
		if len(nums) < 2 {
			break
		}

		on := countBitsSwitched(nums, uint(i))
		
		mostCommon := 0
		if float64(on) >= float64(len(nums)) / 2 {
			mostCommon = 1
		}

		nums = filterByBitValue(nums, uint(i), mostCommon == 1)
	}

	if len(nums) < 1 {
		return -1
	}
	return nums[0]
}

func co2ScrubberRating(report *DiagnosticReport) int {
	nums := report.Measurements
	for i := report.BitLength - 1; i >= 0; i-- {
		if len(nums) < 2 {
			break
		}

		on := countBitsSwitched(nums, uint(i))
		
		leastCommon := 0
		if float64(on) < float64(len(nums)) / 2 {
			leastCommon = 1
		}

		nums = filterByBitValue(nums, uint(i), leastCommon == 1)
	}

	if len(nums) < 1 {
		return -1
	}
	return nums[0]
}

func gamma(report *DiagnosticReport) int {
	gamma := 0
	for i := 0; i < report.BitLength; i++ {
		on := countBitsSwitched(report.Measurements, uint(i))

		if on > len(report.Measurements)/2 {
			// if most bits are on switch the bit on in gamma
			gamma |= (1 << i)
		}
	}

	return gamma
}

func epsilon(report *DiagnosticReport) int {
	epsilon := 0
	// work backwards
	for i := 0; i < report.BitLength; i++ {
		on := countBitsSwitched(report.Measurements, uint(i))

		if on <= len(report.Measurements)/2 {
			// if most bits are off switch the bit on in gamma
			epsilon |= (1 << i)
		}
	}

	return epsilon
}

func filterByBitValue(nums []int, pos uint, switched bool) []int {
	filtered := []int{}
	for _, n := range(nums) {
		if hasBitSwitched(n, pos) == switched {
			filtered = append(filtered, n)
		}
	}
	return filtered
}

func countBitsSwitched(nums []int, pos uint) int {
	on := 0
	for _, n := range nums {
		if hasBitSwitched(n, pos) {
			on += 1
		}
	}

	return on
}

func hasBitSwitched(n int, pos uint) bool {
	return n&(1<<pos) > 0
}

func readDiagnosticReport(path string) (*DiagnosticReport, error) {
	lines, err := file.ReadLines(path)
	if err != nil {
		return nil, err
	}

	report := &DiagnosticReport{}
	for i, line := range lines {
		if report.BitLength == 0 {
			report.BitLength = len(line)
		}
		if report.BitLength != len(line) {
			return nil, fmt.Errorf("line %d: bit length %d did not match prevous length %d", i+1, len(line), report.BitLength)
		}

		measurement, err := strconv.ParseInt(line, 2, 32)
		if err != nil {
			return nil, err
		}
		report.Measurements = append(report.Measurements, int(measurement))
	}

	return report, nil
}
