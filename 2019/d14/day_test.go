package d14

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type chemical string

type reaction struct {
	output int
	input  map[chemical]int
}

func open(path string, cb func(io.Reader) error) error {
	f, err := os.Open("./input")
	if err != nil {
		return err
	}
	return cb(f)
}

func parseInput(r io.Reader) (map[chemical]reaction, error) {
	reactions := make(map[chemical]reaction)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {

		text := strings.SplitN(scanner.Text(), " => ", 2)
		inputs := strings.Split(text[0], ", ")

		r := reaction{
			input: make(map[chemical]int, len(inputs)),
		}

		for _, i := range inputs {
			n, che, err := parseChemical(i)
			if err != nil {
				return reactions, err
			}
			r.input[che] = n
		}

		var err error
		var che chemical
		r.output, che, err = parseChemical(text[1])
		if err != nil {
			return reactions, err
		}

		reactions[che] = r
	}
	return reactions, scanner.Err()
}

func parseChemical(s string) (int, chemical, error) {
	var n int
	var che chemical
	_, err := fmt.Sscanf(s, "%d %s", &n, &che)
	return n, che, err
}

func TestParse(t *testing.T) {
	cc := []struct {
		input                string
		count, subcount, ore int
	}{
		{
			`10 ORE => 10 A
1 ORE => 1 B
7 A, 1 B => 1 C
7 A, 1 C => 1 D
7 A, 1 D => 1 E
7 A, 1 E => 1 FUEL`,
			6,
			10,
			31,
		},
		{
			`9 ORE => 2 A
8 ORE => 3 B
7 ORE => 5 C
3 A, 4 B => 1 AB
5 B, 7 C => 1 BC
4 C, 1 A => 1 CA
2 AB, 3 BC, 4 CA => 1 FUEL`,
			7,
			12,
			165,
		},
	}
	for _, c := range cc {
		t.Run("", func(t *testing.T) {
			a := assert.New(t)
			m, err := parseInput(strings.NewReader(c.input))
			a.NoError(err)
			t.Log(m)
			a.Equal(c.count, len(m))
			subcount := 0
			for _, r := range m {
				subcount += len(r.input)
			}
			a.Equal(c.subcount, subcount)

			d := computeDepth(m)
			t.Log(d)
			a.Equal(len(m)+1, len(d)) // + ORE

			ore := reduceReaction(m, 1)
			a.Equal(c.ore, ore)
		})
	}
}

func computeDepth(r map[chemical]reaction) map[chemical]int {
	reactions := make(map[chemical]reaction, len(r))
	for k, v := range r {
		reactions[k] = v
	}

	depths := map[chemical]int{"ORE": 0}
	cDepth := func(cc map[chemical]int) int {
		dmax := 0
		for c := range cc {
			d, ok := depths[c]
			if !ok {
				return -1
			}
			if d > dmax {
				dmax = d
			}
		}
		return dmax + 1
	}
	for len(reactions) > 0 {
		for c, r := range reactions {
			d := cDepth(r.input)
			if d == -1 {
				continue
			}
			depths[c] = d
			delete(reactions, c)
		}
	}

	return depths
}

func reduceReaction(reactions map[chemical]reaction, fuels int) int {
	depths := computeDepth(reactions)
	maxDepth := func(cc map[chemical]int) int {
		dmax := 0
		for c := range cc {
			d := depths[c]
			if d > dmax {
				dmax = d
			}
		}
		return dmax
	}
	aggregate := make(map[chemical]int, len(reactions["FUEL"].input))
	for i, n := range reactions["FUEL"].input {
		aggregate[i] = fuels * n
	}
	for len(aggregate) > 1 {
		d := maxDepth(aggregate)
		for c, n := range aggregate {
			if depths[c] != d {
				continue
			}
			// need 'n' times the 'c' chemical
			r := reactions[c]
			// this reaction must happen 'factor' times
			factor := ((n - 1) / r.output) + 1

			for subc, subn := range r.input {
				aggregate[subc] += factor * subn
			}
			delete(aggregate, c)
		}
	}
	return aggregate["ORE"]
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	var reactions map[chemical]reaction
	var err error
	err = open("./input", func(r io.Reader) error {
		reactions, err = parseInput(r)
		return err
	})
	a.NoError(err)

	ore := reduceReaction(reactions, 1)
	a.Equal(1582325, ore)
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	var reactions map[chemical]reaction
	var err error
	err = open("./input", func(r io.Reader) error {
		reactions, err = parseInput(r)
		return err
	})
	a.NoError(err)

	maxOre := 1000000000000

	fuels := 1
	usedOre := reduceReaction(reactions, fuels)
	t.Log(fuels, usedOre)
	for usedOre < maxOre {
		orePerFuel := usedOre / fuels
		delta := (maxOre - usedOre) / orePerFuel
		if delta <= 0 {
			delta = 1
		}
		fuels += delta

		usedOre = reduceReaction(reactions, fuels)
		t.Log(fuels, usedOre)
	}

	for usedOre > maxOre {
		fuels--
		usedOre = reduceReaction(reactions, fuels)
	}

	a.Equal(2267486, fuels)
}
