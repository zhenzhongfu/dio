package main

import (
	"fmt"
)

func hammingDistance(x int, y int) int {
	tx := x
	ty := y
	bbx := make([]int, 0)
	bby := make([]int, 0)
	for {
		if tx == 0 && ty == 0 {
			break
		}
		if tx > 0 || ty > 0 {
			bbx = append(bbx, tx%2)
			bby = append(bby, ty%2)
			tx = tx / 2
			ty = ty / 2
		} else if tx > 0 {
			bbx = append(bbx, tx%2)
			tx = tx / 2
		} else if ty > 0 {
			bby = append(bby, ty%2)
			ty = ty / 2
		}
	}
	fmt.Println(bbx, bby)

	i := 0
	sum := 0
	for {
		if i >= len(bbx) && i >= len(bby) {
			break
		}

		if bbx[i] != bby[i] {
			sum += 1
		}

		i++
	}

	return sum
}

func fn(x int, y int) int {
	count := 0
	for x != 0 || y != 0 {
		if x%2 != y%2 {
			count += 1
		}
		x = x >> 1
		y = y >> 1
	}
	return count
}

func main() {
	//	fmt.Println(hammingDistance(1, 3))
	fmt.Println(fn(1, 3))
}
