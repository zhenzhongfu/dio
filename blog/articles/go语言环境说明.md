---
title: go语言环境说明
tags: go,program
grammar_cjkRuby: true
---

## 1.go语言环境说明

### 1.1 go的安装
go的安装步骤算是相对简洁的，[下载](https://golang.org/dl/)对应的版本，解压到指定目录，最后设置环境变量。下面以Linux版本的go1.11为例：
```shell
$ wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
$ tar -C /usr/local -xzf go1.11.linux-amd64.tar.gz
```
修改~/.bashrc，增加2行设置。
```bash
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin
```
然后测试一下是否成功安装。
```shell
$ source ~/.bashrc
$ go version
go version go1.11 linux/amd64
```

### 1.2 go安装包的简单说明
我们看一下刚刚解压安装的go的内部情况，其中包含了所有与go相关的目录和文件：
```go
$ cd /usr/local/go && ls
api  AUTHORS  bin  CONTRIBUTING.md  CONTRIBUTORS  doc  favicon.ico  lib  LICENSE  misc  PATENTS  pkg  README.md  robots.txt  src  test  VERSION
```
简单说明一下各个目录：

- api：这里放的是go的各个版本相对于上个版本的增量API特性，包括公开的变量，函数，参数情况，每个API一行。用于go的api检查。
- bin：go的可执行文件。包含go，godoc，gofmt。其中，[godoc](https://blog.golang.org/godoc-documenting-go-code)解析go源码文件生成html或者纯文本文档，[gofmt](https://blog.golang.org/go-fmt-your-code)是自动格式化go源码文件的工具。
>对于vim用户，这里有个提高效率的工具：[vim-go](https://github.com/fatih/vim-go)包含很多支持go的特性，尤其是goimports，能分析源码自动增加和删除头部的import说明。

>For Eclipse or Sublime Text users, the GoClipse and GoSublime projects add a gofmt facility to those editors.

git的pre-commit钩子也已就绪。
>And for Git aficionados, the misc/git/pre-commit script is a pre-commit hook that prevents incorrectly-formatted Go code from being committed. If you use Mercurial, the hgstyle plugin provides a gofmt pre-commit hook.

- doc: html格式的go的标准库说明文档。默认打开目录为空，我们可以执行godoc启动web server生成并展示这些说明文件。
- lib：库文件。go1.11里放的是时区信息文件。
- misc：一些混杂的库和工具。
- pkg：归档文件。其中类似linux_amd64的这类目录与平台相关，其名称由对应操作系统与处理器架构组合而成。仔细看，这些归档文件的文件名都与标准库的名称一一对应，归档文件的作用体现在，当某一个库文件变更时，只需要编译这个库对应的归档文件，然后再将原来其他的所有归档文件同这个新的归档文件链接在一起生成新的可执行程序，而不需要重新编译所有的库。
- src：go自身和相关工具以及标准库的源码。
- test：包含了go工具链和运行时的测试文件。

### 1.3 go的环境变量说明
在上面的1.1设置了一些环境变量，PATH和GOROOT，以及未提到的GOPATH，GOBIN，都是用来在编写go程序时会用到的环境变量，对于理解go的运作至关重要。
- PATH：环境变量PATH的值是一系列目录，在linux下用“:”分隔，当shell在执行命令或程序时，将会到PATH指定的目录中依次顺序查找，如go version这个命令，shell尝试在PATH目录查找是否go是否存在，你可以echo $PATH看看目前的PATH包含哪些目录。
- GOROOT：go自身的根目录。
- GOPATH：环境变量GOPATH的值也是一系列目录，在linux下用“:”分隔，其中每一个值代表着编写go程序的一个工作区目录。在build或者run的过程中，go tool会在GOPATH所包含的一系列目录中查找源码文件import声明的依赖，默认的查找顺序是GOROOT->GOPATH，然后根据GOPATH中声明的顺序依次查找。
- GOBIN：执行go install时，编译并存放最终可执行程序文件的目录。在GOBIN未被设置时，go install编译出来的文件放在各自工作区目录的bin文件夹中。

### 1.4 什么是工作区
工作区可以理解成项目特定的工程文件目录，我们写一个简单的demo来看一下工作区的概念吧。
```shell
$ mkdir myspace
```
上面的代码，我们创建了一个目录myspace，然后我们在~/.bashrc将myspace设置成工作区。在~/.bashrc增加一行：
```shell
export GOPATH="~/myspace"
```
一个典型工作区至少应该包含以下两个目录：
- src：包含以代码包形式组织的go源码文件。
- bin：包含通过go install编译生成的可执行程序。目录本身可自动生成。
另外还有一个目录：
- pkg：与go的pkg类似。如果执行了go build或者go install，可以看到对应的平台相关归档目录，除了与main包以及test相关的目录和文件不会被归档，其余的归档文件.a都与代码包一一对应。目录本身可自动生成。

什么是代码包？go语言是以代码包为基本组织单位，与目录一一对应，换而言之，一个目录下的所有文件都必须声明为同一个代码包。一般代码包被声明为目录名，方便记忆与查错。

我们来建立一个这样的目录层级关系，
```shell
~/myspace/
    src/
        toy/
            main.go
            /hello
                hello.go
        toy_util/
            main.go
```
在上面的层级关系中，myspace是一个工作区，toy是工作区下面的一个代码仓库(source repository)，main.go是toy仓库下面的命令源码文件，hello.go是toy仓库下面的包文件。命令源码是指用package main在头部声明并包含main函数的文件。

从上面目录层级关系看出，go源码组织形式是以环境变量GOPATH，工作区，工作区目录下的src，代码包等构成一个树状的目录层级结构。

一个典型的工作区包含许多代码仓库，上面的工作区就包含了toy和toy_util两个仓库。关于工作区和代码仓库，普遍的做法是将所有的go源码以及相关的依赖都放在同一个工作区，但是如果设置了多个工作区，并且不同工作区包含了不同功能的同名代码包，就要特别注意go的依赖查找顺序了。

### 1.5 第一个demo
我们修改一下toy和toy_util仓库的源码。
```go
// toy/main.go
package main
import (
    "toy/hello" // import hello包
    "toy_util"
)
func main() {
    result := toy_util.Add(1,2)
    hello.Print(result)
}

// toy/hello/hello.go
package hello
import (
    "fmt"
)
func Print(num int) {
    fmt.Println("hello world! ", num)
}

// toy_util/util.go
package toy_util
func Add(a, b int) int{
    return a + b
}
```
然后回到toy目录下编译和执行：
```shell
# 这步会在当前仓库目录compile出可执行文件toy
$ go build && ls
hello  main.go  toy
$ ./toy
hello world!
# 这步会compile出可执行文件到在当前工作区的bin目录
$ go install
```
顺序执行build和install后，目录树有一些变动：
```shell
~/myspace/
    src/
        toy/
            main.go
            /hello
                hello.go
        toy_util/
            util.go
    bin/
        toy
```
可以看到，build之后toy仓库下面的可执行文件toy不见了，bin和bin/toy被自动生成了。

### 1.6 常用的go tool命令
一般说的go tool是指$GOROOT/bin/go这个可执行程序，在shell中键入go help能看到go tool支持的所有命令。常用的有这些：
- build：compile packages and dependencies
- install：compile and install packages and dependencies
- run：compile and run Go program. 通过run执行Go程序，不会在当前仓库或者工作区bin目录生成可执行文件。
- get：download and install packages and dependencies. 通过go get命令，可以下载并且安装其他开发者和组织编写的包或依赖。

### 1.7 GOPATH带来的问题
若我们要引用其他开发者封装好的包，我们可以通过go get下载，go tool会根据参数提供的地址从github或者其他地址下载，但这里有个问题，get到的仓库版本无法指定，只能是对应仓库的HEAD。所以如果我们需要某个版本的文件，要么手动下载，要么使用其他工具辅助完成。

Go1.5版本开始，支持将当前仓库外部依赖拷贝到本地vendor目录，来实现简单的版本管理。假设toy_util是需要引入的外部包，我们将对应版本下载下来，放入当前项目仓库的vendor目录，那么，引入vendor之后我们项目的目录层级应该如下：
```shell
~/myspace/
    src/
        toy/
            main.go
            /hello
                hello.go
            vendor/     
                toy_util/   (这个是下载下来的1.0版本)
                    util.go
        toy_util/           (这个是刚刚我们修改过的版本)
            util.go
    bin/
        toy
        
```
那么，我们在执行编译时，go tool会优先从vendor目录查找依赖包，然后才是GOPATH目录。有了vendor目录，多人协作搭建开发环境的时候就不用每个人再去跑依赖包的go get了。不过若我们项目规模比较大，并且还引用了较多的外部依赖包，则下载、拷贝、记录版本这个过程就显得颇为琐碎了。

还好，我们还可以通过godep，govendor，glide这些管理工具来完成这件事。

官方觉得vendor还不够优雅，go1.11增加了modules的管理，代替了原来GOPATH的管理方式，而vendor的处理模式依然有效。

### 1.7 modules
####  go mod
Go 1.11包含了对Go modules的初步支持，以及一个新的module识别的go get命令。
官方[文档](https://golang.org/cmd/go/#hdr-Preliminary_module_support)描述，最快让目前项目用上Go  module的方法，将项目移出GOPATH/src到其他目录，在项目根目录创建go.mod文件，然后go build，go tool会自己查找相关的依赖。

Go 1.11增加了一个临时环境变量GO111MODULE，它可以被设置成on，off，auto(default).
- off：GOPATH mode，关闭module特性，继续使用GOPATH的方式管理包。
- on： module-aware mode，使用module功能，忽略GOPATH的设置。
- auto：若当前目录包含go.mod，则使用modle-aware mode，否则使用GOPATH mode。

先修改一下代码：
```go
// main.go
package main
import (
    "github.com/gin-gonic/gin"
)
func main() {
    gin.New()
}

// go.mod
module toy
```
项目层级关系是这样的：
```shell
~/test/
    toy/
        go.mod
        main.go
        /hello
            hello.go
    ...
```
这里可以看到test目录不在原来的GOPATH内，并且我们在代码中引用了别的包。
```shell
$ GO111MODULE=on go build
build toy: cannot find module for path github.com/gin-gonic/gin
```
没有成功。设置一下GOPROXY，就能正常下载第三方依赖包了，默认下载的是latest。当然，在module-ware下，我们还可以用go get来下载指定版本。
```shell
$ go get github.com/gin-gonic/gin@v1.1.4
```
TODO：关于[GOPROXY
](https://baokun.li/archives/goproxyio-intro/)要看一下源码了。
```shell
$ export GOPROXY=https://goproxy.io
$ go build   
```
于是现在的目录层级关系是这样：
```go
~/test/
    toy/
        go.mod
        go.sum
        main.go
        /hello
            hello.go
    ...
```
第三方依赖包则都会被下载到~/go/pkg/mod里。
打开go.mod，我们会看到我们项目所有的依赖包以及当前依赖的版本。
```go
 module toy
 
 require (
     github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
     github.com/gin-gonic/gin v1.3.0
     github.com/golang/protobuf v1.2.0 // indirect
     github.com/manucorporat/sse v0.0.0-20160126180136-ee05b128a739 // indirect
     github.com/mattn/go-isatty v0.0.4 // indirect
     github.com/ugorji/go/codec v0.0.0-20180831062425-e253f1f20942 // indirect
     golang.org/x/net v0.0.0-20180921000356-2f5d2388922f // indirect
     gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
     gopkg.in/yaml.v2 v2.2.1 // indirect
 )
```
#### 替换
若要修改依赖包的版本
```go
go get -u github.com/golang/protobuf@VERSION
```