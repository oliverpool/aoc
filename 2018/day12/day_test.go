package day_test

import (
	"fmt"
	"strings"
	"testing"
)

/*
var in = inputTest

/*/
var in = input

//*/

func parseTransitions(lines []string) [32]bool {
	var transitions [32]bool

	for _, l := range lines {
		if l[9] == '.' {
			continue
		}
		key := 0
		for c := 0; c < 5; c++ {
			key *= 2
			if l[c] == '#' {
				key += 1
			}
		}

		transitions[key] = true
	}

	return transitions
}

func printState(state []bool) {
	nb := 0
	for _, s := range state {
		if s {
			nb += 1
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println(nb)
}

func process(generations int, input string) int {
	lines := strings.Split(input, "\n")
	transitions := parseTransitions(lines[2:])
	var state []bool
	for _, c := range lines[0][15:] {
		state = append(state, c == '#')
	}
	// fmt.Println(transitions, state, 1<<5)
	// printState(state)
	// /*
	accum := 0
	skipped := 0
	for g := 0; g < generations; g++ {
		newState := make([]bool, 0, len(state))
		state = append(state, false, false, false, false, false, false)
		lastFalse := 0
		for _, c := range state {
			accum = accum % (1 << 4)
			accum *= 2
			if c {
				accum += 1
			}
			if !transitions[accum] && len(newState) == 0 {
				skipped += 1
				continue
			}
			newState = append(newState, transitions[accum])

			if !transitions[accum] {
				lastFalse += 1
			} else {
				lastFalse = 0
			}
		}
		state = newState[:len(newState)-lastFalse]

		if g%1000000 == 0 {
			fmt.Println(50000000000 - g)
			printState(state)
		}
	}
	// printState(state)

	sum := 0
	for i, c := range state {
		if c {
			sum += i - 2*generations + skipped
		}
	}
	return sum
}

func processFaster(generations int, input string) int {
	lines := strings.Split(input, "\n")
	transitions := parseTransitions(lines[2:])
	var state []bool
	for _, c := range lines[0][15:] {
		state = append(state, c == '#')
	}
	// fmt.Println(transitions, state, 1<<5)
	// printState(state)
	// /*
	accum := 0
	skipped := 0

	oddState := state
	evenState := make([]bool, 0, len(state))
	var newState []bool
	for g := 0; g < generations; g++ {
		if g%2 == 0 {
			newState = evenState[:0]
		} else {
			newState = oddState[:0]
		}
		state = append(state, false, false, false, false, false, false)
		lastFalse := 0
		for _, c := range state {
			accum = accum % (1 << 4)
			accum *= 2
			if c {
				accum += 1
			}
			if !transitions[accum] && len(newState) == 0 {
				skipped += 1
				continue
			}
			newState = append(newState, transitions[accum])

			if !transitions[accum] {
				lastFalse += 1
			} else {
				lastFalse = 0
			}
		}

		if g%2 == 0 {
			evenState = newState
			state = evenState[:len(newState)-lastFalse]
		} else {
			oddState = newState
			state = oddState[:len(newState)-lastFalse]
		}

		if g%1000000 == 0 {
			fmt.Println(50000000000 - g)
			printState(state)

			sum := 0
			for i, c := range state {
				if c {
					sum += i - 2*(g+1) + skipped
				}
			}
			fmt.Println(g+1, sum)
		}
	}

	// printState(state)

	sum := 0
	for i, c := range state {
		if c {
			sum += i - 2*generations + skipped
		}
	}
	return sum

}

func TestOne(t *testing.T) {
	resTest := processFaster(20, inputTest)
	if resTest != 325 {
		t.Fatalf("expected 325, got: %d", resTest)
	}
	t.Log(processFaster(20, input))
	t.Fatal(processFaster(50000000000, input))
	t.Fatal("doto")
}

var inputTest = `initial state: #..#.#..##......###...###

...## => #
..#.. => #
.#... => #
.#.#. => #
.#.## => #
.##.. => #
.#### => #
#.#.# => #
#.### => #
##.#. => #
##.## => #
###.. => #
###.# => #
####. => #`

var input = `initial state: #..####.##..#.##.#..#.....##..#.###.#..###....##.##.#.#....#.##.####.#..##.###.#.......#............

##... => .
##.## => .
.#.#. => #
#..#. => .
#.### => #
.###. => .
#.#.. => .
##..# => .
..... => .
...#. => .
.#..# => .
####. => #
...## => #
..### => #
#.#.# => #
###.# => #
#...# => #
..#.# => .
.##.. => #
.#... => #
.##.# => #
.#### => .
.#.## => .
..##. => .
##.#. => .
#.##. => .
#..## => .
###.. => .
....# => .
##### => #
#.... => .
..#.. => #`
