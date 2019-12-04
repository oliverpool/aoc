package day_test

import (
	"fmt"
	"strconv"
	"testing"
)

func generateRecipes(minLength int) []int {
	recipes := []int{3, 7}
	elves := [...]int{0, 1}
	for len(recipes) <= minLength+10 {
		score := recipes[elves[0]] + recipes[elves[1]]
		if tens := score / 10; tens > 0 {
			recipes = append(recipes, tens)
		}
		recipes = append(recipes, score%10)

		elves[0] = (elves[0] + recipes[elves[0]] + 1) % len(recipes)
		elves[1] = (elves[1] + recipes[elves[1]] + 1) % len(recipes)
	}

	return recipes
}

func takeTenAt(recipes []int, at int) string {
	res := ""
	for index := at; index < at+10; index++ {
		res += strconv.Itoa(recipes[index])
	}
	return res
}

func TestOne(t *testing.T) {
	cc := []struct {
		in  int
		out string
	}{
		{9, "5158916779"},
		{5, "0124515891"},
		{18, "9251071085"},
		{2018, "5941429882"},
	}

	allPassed := true
	for _, c := range cc {
		allPassed = t.Run(c.out, func(t *testing.T) {
			got := takeTenAt(generateRecipes(c.in), c.in)
			if got != c.out {
				t.Fatalf("expected %s, got: %s", c.out, got)
			}
		}) && allPassed
	}
	if !allPassed {
		t.Fatal()
	}
	t.Fatal(takeTenAt(generateRecipes(580741), 580741))
}

func findPattern(spattern string) int {
	pattern := make([]int, len(spattern))
	for i, s := range spattern {
		pattern[i] = int(s - '0')
	}

	fmt.Println(pattern)

	foundMatching := 0

	recipes := []int{3, 7}
	elves := [...]int{0, 1}
	for foundMatching < len(pattern) {
		score := recipes[elves[0]] + recipes[elves[1]]
		if tens := score / 10; tens > 0 {
			if tens == pattern[foundMatching] {
				foundMatching += 1
			} else if tens == pattern[0] {
				foundMatching = 1
			} else {
				foundMatching = 0
			}
			recipes = append(recipes, tens)
		}
		recipes = append(recipes, score%10)

		if foundMatching >= len(pattern) || (score%10) == pattern[foundMatching] {
			foundMatching += 1
		} else if (score % 10) == pattern[0] {
			foundMatching = 1
		} else {
			foundMatching = 0
		}
		// fmt.Println(recipes, foundMatching)

		elves[0] = (elves[0] + recipes[elves[0]] + 1) % len(recipes)
		elves[1] = (elves[1] + recipes[elves[1]] + 1) % len(recipes)
	}

	return len(recipes) - foundMatching
}

func TestTwo(t *testing.T) {
	cc := []struct {
		in  string
		out int
	}{
		{"51589", 9},
		{"01245", 5},
		{"92510", 18},
		{"59414", 2018},
	}

	allPassed := true
	for _, c := range cc {
		allPassed = t.Run(c.in, func(t *testing.T) {
			got := findPattern(c.in)
			if got != c.out {
				t.Fatalf("expected %d, got: %d", c.out, got)
			}
		}) && allPassed
		if !allPassed {
			t.Fatal()
		}
	}
	if !allPassed {
		t.Fatal()
	}
	t.Fatal(findPattern("580741"))
}
