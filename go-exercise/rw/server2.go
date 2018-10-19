package main

import (
	"fmt"
	"io"
	"net"
)

func handle(conn net.Conn) {
	bytes := []byte("hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world hello world|")
	num, err := conn.Write(bytes)
	if err != nil {
		fmt.Println("write err:", err)
		return
	}
	fmt.Println("write num:", num)

	buf := bytes[0:0]
	tmp := make([]byte, 1)
	for {
		num, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			} else {
				break
			}
		}

		if string(tmp[0]) == "|" {
			break
		}
		buf = append(buf, tmp[:num]...)
		fmt.Println("read num:", num, " msg:", string(tmp[:num]), " || ", len(buf), cap(buf))
	}
	fmt.Println(len(buf), cap(buf), string(buf[:len(buf)]))

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

	//	time.Sleep(time.Second * 10000)
}
