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

	buf := []byte("i am spiderman.")
	n, err := conn.(*net.TCPConn).Write(buf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("send(%d) %s\n", n, buf)

	time.Sleep(time.Second * 3)

	n, err = conn.(*net.TCPConn).Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("recv(%d) %s\n", n, buf)

	time.Sleep(time.Second * 10)
}
