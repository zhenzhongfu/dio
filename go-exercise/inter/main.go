package main

import (
	"inter/b"
	//use of internal package inter/d/internal/c not allowed
	//"inter/d/internal/c"
)

func main() {
	b.Say()
	//c.Hello()
}
