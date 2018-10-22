---
title: golang编写Tcp服务(1)
tags: go,program
grammar_cjkRuby: true
---

Go是一门表达力强，简洁，干净，有效的编程语言，它特有的并发机制使得其能够更为容易的编写充分利用多核机器性能的程序。用Go写一个高并发的server与基于I/O多路复用模型的server，在思路和处理细节上有诸多不同。本篇文章通过编写一个完整可用的例子，来描述如何使用Go编写tcp server。

## 概念
### 多路复用模型
在Linux上基于I/O多路复用模型的server，大体是这样的
```c
int main() {
	listen()
	epoll_create()
	epoll_ctrl()
	while(1) {
		eventNum = epoll_wait()	
		while(eventNum > 0) {
			handleEvent()
			eventNum--
		}
	}
}
```
本质上是事件循环，系统在对应端口上捕获到的事件，均会被放入事件列表，程序编写者需要循环处理这份列表，根据具体的事件进行相应的回调处理。其中epoll_wait可传入时间参数以设定其在等待事件的空窗时间。整个模型循环在一个os进程中，若要利用多核性能，就必须加入多进程和多线程的处理，通常这部分会比较复杂。一些通用的做法是，判断若是读事件，使用main thread去recv数据，raw数据丢到mq，多个worker thread从mq上读取raw数据进行解析并处理。但凡涉及到多进程多线程，就不可避免地需要考虑共享资源的使用，也就不可避免地使用到锁。锁是特别会增加程序员心智负担的一个东西，稍不注意就造成系统锁死，bad guy。所以写过C再写Erlang和Go的时候就会很开心。另外，thread的监控也会比较头疼，若worker thread异常退出，该如何处理？若要通知其他worker thread退出，也很麻烦。

### Go模型
使用Go时，server模型大体是这样的：
```golang?linenums
func main() {
	ln, _ := net.Listen()
	for {
		conn, _ := ln.Accept()
		go handleConn(conn)
	}
}
```
这里最关键的一行line5，关键字go就已经完成了对多核特性的支持。Go将I/O多路复用模型封装在runtime里了，底层的事件再不需要开发者注册和回调。go程并不是一个OS线程，它更为轻量级，这里引用[The Go scheduler](http://morsmachine.dk/go-scheduler)的几张图说明:

![M-P-G](http://morsmachine.dk/in-motion.jpg)
>- M表示OS线程。
>- G表示goroutine，包含独有的stack，指令指针及一些调度信息。
>- P表示一个处于调度的上下文，视作一个局部调度器，用来将goroutine绑定到一个具体的线程上。这是Go从N:1调度器到M:N调度器的关键。

![syscall](http://morsmachine.dk/syscall.jpg)
>Go1.5以上的版本，P被默认为CPU core的数量。
>灰色的G被维护在一个队列里等待被调度，Go在syscall被调用之前，会将P和原本的M0解绑，重新寻找空闲的M1绑定，以便在M0阻塞时还能同时处理队列中的其他G。

![steal](http://morsmachine.dk/steal.jpg)
>当某个P的G队列跑完了，而其他P队列还有G，会尝试进行steal操作，获取其他P的G，保证所有M能够全负荷运行。

goroutine相比thread更为轻量，一个Go程序中可以并发成千上万个goroutine的系统调度和资源占用开销会更小。但有了goroutine，并不代表就不需要处理共享数据和资源，在Go哲学里，强调的是：
>不要通过共享内存来通信，而要通过通信来共享内存。

### 同步
channel作为Go的同步机制，通过传递数据结构的引用来完成goroutine之间的通信，传递的是数据的所有权，无需上锁。
channel在使用上类似mq，channel可以指定容量，当某个channel上被未读数据占满时，向其写入的goroutine会被阻塞。相反，channel为空时，读取的goroutine也会被阻塞。以下代码，容量为0的channel与阻塞操作无异。
```golang?linenums
mq chan int
mq = make(int)	// mq
// mq = make(int, 100)
// goroutine1
go func() {
	for {
		select {
			case msg:= <- mq:
			time.Sleep(time.Minute)
		}
	}
}()
// goroutine2
go func() {
	for {
		mq <- 1
	}
}
```

goroutine在启动后，除非自己退出，否则不能被停止的，唯一的方法就是通过channel，当然实现起来也很容易。
```golang?linenums
func main() {
	done := make(chan int, 1)
	go func() {
		// do sth.
		// 通知主进程退出
		done <- 0
	}()
	
	select {
	case <- done:
		// quit
	}
}
```
以往的thread通信机制，常用的那几种，不管是消息队列，还是共享内存，使用和维护起来还是比较复杂的，尤其是对于锁的争用。
Go提供了sync包，提供基本同步操作，结合goroutine是比较容易写出一个并发程序的，上面的代码引入sync包之后：
```golang?linenums
func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// do sth.
	}

	wg.Wait()
}
```
