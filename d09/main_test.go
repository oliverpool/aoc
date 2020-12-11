package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `35
20
15
25
47
40
62
55
65
95
102
117
150
182
127
219
299
277
309
576`

func TestFirst(t *testing.T) {
	list := parseInput(example)
	invalid := firstInvalid(list, 5)
	require.Equal(t, 127, invalid)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	list = parseInput(string(input))
	invalid = firstInvalid(list, 25)
	require.Equal(t, 776203571, invalid)
}

func firstInvalid(list []int, preamble int) int {
	valid := make(map[int]int, preamble*(preamble-1)/2)
	for i, v := range list[:preamble] {
		for _, w := range list[i+1 : preamble] {
			fmt.Println(v, w)
			valid[v+w]++
		}
	}
	for i, v := range list[preamble:] {
		if valid[v] == 0 {
			return v
		}

		x := list[i]
		for _, w := range list[i+1 : i+preamble] {
			valid[x+w]--
			valid[v+w]++
		}
	}
	return -1
}

func parseInput(s string) []int {
	ss := strings.Split(strings.TrimSpace(s), "\n")
	ii := make([]int, 0, len(ss))
	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		ii = append(ii, i)
	}
	return ii
}
