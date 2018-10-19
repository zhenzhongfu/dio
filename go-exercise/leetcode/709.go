package main

import (
	"bytes"
	"fmt"
)

func toLowerCase(str string) string {
	var bb bytes.Buffer
	var i int
	for _, v := range str {
		i = int(v)
		if i > 64 && i < 91 {
			i = i + 32
			bb.WriteRune(rune(i))
		} else {
			bb.WriteRune(v)
		}
	}
	return bb.String()
}

func main() {
	s := toLowerCase("sdjfJJJdjdjdHHH")
	fmt.Println(s)
}
