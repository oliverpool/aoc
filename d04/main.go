package main

import "fmt"

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
		fmt.Println(field, p[field])
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
