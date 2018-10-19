package main

import (
	"bytes"
	"fmt"
)

var a = []string{".-", "-...", "-.-.", "-..", ".", "..-.", "--.", "....", "..", ".---", "-.-", ".-..", "--", "-.", "---", ".--.", "--.-", ".-.", "...", "-", "..-", "...-", ".--", "-..-", "-.--", "--.."}

func uniqueMorseRepresentations(words []string) int {
	same := 0
	m := make(map[string]bool)
	for _, w := range words {
		var bb bytes.Buffer
		for _, s := range w {
			bb.WriteString(a[s-97])
		}

		s := bb.String()
		if _, ok := m[s]; ok {
		} else {
			m[s] = true
			same += 1
		}
	}
	return same
}

func main() {
	s := []string{"gin", "zen", "gig", "msg"}
	fmt.Println(uniqueMorseRepresentations(s))
}
