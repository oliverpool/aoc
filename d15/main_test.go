package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `0,3,6`

func parseInput(s string) []int {
	var nums []int
	for _, l := range strings.Split(strings.TrimSpace(s), ",") {
		num, err := strconv.Atoi(l)
		if err != nil {
			panic(err)
		}
		nums = append(nums, num)
	}
	return nums
}

func runFor(nums []int, n int) int {
	seen := make(map[int]int)
	i, last := 0, 0
	seenBefore := 0
	for i, last = range nums {
		seenBefore = seen[last]
		seen[last] = i + 1
	}
	i++

	for ; i < n; i++ {
		// consider last
		if seenBefore == 0 {
			last = 0
		} else {
			last = i - seenBefore
		}
		seenBefore = seen[last]
		seen[last] = i + 1
	}
	return last
}

func TestFirst(t *testing.T) {
	cc := []struct {
		input    string
		after    int
		expected int
	}{
		{"0,3,6", 4, 0},
		{"0,3,6", 5, 3},
		{"0,3,6", 6, 3},
		{"0,3,6", 7, 1},
		{"0,3,6", 8, 0},
		{"0,3,6", 9, 4},
		{"0,3,6", 10, 0},
		{"0,3,6", 2020, 436},

		{"1,3,2", 2020, 1},
		{"2,1,3", 2020, 10},
		{"1,2,3", 2020, 27},
		{"2,3,1", 2020, 78},
		{"3,2,1", 2020, 438},
		{"3,1,2", 2020, 1836},

		{"0,20,7,16,1,18,15", 2020, 1025},
	}
	for _, c := range cc {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()
			nums := parseInput(c.input)
			got := runFor(nums, c.after)
			require.Equal(t, c.expected, got)
		})
	}
}

func TestSecond(t *testing.T) {
	cc := []struct {
		input    string
		after    int
		expected int
	}{

		{"0,3,6", 30000000, 175594},
		{"1,3,2", 30000000, 2578},
		{"2,1,3", 30000000, 3544142},
		{"1,2,3", 30000000, 261214},
		{"2,3,1", 30000000, 6895259},
		{"3,2,1", 30000000, 18},
		{"3,1,2", 30000000, 362},

		{"0,20,7,16,1,18,15", 30000000, 129262},
	}
	for _, c := range cc {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()
			nums := parseInput(c.input)
			got := runFor(nums, c.after)
			require.Equal(t, c.expected, got)
		})
	}
}
