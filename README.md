# Altar 祭坛

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
git clone https://gitlab.baidu-shucheng.com/panda/altar.git

cd altar

GOPROXY=https://goproxy.io go mod tidy
```

## 运行项目

```shell
go run main.go
```

也可以下载监听自动重启的开发工具mds：https://gitlab.baidu-shucheng.com/gt/mds  
具体操作方式请参照mds

## END...
结束...


