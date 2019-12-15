package d06

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func orbitSum(r io.Reader) (int, error) {
	pairs := make(map[string]string)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()

		parts := strings.SplitN(text, ")", 2)
		pairs[parts[1]] = parts[0]
	}

	depth := map[string]int{"COM": 0}
	for len(pairs) > 0 {
		for orbiter, orbited := range pairs {
			d, ok := depth[orbited]
			if !ok {
				continue
			}
			depth[orbiter] = d + 1
			delete(pairs, orbiter)
		}
	}

	s := 0
	for _, d := range depth {
		s += d
	}

	return s, scanner.Err()
}

func TestFirstSample(t *testing.T) {
	a := assert.New(t)
	input := strings.NewReader(`COM)B
C)D
B)C
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L`)
	s, err := orbitSum(input)
	a.NoError(err)
	a.Equal(42, s)
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	s, err := orbitSum(f)
	a.NoError(err)
	a.Equal(294191, s)
}

func orbitMinTransfer(r io.Reader) (int, error) {
	pairs := make(map[string]string)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()

		parts := strings.SplitN(text, ")", 2)
		pairs[parts[1]] = parts[0]
	}

	// find all YOU ancestors
	current := pairs["YOU"]
	youAncestors := map[string]int{current: 0}
	for current != "COM" {
		next := pairs[current]
		youAncestors[next] = youAncestors[current] + 1
		current = next
	}

	// find all SAN ancestors
	current = pairs["SAN"]
	sanAncestors := map[string]int{current: 0}
	for current != "COM" {
		if youDepth, ok := youAncestors[current]; ok {
			return youDepth + sanAncestors[current], nil
		}
		next := pairs[current]
		sanAncestors[next] = sanAncestors[current] + 1
		current = next
	}

	return 0, fmt.Errorf("no possible transfer")
}

func TestSecondSample(t *testing.T) {
	a := assert.New(t)
	input := strings.NewReader(`COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
K)YOU
I)SAN`)
	s, err := orbitMinTransfer(input)
	a.NoError(err)
	a.Equal(4, s)
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	s, err := orbitMinTransfer(f)
	a.NoError(err)
	a.Equal(424, s)
}
