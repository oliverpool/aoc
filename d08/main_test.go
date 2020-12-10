package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `nop +0
acc +1
jmp +4
acc +3
jmp -3
acc -99
acc +1
jmp -4
acc +6`

func TestFirst(t *testing.T) {
	prog := parseInput(example)
	require.Equal(t, 9, len(prog))
	acc, i := executeUntilLoop(prog)
	require.Equal(t, 5, acc)
	require.Equal(t, 1, i)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	prog = parseInput(string(input))
	acc, i = executeUntilLoop(prog)
	require.Equal(t, 2025, acc)
}

func executeUntilLoop(prog map[int]instruction) (int, int) {
	acc := 0
	i := 0
	d := 0
	visited := make(map[int]bool)
	for !visited[i] {
		visited[i] = true
		acc, d = prog[i](acc)
		i += d

	}
	return acc, i
}

type instruction func(ac int) (acc int, offset int)

func parseInput(s string) map[int]instruction {
	program := make(map[int]instruction)

	for i, l := range strings.Split(strings.TrimSpace(s), "\n") {
		program[i] = parseInstruction(l)
	}
	return program
}

func parseInstruction(s string) instruction {
	parts := strings.SplitN(s, " ", 2)
	v, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	switch parts[0] {
	case "nop":
		return func(acc int) (int, int) {
			return acc, 1
		}
	case "jmp":
		return func(acc int) (int, int) {
			return acc, v
		}
	case "acc":
		return func(acc int) (int, int) {
			return acc + v, 1
		}
	default:
		panic("unsupported instruction " + parts[0])
	}
}
