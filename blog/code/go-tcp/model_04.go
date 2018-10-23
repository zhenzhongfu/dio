//model_04
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"golang.org/x/sync/errgroup"
)

func main() {
	//pprof
	go func() {
		fmt.Println(http.ListenAndServe(":8887", nil))
	}()

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			handler(ctx, conn)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	select {
	case <-quit:
		{
			fmt.Println("recv quit signal")
			cancel()
		}
	}

	fmt.Println("All done.")
}

var queue = make(chan []byte, 100)

func handler(topCtx context.Context, conn net.Conn) error {
	ctx, cancel := context.WithCancel(context.Background())
	group, newCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return readRoutine(topCtx, newCtx, cancel, conn)
	})

	group.Go(func() error {
		return sendRoutine(topCtx, newCtx, cancel, conn)
	})
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
	return nil
}

func readRoutine(topCtx, ctx context.Context, cancel context.CancelFunc, conn net.Conn) error {
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
			n, err := conn.(*net.TCPConn).Read(buf)
			if err != nil {
				fmt.Println(err)
				return err
			}

			buf = buf[:n]
			fmt.Printf("read(%d) %s\n", n, string(buf))
			if err := process(buf); err != nil {
				return err
			}
		}
	}
	return nil
}

func sendRoutine(topCtx, ctx context.Context, cancel context.CancelFunc, conn net.Conn) error {
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
		case buf := <-queue:
			n, err := conn.(*net.TCPConn).Write(buf)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			fmt.Printf("send(%d) %s\n", n, string(buf))
		}
	}
}

func process(buf []byte) error {
	fmt.Printf("process: %s\n", string(buf))
	queue <- buf
	return nil
}
