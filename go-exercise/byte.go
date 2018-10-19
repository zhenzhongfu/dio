package main

import "fmt"

func main() {
	b := []byte("12345678")
	fmt.Println(b)
	b[0] = b[3]
	copy(b, b[3:])
	b = b[:(8 - 3)]
	fmt.Println(b)
}
