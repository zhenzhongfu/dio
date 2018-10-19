package main

import (
	"fmt"
	"syscall"
)

func main() {
	var err error = syscall.Errno(2)
	fmt.Println(err.Error()) // "no such file or directory"
	fmt.Println(err)         // "no such file or directory
}
