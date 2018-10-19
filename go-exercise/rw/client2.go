package main

import (
	"net"
	"time"
	"fmt"
	"io"
)

func handle(conn net.Conn) {
	//read
	buf := make([]byte, 0, 10)
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
		// 粘包
		if string(tmp[0]) == "|" {
			break
		}
		buf = append(buf, tmp[:num]...)
		fmt.Println("read num:", num, " msg:", string(tmp[:num]), " || ", len(buf), cap(buf))
	}
	fmt.Println(len(buf), cap(buf), string(buf[:len(buf)]))

	//write
	buf = buf[0:0]
	buf = append(buf[0:0], []byte("echo hello|")...)
	num, err := conn.Write(buf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("write: ", num, string(buf), len(buf), cap(buf))
}

func main() {
	conn, err := net.DialTimeout("tcp", ":8888", time.Second * 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	go handle(conn)
	defer conn.Close()
	time.Sleep(time.Second * 10)
}
