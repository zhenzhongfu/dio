package main

import(
	"net"
	"fmt"
	"time"
	"io"
)

func handle(conn net.Conn) {
	buf := make([]byte, 60000)
	time.Sleep(time.Second * 10)
	total := 0
	for {
		time.Sleep(time.Second * 3)
		num, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			} else {
				fmt.Println(err)
				break
			}
		}
		total += num
		fmt.Println("total: ", total)
	}

	defer conn.Close()
}

func main() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			return
		}

		go handle(conn)
	}

	time.Sleep(time.Second * 10000)
}
