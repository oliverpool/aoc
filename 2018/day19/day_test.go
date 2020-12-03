package day19

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"
)

type registers [6]int

type instruction struct {
	a int
	b int
	c int
}

type sample struct {
	opcode      int
	before      registers
	instruction instruction
	after       registers
}

type opcode func(instruction, registers) registers

var opcodes = map[string]opcode{
	"addr": func(i instruction, in registers) registers {
		in[i.c] = in[i.a] + in[i.b]
		return in
	},
	"addi": func(i instruction, in registers) registers {
		in[i.c] = in[i.a] + i.b
		return in
	},
	"mulr": func(i instruction, in registers) registers {
		in[i.c] = in[i.a] * in[i.b]
		return in
	},
	"muli": func(i instruction, in registers) registers {
		in[i.c] = in[i.a] * i.b
		return in
	},
	"banr": func(i instruction, in registers) registers {
		in[i.c] = in[i.a] & in[i.b]
		return in
	},
	"bani": func(i instruction, in registers) registers {
		in[i.c] = in[i.a] & i.b
		return in
	},
	"borr": func(i instruction, in registers) registers {
		in[i.c] = in[i.a] | in[i.b]
		return in
	},
	"bori": func(i instruction, in registers) registers {
		in[i.c] = in[i.a] | i.b
		return in
	},
	"setr": func(i instruction, in registers) registers {
		in[i.c] = in[i.a]
		return in
	},
	"seti": func(i instruction, in registers) registers {
		in[i.c] = i.a
		return in
	},
	"gtir": func(i instruction, in registers) registers {
		if i.a > in[i.b] {
			in[i.c] = 1
		} else {
			in[i.c] = 0
		}
		return in
	},
	"gtri": func(i instruction, in registers) registers {
		if in[i.a] > i.b {
			in[i.c] = 1
		} else {
			in[i.c] = 0
		}
		return in
	},
	"gtrr": func(i instruction, in registers) registers {
		if in[i.a] > in[i.b] {
			in[i.c] = 1
		} else {
			in[i.c] = 0
		}
		return in
	},
	"eqir": func(i instruction, in registers) registers {
		if i.a == in[i.b] {
			in[i.c] = 1
		} else {
			in[i.c] = 0
		}
		return in
	},
	"eqri": func(i instruction, in registers) registers {
		if in[i.a] == i.b {
			in[i.c] = 1
		} else {
			in[i.c] = 0
		}
		return in
	},
	"eqrr": func(i instruction, in registers) registers {
		if in[i.a] == in[i.b] {
			in[i.c] = 1
		} else {
			in[i.c] = 0
		}
		return in
	},
}

func runInstructions(t *testing.T, in io.Reader, opcodes map[string]opcode) registers {
	sc := bufio.NewScanner(in)

	var reg registers

	sc.Scan() // first line
	line := sc.Text()
	var inst int
	_, err := fmt.Sscanf(line, "#ip %d", &inst)
	if err != nil {
		t.Errorf("could not scan: %v", err)
	}
	fmt.Println(inst)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		var current instruction
		var op string
		_, err := fmt.Sscanf(line, "%s %d %d %d", &op, &current.a, &current.b, &current.c)
		if err != nil {
			t.Errorf("could not scan: %v", err)
		}
		reg = opcodes[op](current, reg)
	}

	if sc.Err() != nil {
		t.Errorf("could not scan: %v", sc.Err())
	}
	return reg
}

func TestFirstSample(t *testing.T) {
	// a := assert.New(t)
	sample := strings.NewReader(`#ip 0
	seti 5 0 1
	seti 6 0 2
	addi 0 1 0
	addr 1 2 3
	setr 1 0 0
	seti 8 0 4
	seti 9 0 5`)

	runInstructions(t, sample, opcodes)
}
