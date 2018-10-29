package main

import (
	"fmt"
	"net"
	"time"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe(":8887", nil))
	}()

	ln, err := net.Listen("tcp", ":8886")
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 3)

	go func() {
		for {
			time.Sleep(time.Second * 1)
			_, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	var list []net.Conn
	for i := 0; i < 10000; i++ {
		conn, err := net.Dial("tcp", ":8886")
		if err != nil {
			fmt.Printf("conn: %s\n", err)
		} else {
			list = append(list, conn)
		}
		fmt.Println("connect:", conn, i)
	}

	time.Sleep(time.Hour)
}
