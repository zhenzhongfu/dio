
package main

import (
	"fmt"
	"net"
	_ "bufio"
	"time"
)

func main() {
	var list []net.Conn
	for i:= 1; i<100000; i++ {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:19000", time.Second*2)
		if err != nil {
			fmt.Println("connect error:", err)
		}
		/*
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("err: ", err)
		}
		fmt.Printf("result: ", status, i)
		defer conn.Close()
		*/
		fmt.Println("ok, --i: ", i)
		list = append(list, conn)
	}

	time.Sleep(time.Second * 1000000)
	return
}
