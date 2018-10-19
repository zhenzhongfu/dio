package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

func findStep(len int) int {
	bit := 8
	for {
		value := 1 << (uint)(bit)
		if value > len {
			return bit + 1 // 1<<(bit+1) = 2^bit
		}
		bit += 1
	}
}

func main() {
	b := make([]byte, 12)
	binary.BigEndian.PutUint32(b, 12)
	binary.BigEndian.PutUint32(b, 33)
	copy(b[8:], "aabb")
	fmt.Println(b)

	fmt.Printf("%s\n", time.Now())

	fmt.Println(math.Pow(2, 10))
	fmt.Println("11--", (int)(math.Ceil(math.Log2(1025))))

	fmt.Println(findStep(256))
	fmt.Println(findStep(1024))
	fmt.Println(findStep(1333))
	fmt.Println(findStep(65535))
}
