---
title: golang编写Tcp服务(2)建立模型
tags: go,program
grammar_cjkRuby: true
---

## 
要完成一个tcp server服务，主要依赖[net](https://golang.org/pkg/net)包。
我们只需要三类goroutine即能完成一个简单的模型。
```golang?linenums
// model_01.go
package main   
import (    
    "fmt"      
    "net"      
    "os"    
    "os/signal"
    "syscall"
)       
func main() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {                          
		fmt.Println(err)                     
		return                               
	}   
	go func() {                     
    	for {                       
			conn, err := ln.Accept()
			if err != nil {         
				fmt.Println(err)    
				continue            
			}                       
			go func() {
				// recv and send from conn.
				fmt.Println(conn)
			}()
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
		}
	}   
}
```
- 主协程：完成Listen，等待退出信号。
- Accept协程：完成Accept动作。
- Handler协程：处理连接conn上的读写事件。

启动。
```golang?linenums
$ go run main.go &
$ netstat -anp|grep 8000
tcp      0      0 :::8000    :::*       LISTEN      37512/main 
```
kill进程，主协程退出，但不会通知和等待其他协程，程序终止。但在实际应用中，handler go经常会需要做一些收尾工作，比如回收资源以及通知其他服务，此时，我们需要借助sync包来完成。
```golang?linenums
// model_02.go
package main
import (
    "context"
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"

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
						time.Sleep(time.Second)
                        // recv and send from conn.        
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
```
相比model_01，
- 主协程多了三个动作，1）创建context并将handler加入到WaitGroup中；2）quit时执行cancel；3）wait所有的handler执行完毕。
- Handler协程多了一个动作，等待context的cancel消息。
测试一下connect的情况。
```golang?linenums
// client.go
package main
import (
    "fmt"
    "net"
    "time"
)
func main() {
    conn, err := net.Dial("tcp", ":8000")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(conn)
    time.Sleep(time.Second * 10)
}
```
分别run model_02.go和client.go，然后kill model_02。
```shell
$ go run model_02.go 
&{{0xc0000b2080}}
&{{0xc0000b2080}}
&{{0xc0000b2080}}
&{{0xc0000b2080}}
^Crecv quit signal
&{{0xc0000b2080}}
handler done.
All done.
```
“recv quit signal”，"handler done."，“All done.”依次输出。主协程在收到退出信号时，调用cancel()向context的quit channel发送消息，group有多少个成员，发送多少quit消息，quit消息的类型是struct{}。每个handler协程都从quit channel中获取1个quit消息，然后走退出流程。

将handler处理分离成函数。
```golang?linenums
// model_02.go
...
func main() {
...
	group.Go(func() error {   
		return handler(newCtx)
	})       
...

func hanlder(ctx context.Context) error {
    for {                                
        select {                         
        case <-ctx.Done():               
            fmt.Println("handler done.") 
            return nil                   
        default:                         
            // recv and send from conn.  
            time.Sleep(time.Second * 1)  
            fmt.Println(conn)            
        }                                
    }                                    
    return nil   
}     
```
