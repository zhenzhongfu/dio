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
	/*
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("err: ", err)
	}
		fmt.Printf("result: ", status)
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n\r\n")

	var read,write []byte
	num, err := conn.Read(read)
	if err !=nil {
		fmt.Println(err)
	}
	fmt.Println(read, num)
	fmt.Println("000")
	read, err := ioutil.ReadAll(conn)
	if err !=nil {
		fmt.Println(err)
	}
	fmt.Println("111")
	fmt.Println(string(read))
	*/

	/*
	read := make([]byte, 1024)
	num, err := conn.Read(read)
	if err !=nil {
		fmt.Println(err)
	}
	fmt.Println(string(read), num)

	var write []byte
	write = []byte("HTTP/1.1 200 OK\r\n\r\n")
	conn.Write(write)
	*/
	fmt.Println("connect success")
	return
}

func main() {
	ln, err := net.Listen("tcp", ":19000")
	if err != nil {
		fmt.Println("listen err:", err)
	}
	i := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
		}
		i++
		fmt.Println("gogogo -- ", i)
		go handleConnection(conn)
		time.Sleep(10 * time.Second)
	}
	return
}
