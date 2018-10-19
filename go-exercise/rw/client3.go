package main

import (
	"net"
	"time"
	"fmt"
)

func handle(conn net.Conn) {
	//write
	total := 0
	buf := make([]byte, 65535)
	for {
		num, err := conn.Write(buf)
		if err != nil {
			fmt.Println(err)
		}
		total += num
		fmt.Println("write: ", num, total)
	}
}

func main() {
	conn, err := net.DialTimeout("tcp", ":8888", time.Second * 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			    //error handle
				fmt.Println("tcp conn wrong")
			}
		tcpConn.SetNoDelay(true)

		go handle(conn)
		defer conn.Close()
		time.Sleep(time.Second * 10)
	}
}
