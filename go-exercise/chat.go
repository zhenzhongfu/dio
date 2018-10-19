package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8003")
	if err != nil {
		log.Fatal(err)
	}
	// 广播
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// 连接处理
		go handleConn(conn)
	}
}

//----------------------
type client chan<- string

var (
	msgChan = make(chan string)
	inChan  = make(chan client)
	outChan = make(chan client)
)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-msgChan:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-inChan:
			clients[cli] = true
		case cli := <-outChan:
			delete(clients, cli)
			// 关闭chan
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	// TODO in
	ch := make(chan string)
	// 从ch读数据
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	msgChan <- who + " has arrived"
	// 这里是将chan放入chan
	inChan <- ch

	// TODO msg
	input := bufio.NewScanner(conn)
	for input.Scan() {
		// 读单行
		msgChan <- who + " : " + input.Text()
	}

	// TODO out
	outChan <- ch
	msgChan <- who + " has left"
	if err := conn.Close(); err != nil {
		fmt.Println("close 1:", err)
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}

	if err := conn.Close(); err != nil {
		fmt.Println("close 2:", err)
	}
}
