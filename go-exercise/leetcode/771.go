package main

import (
	"fmt"
	"strings"
)

func numJewelsInStones(J string, S string) int {
	m := make(map[string]int)
	var b string
	for i := 0; i < len(J); i++ {
		b = J[i : i+1]
		_, ok := m[b]
		if !ok {
			m[b] = 1
		}
	}

	sum := 0
	for i := 0; i < len(S); i++ {
		b = S[i : i+1]
		_, ok := m[b]
		if ok {
			sum += 1
		}
	}
	return sum
}

func num2(J string, S string) int {
	m := make(map[rune]bool)
	for _, b := range J {
		m[b] = true
	}
	var sum int
	for _, b := range S {
		_, ok := m[b]
		if ok {
			sum += 1
		}
	}
	return sum
}

func num3(J string, S string) int {
	var sum int
	for _, b := range S {
		sum += strings.Count(J, string(b))
	}
	return sum
}

func main() {
	//sum := numJewelsInStones("abc", "cccddjddjdjdj")
	sum := num3("abc", "cccddjddjdjdj")
	fmt.Println(sum)
}
