package main

import (
	"cs/slice1"
	"cs/slice2"
	"fmt"
)

func showPt(s []int) {
	fmt.Printf("data 2: %p \n", &s)
}

func init() {
	fmt.Println("slice")
	slice2.Void()
	slice1.Void()
}

func main() {
	data := []int{1, 2, 3, 34, 5, 5, 6, 67, 77, 87}
	/*
		fmt.Println(data)
		s1 := data[:0:0]
		s1 = append(s1, 123)
		s2 := data[:0]
		s2 = append(s2, 10)
		fmt.Printf("data %p %d %d\n", &data[0], len(data), cap(data))
		fmt.Printf("s2 %p %d %d\n", &s2[0], len(s2), cap(s2))
		fmt.Printf("s1 %p %d %d\n", &s1[0], len(s1), cap(s1))


		data2 := make([]int, 0, 10)
		data2 = append(data2, []int{1,2,3,4,5}...)
		s1 = data2[0:0]
		s2 = data2[:0]
		fmt.Println("s1: ", s1, " | ", len(s1), cap(s1))
		fmt.Println("s2: ", s2, " | ", len(s2), cap(s2))

		fmt.Printf("data 1: %p \n", &data)
		showPt(data)

		fmt.Println(data[1:])
	*/

	data2 := make([]int, 20, 20)
	copy(data2, data)
	data[1] = 1000
	fmt.Println(data)
	fmt.Println(data2)
}
