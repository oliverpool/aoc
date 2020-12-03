package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#`

func TestFirst(t *testing.T) {
	forest := parseInput(example)
	require.Equal(t, 7, treeEncounters(forest, 3, 1))

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	forest = parseInput(string(input))

	require.Equal(t, 178, treeEncounters(forest, 3, 1))
}

func TestSecond(t *testing.T) {
	forest := parseInput(example)
	require.Equal(t, 336, treeEncountersProduct(forest))

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	forest = parseInput(string(input))

	require.Equal(t, 3492520200, treeEncountersProduct(forest))
}
func parseInput(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n")
}
