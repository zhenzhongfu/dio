package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/garyburd/redigo/redis"
)

func main() {
	RedisConn := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	/*
		conn := RedisConn.Get()
		defer conn.Close()

		// set
		_, err := conn.Do("SET", "key1", "value1")
		if err != nil {
			fmt.Println("set fail ", err)
		}

		// get
		reply, err := conn.Do("GET", "key1")
		if err != nil {
			fmt.Println("get fail ", err)
		}
		rt, err := redis.Bytes(conn.Do("GET", "key1"))
		fmt.Println("reply type ", reflect.TypeOf(reply), string(rt))

	*/
	ctx, _ := context.WithCancel(context.Background())
	group, _ := errgroup.WithContext(ctx)

	//sub
	onMessage := func(channel string, msg []byte) error {
		fmt.Printf("channel: %s, message: %s\n", channel, msg)
		return nil
	}
	done := make(chan error, 1)
	group.Go(func() error {
		conn := RedisConn.Get()
		defer conn.Close()
		psc := redis.PubSubConn{Conn: conn}
		if err := psc.Subscribe("chan1"); err != nil {
			return err
		}

		go func() {
			for {
				switch n := psc.Receive().(type) {
				case error:
					done <- n
					return
				case redis.Message:
					if err := onMessage(n.Channel, n.Data); err != nil {
						done <- err
						return
					}
				case redis.Subscription:
					switch n.Count {
					case 1:
						fmt.Println("all done")
					case 0:
						fmt.Println("return")
						done <- nil
						return
					}
				}
			}
		}()

		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
	loop:
		for {
			select {
			case <-ticker.C:
				fmt.Println("ping...")
				if err := psc.Ping(""); err != nil {
					break loop
				}
			case <-ctx.Done():
				break loop
			case err := <-done:
				return err
			}
		}

		psc.Unsubscribe()
		return nil
	})

	time.Sleep(time.Second * 3)

	//pub
	group.Go(func() error {
		conn := RedisConn.Get()
		defer conn.Close()
		for i := 0; i < 3; i++ {
			reply, err := conn.Do("PUBLISH", "chan1", "woooooo")
			if err != nil {
				fmt.Println("publish err:", err)
			}
			fmt.Println("publish ", reply)
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		fmt.Println("wait --", err)
	}
}
