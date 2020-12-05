package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirst(t *testing.T) {
	cc := []struct {
		input string
		id    int64
	}{
		{"BFFFBBFRRR", 567},
		{"FFFBBBFRRR", 119},
		{"BBFFBBFRLL", 820},
	}
	for _, c := range cc {
		t.Run(c.input, func(t *testing.T) {
			require.Equal(t, c.id, seatID(c.input))
		})
	}

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	ss := parseInput(string(input))
	require.Equal(t, int64(928), getMax(ss))
}

func TestSecond(t *testing.T) {
	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	ss := parseInput(string(input))

	require.Equal(t, int64(610), getEmpty(ss))
}

func getEmpty(ss []string) int64 {
	occupied := map[int64]bool{}
	for _, s := range ss {
		occupied[seatID(s)] = true
	}
	for i := getMax(ss); i > 0; i-- {
		if occupied[i] {
			continue
		}
		if !occupied[i+1] || !occupied[i-1] {
			continue
		}
		return i
	}
	return -1
}

func getMax(ss []string) int64 {
	m := seatID(ss[0])
	for _, s := range ss {
		i := seatID(s)
		if i > m {
			m = i
		}
	}
	return m
}

func parseInput(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n")
}
