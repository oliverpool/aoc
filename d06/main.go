package main

import (
	"strconv"
	"strings"
)

func seatID(s string) int64 {
	r := strings.NewReplacer("F", "0", "B", "1", "L", "0", "R", "1")
	b := r.Replace(s)
	i, _ := strconv.ParseInt(b, 2, 64)
	return i
}
