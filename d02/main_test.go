package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const example = `1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`

func TestOneExample(t *testing.T) {
	input, err := parseInput(strings.NewReader(example))
	assert.NoError(t, err)

	assert.Equal(t, 2, countValid(input))
}

func TestOne(t *testing.T) {
	var input []passwordRow
	var err error
	err = open("./input", func(r io.Reader) error {
		input, err = parseInput(r)
		return err
	})
	assert.NoError(t, err)

	assert.Equal(t, 422, countValid(input))
}
func TestTwoExample(t *testing.T) {
	input, err := parseInput(strings.NewReader(example))
	assert.NoError(t, err)

	assert.Equal(t, 1, countValidTwo(input))
}

func TestTwo(t *testing.T) {
	var input []passwordRow
	var err error
	err = open("./input", func(r io.Reader) error {
		input, err = parseInput(r)
		return err
	})
	assert.NoError(t, err)

	assert.Equal(t, 422, countValidTwo(input))
}

func open(path string, cb func(io.Reader) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return cb(f)
}

func parseInput(r io.Reader) ([]passwordRow, error) {
	var input []passwordRow

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()

		var v passwordRow
		_, err := fmt.Sscanf(text, "%d-%d %s %s", &v.Min, &v.Max, &v.Letter, &v.Password)
		v.Letter = v.Letter[:1]
		if err != nil {
			return input, err
		}
		input = append(input, v)
	}
	return input, scanner.Err()
}
