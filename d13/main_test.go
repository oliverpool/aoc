package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `939
7,13,x,x,59,x,31,19`

func parseInput(s string) (int64, []int64) {
	parts := strings.Split(strings.TrimSpace(s), "\n")
	ts, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		panic(err)
	}

	var buses []int64
	for _, l := range strings.Split(parts[1], ",") {
		var b int64
		if l != "x" {
			b, err = strconv.ParseInt(l, 10, 64)
		}
		buses = append(buses, b)
	}
	return ts, buses
}

func findNext(ts int64, buses []int64) (id int64, delay int64) {
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
	require.Equal(t, int64(59), id)
	require.Equal(t, int64(5), delay)
	require.Equal(t, int64(295), id*delay)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)

	ts, buses = parseInput(string(input))
	id, delay = findNext(ts, buses)
	require.Equal(t, int64(3246), id*delay)
}

func euclidex(a, b int64) (int64, int64, int64) {
	r, u, v, r2, u2, v2 := a, int64(1), int64(0), b, int64(0), int64(1)
	for r2 != 0 {
		q := r / r2
		r, u, v, r2, u2, v2 = r2, u2, v2, r-q*r2, u-q*u2, v-q*v2
	}
	//  r = pgcd(a, b) et r = a*u+b*v
	return r, u, v
}

func findFirst(buses []int64) int64 {
	var n, offset int64

	for i, b := range buses {
		if b == 0 {
			continue
		}
		if n == 0 {
			n = b
			offset = 0
			if i != 0 {
				panic("unsupported")
			}
			continue
		}
		gcd, u, v := euclidex(n, b)
		pcm := b * (n / gcd)
		// fmt.Println(b, n, pcm)
		// gcd = n*u+b*v
		// fmt.Printf("%d = %d x %d  +  %d x %d\n", gcd, u, n, v, b)

		goal := offset + int64(i)
		if goal%gcd != 0 {
			fmt.Println(goal, gcd, goal/gcd, v)
			panic("impossible")
		}
		// n*X + offset + i = Y*b
		// offset2 = n*X + offset
		// fmt.Println(goal, gcd, u, n, "=", goal*u*(n/gcd))

		dBig := big.NewInt(-goal)
		dBig = dBig.Mul(dBig, big.NewInt(u))
		dBig = dBig.Mul(dBig, big.NewInt(n/gcd))
		var m big.Int
		dBig.DivMod(dBig, big.NewInt(pcm), &m)

		// d := (-goal * u * (n / gcd)) % pcm
		// for d < 0 {
		// 	d += pcm
		// }
		offset += m.Int64()

		n = pcm
	}
	// fmt.Println(n)
	return offset
}

func TestSecond(t *testing.T) {
	cc := []struct {
		buses    []int64
		expected int64
	}{
		{[]int64{7, 8}, 7},
		{[]int64{7, 13, 79}, 77},

		{[]int64{17, 0, 13, 19}, 3417},
		{[]int64{67, 7, 59, 61}, 754018},
		{[]int64{67, 0, 7, 59, 61}, 779210},
		{[]int64{67, 7, 0, 59, 61}, 1261476},
		{[]int64{1789, 37, 47, 1889}, 1202161486},
	}
	for _, c := range cc {
		t.Run(fmt.Sprint(c.buses), func(t *testing.T) {
			// fmt.Println(c.buses)
			ts := findFirst(c.buses)
			require.Equal(t, c.expected, ts)
		})
	}

	_, buses := parseInput(example)
	ts := findFirst(buses)
	require.Equal(t, int64(1068781), ts)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)

	ts, buses = parseInput(string(input))
	ts = findFirst(buses)
	require.Equal(t, int64(1010182346291467), ts)
}
