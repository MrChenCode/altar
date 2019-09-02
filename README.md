# Altar 祭坛 [![pipeline status](https://gitlab.baidu-shucheng.com/panda/altar/badges/master/pipeline.svg)](https://gitlab.baidu-shucheng.com/panda/altar/commits/master)

altar祭坛是一个制造生产英雄的起点...

## 安装使用

go version >= 1.12  

先下载安装GO(如已经安装，则跳过)：  

MacOS: https://dl.google.com/go/go1.12.9.darwin-amd64.pkg  
Linux: https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz  
Windows: https://dl.google.com/go/go1.12.9.windows-amd64.msi  

源码：https://dl.google.com/go/go1.12.9.src.tar.gz  

更新版本请移步至：  

golang官网：https://golang.org/dl/  
谷歌国内镜像: https://golang.google.cn/dl/  


## 设置环境

安装完毕后，需要设置环境变量GOROOT、GOPATH、GO111MODULE、GOPROXY

GOROOT: 安装的go的根目录(安装到哪里就设置哪里，比如GOROOT=/usr/local/go)  
GOPATH: 放一些依赖代码缓存的目录(自行随便设置，比如GOPATH=/home/work/gopath)  

本项目基于golang module模式开发，所以需要固定设置环境变量GO111MODULE=on(永久设置)  

因为众所周知的原因，一些墙外的依赖包无法直接获取，golang提供了一个环境变量GOPROXY  
GOPROXY: 下载依赖或者更新依赖库时，以GOPROXY为跳板下载，这里提供一个地址https://goproxy.io  

当需要下载第三方库时，可使用:  

```shell
GOPROXY=https://goproxy.io go get golang.org/x/net
```


也可以在更新module依赖时直接使用:  

```shell
GOPROXY=https://goproxy.io go mod tidy
```

如果觉得每次都在命令前面加GOPROXY，也可以把GOPROXY设置为固定的环境变量  


## 下载项目

```shell
> git clone https://gitlab.baidu-shucheng.com/panda/altar.git
> 
> cd altar
> 
> GOPROXY=https://goproxy.io go mod tidy
```

## 配置文件

首先拷贝一个altar_default.ini到altar.ini   
针对当前的环境，设置altar.ini里面的running为qa或者online模式, 会针对不同的模式，加载不同的配置文件

## 开发运行项目

####运行参数说明

./altar [-chtv] args
* -c ini_path  手动设置配置文件路径
* -h 显示帮助
* -t 检测配置文件
* -v 显示编译等信息
* stop 关闭停止
* restart 重启

示例：  
./appname -c /etc/altar.ini  
./appname -h  
./appname -t  
./appname -v  
./appname stop  
./appname restart  

####开发测试(qa)

配置文件running设置为qa
```shell
//命令行模式
> go run main.go
```
也可以使用ide的run运行   
命令行和ide run为阻塞模式，使用ctrl-c即可停止

开发模式下也可以下载监听自动重启的开发工具mds：https://gitlab.baidu-shucheng.com/gt/mds  
具体操作方式请参照mds


####正式环境(online)

首先在Makefile中设置正式环境的项目名(APP_NAME = appname)  
修改ini配置文件的running为online(必须)  
然后编译运行：
```shell
> make
> ./appname
```
正式环境使用子进程非阻塞启动，如需要重启或者关闭，可以使用以下方式：
```shell
//普通优雅关闭, 优雅关闭需要等待所有处理连接正常处理完成，可能存在延迟关闭
//如果需要立即关闭，请使用kill -9 pid

> kill pid
//或者
> ./appname stop

//重启操作, 重启必须等待所有连接处理完成
> kill -USR2 pid
//或者
> ./appname restart
```
针对http keep-alive的情况，停止或者重启时会对http请求添加header头 Connection: close   
对于短时间内没有请求的keep-alive长连接，会最多等待60秒，超时后直接关闭

## END...
结束...


