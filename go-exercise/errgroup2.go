package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, errCtx := errgroup.WithContext(ctx)

	for i := 0; i < 10; i++ {
		// goroutine闭包使用的变量全用copy出来的新值，
		//因为创建go需要时间，在go未被创建时for循环就跑完了
		//闭包拿到的值其实是最终的i
		tmp := i
		g.Go(func() error {
			if tmp == 2 {
				fmt.Println("index ", tmp)
				//这里一般都是某个协程发生异常之后，调用cancel()
				//这样别的协程就可以通过errCtx获取到err信息，以便决定是否需要取消后续操作
				cancel()
				// 这里在调用cancel之后还会执行，异常之后的收尾工作
				fmt.Println("err index ", tmp)
				return errors.New("errrrrrrrrr ")
			} else if tmp == 7 || tmp == 8 || tmp == 9 {
				time.Sleep(time.Second * 3)
				//检查 其他协程已经发生错误，如果已经发生异常，则不再执行下面的代码
				err := CheckGoErr(errCtx)
				if err != nil {
					fmt.Println("check err:", err, tmp)
					return err
				}
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("wait err :", err)
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
