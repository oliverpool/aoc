package main

import (
	"bufio"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func requiredFuel(mass int) int {
	return mass/3 - 2
}

func requiredFuelCompenssated(mass int) int {
	fuel := requiredFuel(mass)
	sum := 0
	for fuel > 0 {
		sum += fuel
		fuel = requiredFuel(fuel)
	}
	return sum
}

func TestRequiredFuel(t *testing.T) {
	cc := []struct {
		mass int
		fuel int
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}
	for _, c := range cc {
		t.Run("Mass "+strconv.Itoa(c.mass), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(c.fuel, requiredFuel(c.mass))
		})
	}
}

func TestFirst(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	sum := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		mass, err := strconv.Atoi(text)
		a.NoError(err)
		sum += requiredFuel(mass)
	}
	a.NoError(scanner.Err())
	a.Equal(3490763, sum)
}

func TestRequiredFuelCompenssated(t *testing.T) {
	cc := []struct {
		mass int
		fuel int
	}{
		{12, 2},
		{14, 2},
		{1969, 966},
		{100756, 50346},
	}
	for _, c := range cc {
		t.Run("Mass "+strconv.Itoa(c.mass), func(t *testing.T) {
			a := assert.New(t)
			a.Equal(c.fuel, requiredFuelCompenssated(c.mass))
		})
	}
}

func TestSecond(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./input")
	a.NoError(err)
	defer f.Close()

	sum := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		mass, err := strconv.Atoi(text)
		a.NoError(err)
		sum += requiredFuelCompenssated(mass)
	}
	a.NoError(scanner.Err())
	a.Equal(5233250, sum)
}
