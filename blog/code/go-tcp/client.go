package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conn)
	time.Sleep(time.Second * 10)
}
