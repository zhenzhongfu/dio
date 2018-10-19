package main

import (
	"fmt"
	"strings"
)

func judgeCircle(moves string) bool {
	u := strings.Count(moves, "U")
	d := strings.Count(moves, "D")
	l := strings.Count(moves, "L")
	r := strings.Count(moves, "R")
	if u == d && l == r {
		return true
	}
	return false
}

func main() {
	r := judgeCircle("RLUURDDDLU")
	fmt.Println(r)
}
