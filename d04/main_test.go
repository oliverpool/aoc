package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const example = `ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in`

func TestFirst(t *testing.T) {
	batch := parseInput(example)
	require.Equal(t, 2, countValidPassports(batch))

	input, err := ioutil.ReadFile("./input")
	require.NoError(t, err)
	batch = parseInput(string(input))

	require.Equal(t, 237, countValidPassports(batch))
}

func TestSecond(t *testing.T) {
	// forest := parseInput(example)
	// require.Equal(t, 336, treeEncountersProduct(forest))

	// input, err := ioutil.ReadFile("./input")
	// require.NoError(t, err)
	// forest = parseInput(string(input))

	// require.Equal(t, 3492520200, treeEncountersProduct(forest))
}
func parseInput(s string) []passport {
	ss := strings.Split(strings.TrimSpace(s), "\n\n")
	passports := make([]passport, 0, len(ss))
	for _, s := range ss {
		fields := strings.Split(strings.ReplaceAll(s, "\n", " "), " ")
		p := make(passport, len(fields))
		for _, f := range fields {
			kv := strings.SplitN(f, ":", 2)
			fmt.Println(kv)
			p[kv[0]] = kv[1]
		}
		passports = append(passports, p)
	}
	return passports
}
