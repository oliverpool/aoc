package d04

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isIncreasing(n int) bool {
	previous := 9
	var last int
	for n > 0 {
		last, n = (n % 10), n/10
		if last > previous {
			return false
		}
		previous = last
	}
	return true
}

func hasDouble(n int) bool {
	previous := -1
	var last int
	for n > 0 {
		last, n = (n % 10), n/10
		if last == previous {
			return true
		}
		previous = last
	}
	return false
}

func hasDoubleNoTriple(n int) bool {
	previous := -1
	var last int
	for n > 0 {
		last, n = (n % 10), n/10
		if last == previous {
			if n%10 == last {
				for n%10 == last {
					n = n / 10
				}
				continue
			}
			return true
		}
		previous = last
	}
	return false
}

func countValid(start, end int) int {
	c := 0
	for i := start; i <= end; i++ {
		if hasDouble(i) && isIncreasing(i) {
			c++
		}
	}
	return c
}

func countValidNoTriple(start, end int) int {
	c := 0
	for i := start; i <= end; i++ {
		if hasDoubleNoTriple(i) && isIncreasing(i) {
			c++
		}
	}
	return c
}

func TestIsIncreasing(t *testing.T) {
	cc := []struct {
		n          int
		increasing bool
	}{
		{1, true},
		{12, true},
		{11, true},
		{21, false},
		{22, true},
		{123445, true},
		{1234450, false},
	}
	for _, c := range cc {
		t.Run(fmt.Sprintf("%d", c.n), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(c.increasing, isIncreasing(c.n))
		})
	}
}

func TestHasDouble(t *testing.T) {
	cc := []struct {
		n      int
		double bool
	}{
		{1, false},
		{12, false},
		{11, true},
		{21, false},
		{22, true},
		{123445, true},
		{1234450, true},
	}
	for _, c := range cc {
		t.Run(fmt.Sprintf("%d", c.n), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(c.double, hasDouble(c.n))
		})
	}
}

func TestHasDoubleNoTriple(t *testing.T) {
	cc := []struct {
		n      int
		double bool
	}{
		{1, false},
		{12, false},
		{11, true},
		{21, false},
		{22, true},
		{123445, true},
		{12344450, false},
		{11122, true},
		{111222, false},
		{11222, true},
	}
	for _, c := range cc {
		t.Run(fmt.Sprintf("%d", c.n), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(c.double, hasDoubleNoTriple(c.n))
		})
	}
}

func TestFirstBrute(t *testing.T) {
	a := assert.New(t)
	a.Equal(2814, countValid(109165, 576723))
}

func TestSecondBrute(t *testing.T) {
	a := assert.New(t)
	a.Equal(1991, countValidNoTriple(109165, 576723))
}
