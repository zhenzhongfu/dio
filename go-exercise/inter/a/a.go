package a

import (
	"fmt"
	_ "unsafe"
)

//go:linkname say inter/b.Say
func say() {
	fmt.Println("hello.")
}
