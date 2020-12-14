package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `939
7,13,x,x,59,x,31,19`

func parseInput(s string) (int, []int) {
	parts := strings.Split(strings.TrimSpace(s), "\n")
	ts, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}

	var buses []int
	for _, l := range strings.Split(parts[1], ",") {
		b := 0
		if l != "x" {
			b, err = strconv.Atoi(l)
		}
		buses = append(buses, b)
	}
	return ts, buses
}

func findNext(ts int, buses []int) (id int, delay int) {

	for _, b := range buses {
		if b == 0 {
			continue
		}
		d := b - (ts % b)
		if delay == 0 || d < delay {
			id, delay = b, d
		}
	}

	return id, delay
}

func TestFirst(t *testing.T) {
	ts, buses := parseInput(example)
	id, delay := findNext(ts, buses)
	require.Equal(t, id, 59)
	require.Equal(t, delay, 5)
	require.Equal(t, id*delay, 295)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)

	ts, buses = parseInput(string(input))
	id, delay = findNext(ts, buses)
	require.Equal(t, id*delay, 295)
}
