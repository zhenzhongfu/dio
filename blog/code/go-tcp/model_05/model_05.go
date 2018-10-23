//model_05
package main

import (
	"context"
	"fmt"
	"go-tcp/model_05/network"
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

	c := network.NewNetConn(conn, 100)
	group.Go(func() error {
		return c.ReadRoutine(topCtx, newCtx, cancel)
	})

	group.Go(func() error {
		return c.SendRoutine(topCtx, newCtx, cancel)
	})
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
	return nil
}
