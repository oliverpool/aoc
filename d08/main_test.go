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
	acc, i := 0, 0
	visited := make(map[int]bool)
	for !visited[i] {
		visited[i] = true
		acc, i = prog[i].execute(acc, i)
	}
	return acc, i
}

func TestSecond(t *testing.T) {
	prog := parseInput(example)
	require.Equal(t, 9, len(prog))
	acc, i := executeUntilLoopWithSwap(prog)
	require.Equal(t, 8, acc)
	require.Equal(t, 9, i)

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	prog = parseInput(string(input))
	acc, i = executeUntilLoopWithSwap(prog)
	require.Equal(t, 2001, acc)
}

func executeUntilLoopWithSwap(prog map[int]instruction) (int, int) {
	acc, i := 0, 0
	visited := make(map[int]bool)
	for !visited[i] {
		visited[i] = true
		inst := prog[i]
		switch in := inst.(type) {
		case jmp:
			facc, fi := nop(int(in)).execute(acc, i)
			facc, fi = continueExecution(prog, facc, fi, visited)
			if fi == len(prog) {
				return facc, fi
			}
		case nop:
			facc, fi := jmp(int(in)).execute(acc, i)
			facc, fi = continueExecution(prog, facc, fi, visited)
			if fi == len(prog) {
				return facc, fi
			}
		}
		acc, i = inst.execute(acc, i)
	}
	return acc, i
}

func continueExecution(prog map[int]instruction, acc, i int, visitedBefore map[int]bool) (int, int) {
	visited := make(map[int]bool)
	for !visited[i] && !visitedBefore[i] {
		visited[i] = true
		if i == len(prog) {
			return acc, i
		}
		acc, i = prog[i].execute(acc, i)
	}
	return acc, i
}

func parseInput(s string) map[int]instruction {
	program := make(map[int]instruction)

	for i, l := range strings.Split(strings.TrimSpace(s), "\n") {
		program[i] = parseInstruction(l)
	}
	return program
}

type instruction interface {
	execute(acc, i int) (int, int)
}

type nop int

func (nop) execute(acc, i int) (int, int) {
	return acc, i + 1
}

type jmp int

func (j jmp) execute(acc, i int) (int, int) {
	return acc, i + int(j)
}

type acc int

func (a acc) execute(acc, i int) (int, int) {
	return acc + int(a), i + 1
}

func parseInstruction(s string) instruction {
	parts := strings.SplitN(s, " ", 2)
	v, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	switch parts[0] {
	case "nop":
		return nop(v)
	case "jmp":
		return jmp(v)
	case "acc":
		return acc(v)
	default:
		panic("unsupported instruction " + parts[0])
	}
}
