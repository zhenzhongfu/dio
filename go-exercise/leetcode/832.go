package main

import "fmt"

func flipAndInvertImage(A [][]int) [][]int {
	for _, item := range A {
		length := len(item)
		t := 0
		for i := 0; i < len(item)/2; i++ {
			t = item[i]
			item[i] = item[length-i-1]
			item[length-i-1] = t
		}

		for i := 0; i < len(item); i++ {
			switch item[i] {
			case 0:
				item[i] = 1
			case 1:
				item[i] = 0
			}
		}

		fmt.Println(item)
	}
	fmt.Println(A)
	return A
}

func main() {
	a := [][]int{{1, 1, 0}, {1, 0, 1}, {0, 0, 0}}
	flipAndInvertImage(a)
}
