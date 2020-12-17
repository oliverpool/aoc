package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `class: 1-3 or 5-7
row: 6-11 or 33-44
seat: 13-40 or 45-50

your ticket:
7,1,14

nearby tickets:
7,3,47
40,4,50
55,2,20
38,6,12`

type intervals []interval

func (i intervals) isValid(n int) bool {
	for _, r := range i {
		if r.isValid(n) {
			return true
		}
	}
	return false
}

type interval struct {
	min, max int
}

func (i interval) isValid(n int) bool {
	return i.min <= n && n <= i.max
}

type tickets struct {
	rules  map[string]intervals
	my     []int
	nearby [][]int
}

func parseInput(s string) tickets {
	parts := strings.Split(strings.TrimSpace(s), "\n\n")
	return tickets{
		rules:  parseRules(parts[0]),
		my:     parseMy(parts[1]),
		nearby: parseNearby(parts[2]),
	}
}

func parseRules(s string) map[string]intervals {
	parts := strings.Split(strings.TrimSpace(s), "\n")
	rules := make(map[string]intervals, len(parts))
	for _, p := range parts {
		sparts := strings.SplitN(p, ": ", 2)
		for _, ssp := range strings.Split(sparts[1], " or ") {
			rules[sparts[0]] = append(rules[sparts[0]], parseInterval(ssp))
		}
	}
	return rules
}
func parseInterval(s string) interval {
	parts := strings.SplitN(s, "-", 2)
	min, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return interval{min, max}
}

func parseMy(s string) []int {
	parts := strings.Split(strings.TrimSpace(s), "\n")
	return parseTickets(parts[1])
}

func parseNearby(s string) [][]int {
	parts := strings.Split(strings.TrimSpace(s), "\n")
	tt := make([][]int, 0, len(parts)-1)
	for _, p := range parts[1:] {
		tt = append(tt, parseTickets(p))
	}
	return tt
}

func parseTickets(s string) []int {
	parts := strings.Split(s, ",")
	tt := make([]int, 0, len(parts))
	for _, p := range parts {
		t, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		tt = append(tt, t)
	}
	return tt
}

func (t tickets) checkNearby() int {
	sum := 0
	for _, nn := range t.nearby {
		for _, n := range nn {
			if !t.isValid(n) {
				sum += n
			}
		}
	}
	return sum
}
func (t tickets) isValid(n int) bool {
	for _, rr := range t.rules {
		if rr.isValid(n) {
			return true
		}
	}
	return false
}

func TestFirst(t *testing.T) {
	tickets := parseInput(example)
	invalid := tickets.checkNearby()
	require.Equal(t, 71, invalid)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	tickets = parseInput(string(input))
	invalid = tickets.checkNearby()
	require.Equal(t, 29851, invalid)
}

const example2 = `class: 0-1 or 4-19
row: 0-5 or 8-19
seat: 0-13 or 16-19

your ticket:
11,12,13

nearby tickets:
3,9,18
15,1,5
5,14,9`

func (t tickets) checkMy(prefix string) int {
	valid := make(map[string]map[int]bool)
	for r, rr := range t.rules {
		valid[r] = make(map[int]bool, len(t.my))
		for i, v := range t.my {
			if rr.isValid(v) {
				valid[r][i] = true
			}
		}
	}
	for _, nn := range t.nearby {
		invalids := t.invalids(nn)
		for i, names := range invalids {
			for _, name := range names {
				delete(valid[name], i)
			}
		}
	}

	// fmt.Println(valid)
	positions := positions(valid)
	// fmt.Println(positions)
	prod := 1
	for i, name := range positions {
		if strings.HasPrefix(name, prefix) {
			prod *= t.my[i]
		}
	}
	return prod
}

func positions(valid map[string]map[int]bool) map[int]string {
	positions := make(map[int]string)
	l := len(valid)
	for len(positions) != l {
		for name, possibilites := range valid {
			if len(possibilites) == 1 {
				for i := range possibilites {
					positions[i] = name
				}
				delete(valid, name)
				continue
			}
			for i := range possibilites {
				if positions[i] != "" {
					delete(valid[name], i)
				}
			}
		}
	}
	return positions
}

func (t tickets) invalids(nn []int) [][]string {
	var invalids [][]string
	for _, n := range nn {
		var names []string
		for name, intervals := range t.rules {
			if !intervals.isValid(n) {
				names = append(names, name)
			}
		}
		if len(names) == len(t.rules) {
			return nil
		}
		invalids = append(invalids, names)
	}
	return invalids
}

func TestSecond(t *testing.T) {
	tickets := parseInput(example2)
	invalid := tickets.checkMy("ro")
	require.Equal(t, 11, invalid)
	invalid = tickets.checkMy("")
	require.Equal(t, 11*12*13, invalid)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	tickets = parseInput(string(input))
	invalid = tickets.checkMy("departure")
	require.Equal(t, 3029180675981, invalid)
}
