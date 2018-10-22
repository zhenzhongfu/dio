//model_02
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	group, newCtx := errgroup.WithContext(ctx)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			group.Go(func() error {
				for {
					select {
					case <-newCtx.Done():
						fmt.Println("handler done.")
						return nil
					default:
						// recv and send from conn.
						time.Sleep(time.Second * 1)
						fmt.Println(conn)
					}
				}
				return nil
			})
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

	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("All done.")
}
