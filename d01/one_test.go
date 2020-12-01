package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstCase(t *testing.T) {
	assert.Equal(t, 514579, findProduct(2020,
		1721,
		979,
		366,
		299,
		675,
		1456))
}

func TestFirst(t *testing.T) {
	var expenses []int
	var err error
	err = open("./input", func(r io.Reader) error {
		expenses, err = parseInput(r)
		return err
	})
	assert.NoError(t, err)

	assert.Equal(t, 1010299, findProduct(2020, expenses...))
}

func open(path string, cb func(io.Reader) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	return cb(f)
}

func parseInput(r io.Reader) ([]int, error) {
	var expenses []int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()

		var v int
		_, err := fmt.Sscanf(text, "%d", &v)
		if err != nil {
			return expenses, err
		}
		expenses = append(expenses, v)
	}
	return expenses, scanner.Err()
}
