package network

import (
	"context"
	"fmt"
	"net"
)

type netConn struct {
	conn  net.Conn
	queue chan []byte
}

func NewNetConn(conn net.Conn, queueLen int) *netConn {
	return &netConn{
		conn:  conn,
		queue: make(chan []byte, queueLen),
	}
}

func (c *netConn) ReadRoutine(topCtx, ctx context.Context, cancel context.CancelFunc) error {
	defer func() {
		fmt.Println("read cancel")
		cancel()
	}()

	for {
		select {
		case <-topCtx.Done():
			fmt.Println("readRoutine top done.")
			return nil
		case <-ctx.Done():
			fmt.Println("readRoutine done.")
			return nil
		default:
			buf := make([]byte, 100)
			n, err := c.conn.(*net.TCPConn).Read(buf)
			if err != nil {
				fmt.Println(err)
				return err
			}

			buf = buf[:n]
			fmt.Printf("read(%d) %s\n", n, string(buf))
			if err := c.process(buf); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *netConn) SendRoutine(topCtx, ctx context.Context, cancel context.CancelFunc) error {
	defer func() {
		fmt.Println("send cancel")
		cancel()
	}()

	for {
		select {
		case <-topCtx.Done():
			fmt.Println("sendRoutine top done.")
			return nil
		case <-ctx.Done():
			fmt.Println("sendRoutine done.")
			return nil
		case buf := <-c.queue:
			n, err := c.conn.(*net.TCPConn).Write(buf)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Printf("send(%d) %s\n", n, string(buf))
		}
	}
}

func (c *netConn) process(msg []byte) error {
	fmt.Printf("process: %s\n", string(msg))
	c.queue <- msg
	return nil
}
