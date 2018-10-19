package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

//var g errgroup.Group

func main() {
	ctx, _ := context.WithCancel(context.Background())
	group, errCtx := errgroup.WithContext(ctx)

	i := 0
	for ; i < 10; i++ {
		j := i
		group.Go(func() error {
			index := fmt.Sprintf("%d", j)
			fmt.Println("entry ", index)
			if j == 2 {
				fmt.Println("-----222 ", index)
				//这里一般都是某个协程发生异常之后，调用cancel()
				//这样别的协程就可以通过errCtx获取到err信息，以便决定是否需要取消后续操作
				return errors.New("md error " + index)
				//cancel()
			} else if j == 7 || j == 8 || j == 9 {
				time.Sleep(time.Second * 3)
				//检查 其他协程已经发生错误，如果已经发生异常，则不再执行下面的代码
				err := CheckGoErr(errCtx)
				if err != nil {
					fmt.Println("wooooooo -- err:", err, j)
					return err
				}
				fmt.Println("wooooooo")
			}
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		fmt.Println("wait -- ", err)
	}
}

func CheckGoErr(errContext context.Context) error {
	select {
	case <-errContext.Done():
		return errContext.Err()
	default:
		return nil
	}
}
