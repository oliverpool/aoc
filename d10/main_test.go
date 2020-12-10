package main

import (
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `28
33
18
42
31
14
46
20
48
47
24
23
49
45
19
38
39
11
1
32
25
35
8
17
7
9
4
2
34
10
3`

func TestHelloWorld(t *testing.T) {
	adapters := parseInput(example)
	require.Len(t, adapters, 31)

	steps := countSteps(adapters)
	steps[3]++
	p := steps[1] * steps[3]
	require.Equal(t, 220, p)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	adapters = parseInput(string(input))

	steps = countSteps(adapters)
	steps[3]++
	p = steps[1] * steps[3]
	require.Equal(t, 2516, p)
}

func parseInput(s string) []int {
	nn := strings.Split(strings.TrimSpace(s), "\n")
	out := make([]int, 0, len(nn))
	for _, n := range nn {
		i, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		out = append(out, i)
	}
	sort.Ints(out)
	return out
}
