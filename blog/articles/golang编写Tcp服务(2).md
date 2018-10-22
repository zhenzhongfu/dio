---
title: golang编写Tcp服务(2)
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
                        return nil                         
                    default:                               
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


