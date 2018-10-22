---
title: golang编写Tcp服务(2)
tags: go,program
grammar_cjkRuby: true
---

## 
要完成一个tcp server服务，主要依赖[net](https://golang.org/pkg/net)包。
我们只需要三类goroutine即能完成一个简单的模型。
```golang?linenums
// main goroutine
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
- 主协程：完成Listen和wait动作。
- Accept协程：完成Accept动作。
- Handler协程：处理连接conn上的读写事件。

