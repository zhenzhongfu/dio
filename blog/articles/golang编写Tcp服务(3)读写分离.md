---
title: golang编写Tcp服务(3)读写分离
tags: go,program
grammar_cjkRuby: true
---

本节讲解一下读写分离。

仅靠Handler协程是无法完成在socket上同时处理读写的。[net](https://golang.org/pkg/net/)包的网络IO操作，比如Dial，Read，Write均为block，虽然可以使用类似SetReadDeadline这类接口设置block的超时时间，但是在没有IO请求的时候，调用接口的协程处于挂起状态，不能提供服务。
来看一个例子：
```golang?linenums
// model_03.go
...
func handler(ctx context.Context, conn net.Conn) error {
    for {                                               
        buf := make([]byte, 100)                        
        select {                                        
        case <-ctx.Done():                              
            fmt.Println("handler done.")                
            return nil                                  
        default:                                        
            // recv and send from conn.                 
            n, err := conn.(*net.TCPConn).Read(buf)     
            if err != nil {                             
                fmt.Println(err)                        
                return err                              
            }                                           
                                                        
            fmt.Printf("recv(%d) %s\n", n, buf)         
            buf = buf[:n]                               
            n, err = conn.(*net.TCPConn).Write(buf)     
            if err != nil {                             
                fmt.Println(err)                        
            }                                           
            fmt.Printf("send(%d) %s\n", n, buf)         
        }                                               
    }                                                   
    return nil                                          
}
...

// client.go
...
func main() {                               
    conn, err := net.Dial("tcp", ":8000")   
    if err != nil {                         
        fmt.Println(err)                    
    }                                       
                                            
    buf := []byte("i am spiderman.")        
    n, err := conn.(*net.TCPConn).Write(buf)
    if err != nil {                         
        fmt.Println(err)                    
    }                                       
    fmt.Printf("send(%d) %s\n", n, buf)  
	time.Sleep(time.Second * 3)
                                            
    n, err = conn.(*net.TCPConn).Read(buf)  
    if err != nil {                         
        fmt.Println(err)                    
    }                                       
    fmt.Printf("recv(%d) %s\n", n, buf)     
                                            
    time.Sleep(time.Second * 10)            
}                
...
```
分别启动server和client。server的输出：
```shell?linenums
$ go run model_03.go 
recv(15) i am spiderman.
send(15) i am spiderman.
^Crecv quit signal
EOF
EOF
All done.
```
在recv和send的内容输出之后，马上ctrl+c终止server时会发现程序不会退出，而是在输出“recv quit signal”后继续被阻塞，说明并没有走到ctx.Done的case，而是依然阻塞在Read调用处，等待连接上的后续请求，直到client超时退出，Read请求返回io.EOF，才会继续后面的流程。

使用pprof验证我们的猜想。[pprof](https://golang.org/pkg/net/http/pprof/)包通过其HTTP服务器运行时分析数据提供pprof可视化工具所期望的格式。
```golang?linenums
...
import (
	...
	_ "net/http/pprof"
)
...
func main() {
	//pprof                                           
	go func() {                                       
		fmt.Println(http.ListenAndServe(":8887", nil))
	}() 
...
```
这里我们启动了一个web server，访问http://localhost:8887/debug/pprof/ ，在浏览器可以查看到goroutine、heap、profile等程序在运行时收集到的分析数据。再次运行上个例子，ctrl+c之后查看goroutine信息，能够看到其中的一个goroutine被阻塞到model_03.go的74行，也就是Read调用的地方，上面的每一行都是在调用过程中的堆栈信息。
![model_03](./images/1540265979478.png)
打开源码/usr/local/go/src/runtime/netpoll.go的173行，可以看到是在IO等待中，这步操作是阻塞的，Go将基于事件的IO复用模型封装在Runtime里，底层还是基于epoll的。
![netpoll.go](./images/1540266284935.png)

所以，一个Handler协程还需要多个协程配合共同完成在同一socket连接上的收发处理。服务器普遍处理的模式都是一问一答，一条上行消息，解析处理完毕后回复一条下行消息，至少需要两个协程处理阻塞，一读一写。

修改handler协程，
```golang?linenums
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
```
在handler协程创建读写协程时，将其放入同一个group的好处是，若其中之一退出，会调用cancel通知另一协程同步退出，handler协程的group.Wait便会返回，不需要在handler里再监听退出信号。
在这里，将handler协程、读协程、写协程看成一组，负责处理同一个连接上的消息，同时创建，同时退出。
```golang?linenums
// 读协程
func readRoutine(topCtx, ctx context.Context, cancel context.CancelFunc, conn net.Conn) error {
    defer func() {                                                                                 
		// 调用cancel通知group退出
        cancel()                                                                               
    }()                                                                                        
                                                                                               
    for {                                                                                      
        select {   
		// 监听程序退出信号
        case <-topCtx.Done():                                                                  
            fmt.Println("readRoutine top done.")                                               
            return nil  
		// 监听group退出信号
        case <-ctx.Done():                                                                     
            fmt.Println("readRoutine done.")                                                   
            return nil                                                                         
        default:                                                                               
       		// TODO read from conn
			// TODO unpack message
			// TODO process message                                                                            
        }                                                                                      
    }                                                                                          
    return nil              
}                                                                                              
```
读协程逻辑简单，收消息，解析，处理，再收下一条消息。
```golang?linenums
func sendRoutine(topCtx, ctx context.Context, cancel context.CancelFunc, conn net.Conn) error { 
    defer func() {                                                                              
		// 调用cancel通知group退出                                                             
        cancel()                                                                                
    }()                                                                                         
                                                                                                
    for {                                                                                       
        select {                                                                                
		// 监听程序退出信号
        case <-topCtx.Done():                                                                  
            fmt.Println("readRoutine top done.")                                               
            return nil  
		// 监听group退出信号
        case <-ctx.Done():                                                                     
            fmt.Println("readRoutine done.")                                                   
            return nil                                                                        
        case message, ok := <-queue: 
			// TODO send message to socket
        }                                                                                       
    }                                                                                           
}                                                                                               
```
写协程与读协程类似，但数据来自channel queue，若queue被关闭，ok为false。
queue应该与conn封装在同一个struct里，这样便于处理。


Go包含了对OOP的支持。用代码来说明。
```golang?linenums
type netConn struct {
	conn net.Conn
	queue chan []byte
}
func (c *netConn) test() {
}
```
上面的代码定义了一个类型netConn，而test是一个与netConn类型关联的函数也就是一个方法。一个面向对象的程序借助方法来表达及操作对象属性，而不是直接去操作对象。
test前面附加的参数叫做方法的receiver，receiver可以是指针或非指针类型，取决于使用方式。
- 指针类型的receiver，需要注意指针指向的还是原地址，对于receiver的修改也是在原地址上的修改。
- 非指针类型的receiver，进行了一次传值拷贝，需考虑对象大小造成的影响。

在前面的例子中，都是用函数去直接操作一个个单个的对象，这里将他们封装起来，使用方法处理。
单独声明一个package。
```golang?linenums
package network
...
type netConn struct {
	conn net.Conn
	queue chan []byte
}

func NewNetConn(conn net.Conn, queueLen int) *netConn {
    return &netConn{                                   
        conn:  conn,                                   
        queue: make(chan []byte, queueLen),            
    }                                                  
}                                                      

func (c *netConn) ReadRoutine(topCtx, ctx context.Context, cancel context.CancelFunc) error {
...
	n, err := c.conn.(*net.TCPConn).Read(buf) 
...
}

func (c *netConn) SendRoutine(topCtx, ctx context.Context, cancel context.CancelFunc) error {
...
	case buf := <-c.queue:
		n, err := c.conn.(*net.TCPConn).Write(buf)   
...
}
```
Go使用首字母大小写来控制包内的类型，变量，函数，方法对外部包是否可见，大写可见，小写为不可见。使用大写首字母的标识符均会从定义它们的包中导出。



