package main

import (
	"fmt"
	"net"
	"time"
	_ "bytes"
	_ "io/ioutil"
	_ "bufio"
)

func handleConnection(conn net.Conn) {
	fmt.Println("conn handle")
}

func main() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen err:", err)
	}
	count := 1
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
		}
		count++
		fmt.Println("accept count: ", count)
		go handleConnection(conn)
		time.Sleep(10 * time.Second)
	}
	return
}
