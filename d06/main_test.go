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
	s := sumUniqueQuestions(groups)
	require.Equal(t, 11, s)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	groups = parseInput(string(input))
	s = sumUniqueQuestions(groups)
	require.Equal(t, 6170, s)
}

func TestSecond(t *testing.T) {
}

func sumUniqueQuestions(groups []string) int {
	s := 0
	for _, g := range groups {
		s += len(uniqueQuestions(g))
	}
	return s
}

func uniqueQuestions(s string) map[byte]bool {
	m := make(map[byte]bool)
	for _, b := range s {
		if b == '\n' {
			continue
		}
		m[byte(b)] = true
	}
	return m
}

func parseInput(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n\n")
}
