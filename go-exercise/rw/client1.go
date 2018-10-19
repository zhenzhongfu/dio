
package main

import (
	"fmt"
	"net"
	"time"
)

func handle(conn net.Conn) {
	fmt.Println("conn handle")
}

func main() {
	var list []net.Conn
	for i:= 1; i<100000; i++ {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:8888", time.Second * 2)
		if err != nil {
			fmt.Println("connect error:", err)
			continue
		}

		list = append(list, conn)
		fmt.Println("conn count: ", i)
		go handle(conn)
	}

	time.Sleep(time.Second * 1000000)
	return
}
