package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type passport map[string]string

func (p passport) IsValid() bool {
	fields := []string{"byr",
		"iyr",
		"eyr",
		"hgt",
		"hcl",
		"ecl",
		"pid",
		//		"cid",
	}
	for _, field := range fields {
		if p[field] == "" {
			return false
		}
	}
	return true

}

func countValidPassports(pp []passport) int {
	valid := 0
	for _, p := range pp {
		if p.IsValid() {
			valid++
		}
	}
	return valid
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
func (p passport) IsValidStrict() bool {
	fields := map[string]func(string) bool{
		"byr": func(s string) bool {
			i := parseInt(s)
			return i >= 1920 && i <= 2002
		},
		"iyr": func(s string) bool {
			i := parseInt(s)
			return i >= 2010 && i <= 2020
		},
		"eyr": func(s string) bool {
			i := parseInt(s)
			return i >= 2020 && i <= 2030
		},
		"hgt": func(s string) bool {
			if len(s) < 3 {
				return false
			}
			unit := s[len(s)-2:]
			i := parseInt(s[:len(s)-2])
			fmt.Println(unit, i)
			if unit == "cm" {
				return i >= 150 && i <= 193
			}
			if unit == "in" {
				return i >= 59 && i <= 76
			}
			return false
		},
		"hcl": func(s string) bool {
			m, err := regexp.MatchString(`^#[0-9a-f]{6}$`, s)
			if err != nil {
				panic(err)
			}
			return m
		},
		"ecl": func(s string) bool {
			for _, e := range []string{
				"amb", "blu", "brn", "gry", "grn", "hzl", "oth",
			} {
				if s == e {
					return true
				}
			}
			return false
		},
		"pid": func(s string) bool {
			m, err := regexp.MatchString(`^[0-9]{9}$`, s)
			if err != nil {
				panic(err)
			}
			return m
		},
		//		"cid",
	}
	for field, valid := range fields {
		if valid == nil {
			continue
		}
		if !valid(p[field]) {
			return false
		}
	}
	return true

}
func countValidStrictPassports(pp []passport) int {
	valid := 0
	for _, p := range pp {
		if p.IsValidStrict() {
			valid++
		}
	}
	return valid
}
