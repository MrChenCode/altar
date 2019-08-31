// +build windows

package fore

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	NoSupport = errors.New("windows system is not supported")
)

//初始化、启动、重启、停止服务接口
//本接口的所有函数，必须是非阻塞的
type Service interface {
	//初始化操作
	//主要用于启动时初始化服务使用，启动后会首次调用init函数
	Init() error

	//启动操作
	//初始化init之后，会执行start操作，用于启动服务
	Start() error

	//获取要启动的http server，如果不启动httpserver，返回nil
	HttpServers() []*http.Server

	//启动完成，监听信号之前调用
	RunSuccess() error

	//重启操作
	//当捕获重启信号之后，在操作重启操作前，会调用Restart函数
	Restart() error

	//触发重启操作时，传递给下方的环境变量(xxx=yyy)
	RestartEnvs() []string

	//停止操作
	//当捕获退出信号时，退出服务之前调用此函数，
	Stop()

	//当遇到任何错误时，都会调用此函数，传递具体错误信息
	Error(error)
}

//启动服务
//默认监听syscall.SIGINT, syscall.SIGTERM信号退出服务
//同时监听syscall.SIGUSR2信号进行服务重启
func run(service Service) error {
	var err error
	if err := service.Init(); err != nil {
		return err
	}

	if err := service.Start(); err != nil {
		return err
	}

	var app *App = nil
	httpServers := service.HttpServers()
	if len(httpServers) > 0 {
		app, err = HttpServe(httpServers)
		if err != nil {
			return err
		}
		go func(app *App) {
			errChan := app.ErrorChan()
			for {
				err, ok := <-errChan
				if !ok {
					return
				}
				service.Error(err)
			}
		}(app)
	}

	if err := service.RunSuccess(); err != nil {
		return err
	}

	//到此处，所有服务已经初始化完成，此时可以干掉父进程

	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			signal.Stop(ch)
			if app != nil {
				app.Stop()
			}
			service.Stop()
			return nil
		}
	}
}

//启动
func Run(service Service) error {
	return run(service)
}

//停止
func RunStop(pid int) error {
	return NoSupport
}

//重启
//pid 要重启的服务pid
func RunRestart(pid int) error {
	return NoSupport
}
