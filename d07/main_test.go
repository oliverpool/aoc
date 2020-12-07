package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `light red bags contain 1 bright white bag, 2 muted yellow bags.
dark orange bags contain 3 bright white bags, 4 muted yellow bags.
bright white bags contain 1 shiny gold bag.
muted yellow bags contain 2 shiny gold bags, 9 faded blue bags.
shiny gold bags contain 1 dark olive bag, 2 vibrant plum bags.
dark olive bags contain 3 faded blue bags, 4 dotted black bags.
vibrant plum bags contain 5 faded blue bags, 6 dotted black bags.
faded blue bags contain no other bags.
dotted black bags contain no other bags.`

func TestFirst(t *testing.T) {
	bags := parseInput(example)
	s := validOutermost("shiny gold", bags)
	require.Equal(t, 4, s)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	bags = parseInput(string(input))
	s = validOutermost("shiny gold", bags)
	require.Equal(t, 192, s)
}

func TestSecond(t *testing.T) {
}

func parseInput(s string) map[string][]string {
	contains := make(map[string][]string)
	for _, l := range strings.Split(strings.TrimSpace(s), "\n") {
		parts := strings.SplitN(l, " contain ", 2)
		container := strings.TrimSuffix(parts[0], " bags")
		for _, contained := range strings.Split(parts[1], ", ") {
			if contained == "no other bags." {
				continue
			}
			contained = strings.TrimSuffix(contained, ".")
			contained = strings.TrimSuffix(contained, "s")
			contained = strings.TrimSuffix(contained, " bag")
			contained = strings.Join(strings.Split(contained, " ")[1:], " ")
			contains[contained] = append(contains[contained], container)
		}
	}
	return contains
}
