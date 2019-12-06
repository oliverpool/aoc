package d06

import (
	"bufio"
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
