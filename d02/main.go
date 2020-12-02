package main

import (
	"strings"
)

type passwordRow struct {
	Min      int
	Max      int
	Letter   string
	Password string
}

func (p passwordRow) IsValid() bool {
	c := strings.Count(p.Password, p.Letter)
	return c >= p.Min && c <= p.Max
}

func countValid(passwords []passwordRow) int {
	c := 0
	for _, v := range passwords {
		if v.IsValid() {
			c++
		}
	}
	return c
}

func (p passwordRow) IsValidTwo() bool {
	l := p.Letter[0]
	first := p.Password[p.Min-1] == l
	second := p.Password[p.Max-1] == l
	return first != second
}

func countValidTwo(passwords []passwordRow) int {
	c := 0
	for _, v := range passwords {
		if v.IsValidTwo() {
			c++
		}
	}
	return c
}
