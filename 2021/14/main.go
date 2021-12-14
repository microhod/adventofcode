package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/microhod/adventofcode/internal/puzzle"
)

const (
	InputFile = "input.txt"
	TestFile  = "test.txt"
)

func main() {
	puzzle.NewSolution("Extended Polymerization", part1, part2).Run()
}

func part1() error {
	polymer, rules, err := readInstructions(InputFile)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		polymer.Execute(rules)
	}

	max := polymer.MaxChar()
	min := polymer.MinChar()
	fmt.Printf("max character - min character = %d - %d = %d\n", max, min, max-min)

	return nil
}

func part2() error {
	polymer, rules, err := readInstructions(InputFile)
	if err != nil {
		return err
	}

	for i := 0; i < 40; i++ {
		polymer.Execute(rules)
	}

	max := polymer.MaxChar()
	min := polymer.MinChar()
	fmt.Printf("max character - min character = %d - %d = %d\n", max, min, max-min)

	return nil
}

func readInstructions(path string) (*Polymer, []Rule, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	parts := strings.Split(string(bytes), "\n\n")
	if len(parts) < 2 {
		return nil, nil, fmt.Errorf("expected 2 sections separated by a blank line only got %d", len(parts))
	}

	polymer := readPolymer(parts[0])
	rules := readRules(strings.Split(parts[1], "\n"))

	return polymer, rules, nil
}

func readPolymer(str string) *Polymer {
	str = strings.TrimSpace(str)

	polymer := &Polymer{
		Characters: map[rune]int{},
		Pairs:      map[string]int{},
	}
	for _, c := range str {
		polymer.Characters[c] += 1
	}

	for i := range str[1:] {
		pair := string(str[i]) + string(str[i+1])
		polymer.Pairs[pair] += 1
	}

	return polymer
}

func readRules(lines []string) []Rule {
	rules := []Rule{}
	for _, line := range lines {
		parts := strings.Split(line, " -> ")
		if len(parts) < 2 {
			continue
		}
		if len(parts[0]) != 2 {
			continue
		}
		if len(parts[1]) != 1 {
			continue
		}

		rules = append(rules, Rule{
			Pair:      parts[0],
			Character: rune(parts[1][0]),
		})
	}

	return rules
}

type Polymer struct {
	Characters map[rune]int
	Pairs      map[string]int
}

type Modification func(*Polymer)

func (p *Polymer) Execute(rules []Rule) {
	modifications := []Modification{}

	for _, rule := range rules {
		modifications = append(modifications, p.executeRule(rule))
	}

	for _, modification := range modifications {
		modification(p)
	}
}

func (p *Polymer) executeRule(rule Rule) Modification {
	pairs := p.Pairs[rule.Pair]
	if pairs == 0 {
		return func(p *Polymer) {}
	}
	return func(p *Polymer) {
		// insert character as many times as there are pairs
		p.Characters[rule.Character] += pairs
		// split all pairs (as we have inserted a new character in between)
		p.Pairs[rule.Pair] -= pairs
		// record new pairs
		p.Pairs[string(rule.Pair[0])+string(rule.Character)] += pairs
		p.Pairs[string(rule.Character)+string(rule.Pair[1])] += pairs
	}
}

func (p *Polymer) MaxChar() int {
	max := -1
	for _, count := range p.Characters {
		if count > max {
			max = count
		}
	}
	return max
}

func (p *Polymer) MinChar() int {
	min := math.MaxInt
	for _, count := range p.Characters {
		if count > 0 && count < min {
			min = count
		}
	}
	return min
}

func (p *Polymer) String() string {
	lines := []string{}

	lines = append(lines, "characters:\n")
	for c, count := range p.Characters {
		lines = append(lines, fmt.Sprintf("%s: %d", string(c), count))
	}

	lines = append(lines, "\npairs:\n")
	for p, count := range p.Pairs {
		lines = append(lines, fmt.Sprintf("%s: %d", p, count))
	}

	return strings.Join(lines, "\n")
}

type Rule struct {
	Pair      string
	Character rune
}

func (r Rule) String() string {
	return fmt.Sprintf("%s -> %s", r.Pair, string(r.Character))
}
