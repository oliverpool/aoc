package day_test

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type registers [4]int

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

func parseOpcodes(t *testing.T, in io.Reader) []sample {
	sc := bufio.NewScanner(in)

	var current sample
	var samples []sample
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "Before") {
			fmt.Sscanf(line, "Before: [%d, %d, %d, %d]", &current.before[0], &current.before[1], &current.before[2], &current.before[3])
		} else if strings.HasPrefix(line, "After") {
			fmt.Sscanf(line, "After: [%d, %d, %d, %d]", &current.after[0], &current.after[1], &current.after[2], &current.after[3])
			samples = append(samples, current)
			current = sample{}
		} else {
			fmt.Sscanf(line, "%d %d %d %d", &current.opcode, &current.instruction.a, &current.instruction.b, &current.instruction.c)
		}
	}

	if sc.Err() != nil {
		t.Errorf("could not scan: %v", sc.Err())
	}
	return samples
}

func TestOne(t *testing.T) {

	f, err := os.Open("opcodes.txt")
	if err != nil {
		t.Fatal("could not open opcodes")
	}
	defer f.Close()
	samples := parseOpcodes(t, f)

	n := 0
	for _, s := range samples {
		valid := 0
		for _, f := range opcodes {
			if f(s.instruction, s.before) == s.after {
				valid += 1
			}
			if valid >= 3 {
				n += 1
				break
			}
		}
	}
	t.Fatal(n)
}

func runInstructions(t *testing.T, in io.Reader, opcodes map[int]opcode) registers {
	sc := bufio.NewScanner(in)

	var reg registers

	for sc.Scan() {
		line := sc.Text()
		var current instruction
		var op int
		fmt.Sscanf(line, "%d %d %d %d", &op, &current.a, &current.b, &current.c)
		reg = opcodes[op](current, reg)
	}

	if sc.Err() != nil {
		t.Errorf("could not scan: %v", sc.Err())
	}
	return reg
}

func TestTwo(t *testing.T) {
	f, err := os.Open("opcodes.txt")
	if err != nil {
		t.Fatal("could not open opcodes")
	}
	defer f.Close()
	samples := parseOpcodes(t, f)

	// reduce based on the before/after operations
	possibilites := make(map[int]map[string]opcode)
	for _, s := range samples {
		if possibilites[s.opcode] == nil {
			possibilites[s.opcode] = make(map[string]opcode, len(opcodes))
			for op, f := range opcodes {
				if f(s.instruction, s.before) == s.after {
					possibilites[s.opcode][op] = f
				}
			}
			continue
		}
		for op, f := range possibilites[s.opcode] {
			if f(s.instruction, s.before) != s.after {
				delete(possibilites[s.opcode], op)
			}
		}
	}

	matches := make(map[string]int)
	for len(possibilites) > 0 {
		for i, p := range possibilites {
			for op := range p {
				if _, ok := matches[op]; ok {
					delete(p, op)
				}
			}

			if len(p) == 1 {
				for op := range p {
					matches[op] = i
				}
				delete(possibilites, i)
			}
		}
	}

	operations := make(map[int]opcode)
	for name, i := range matches {
		operations[i] = opcodes[name]
	}

	f, err = os.Open("instructions.txt")
	if err != nil {
		t.Fatal("could not open instructions")
	}
	defer f.Close()
	res := runInstructions(t, f, operations)
	t.Fatal(res)
}
