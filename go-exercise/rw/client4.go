package main

import (
	"net"
	"io"
	"fmt"
	"time"
)

func handle() {

}

func main() {
	conn, err := net.DialTimeout("tcp", ":8888", time.Second * 2)
	if err != nil {
		fmt.Println("cannot connect")
		return
	}

	go handle()
	defer conn.Close()
	time.Sleep(time.Second * 10)
}
