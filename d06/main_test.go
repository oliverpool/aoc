package main

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `abc

a
b
c

ab
ac

a
a
a
a

b`

func TestFirst(t *testing.T) {
	groups := parseInput(example)
	require.Len(t, groups, 5)
	s := sumAnyoneQuestions(groups)
	require.Equal(t, 11, s)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	groups = parseInput(string(input))
	s = sumAnyoneQuestions(groups)
	require.Equal(t, 6170, s)
}

func TestSecond(t *testing.T) {
	groups := parseInput(example)
	require.Len(t, groups, 5)
	s := sumEveryoneQuestions(groups)
	require.Equal(t, 6, s)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	groups = parseInput(string(input))
	s = sumEveryoneQuestions(groups)
	require.Equal(t, 2947, s)
}

func sumEveryoneQuestions(groups []string) int {
	s := 0
	for _, g := range groups {
		m, c := uniqueQuestions(g)
		for _, q := range m {
			if c == q {
				s++
			}
		}
	}
	return s
}
func sumAnyoneQuestions(groups []string) int {
	s := 0
	for _, g := range groups {
		m, _ := uniqueQuestions(g)
		s += len(m)
	}
	return s
}

func uniqueQuestions(s string) (map[byte]int, int) {
	c := 1
	m := make(map[byte]int)
	for _, b := range s {
		if b == '\n' {
			c++
			continue
		}
		m[byte(b)]++
	}
	return m, c
}

func parseInput(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n\n")
}
