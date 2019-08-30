package fore

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
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
	//启动时，是否是子进程启动，如果不是，则直接启动一个新子进程
	if os.Getenv(altarReactorChildProcess) == "" {
		pid, err := startProcess(nil, nil)
		if err != nil {
			return err
		}
		if runtime.GOOS == "windows" {
			return nil
		}
		rch := make(chan os.Signal, 10)
		signal.Notify(rch, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-time.After(10 * time.Second):
			//非windows系统此函数不会返回err
			p, err := os.FindProcess(pid)
			if err == nil {
				//释放子进程的资源
				_, _ = p.Wait()
			}
		case <-rch:
			return nil
		}
		return nil
	}

	//defer func() {
	//	if os.Getppid() != 1 {
	//		syscall.Kill(os.Getppid(), syscall.SIGTERM)
	//	}
	//}()

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
	ppid := os.Getppid()
	if ppid != 1 {
		if err := syscall.Kill(ppid, syscall.SIGINT); err != nil {
			return fmt.Errorf("failed to close parent: %s", err)
		}
	}

	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	/*--------------------------------------------------*/
	//处理子进程重启失败后，存在僵尸进程的问题
	restartZP := make(chan int, 10)
	closeZP := make(chan struct{}, 10)
	go func(rzp chan int, czp chan struct{}) {
		for {
			select {
			case pid := <-rzp:
				//如果收到了重启的信号，30s后，如果还没有收到停止信号，则尝试回收子进程
				select {
				case <-czp:
					return
				case <-time.After(30 * time.Second):
					//非windows系统此函数不会返回err
					p, err := os.FindProcess(pid)
					if err == nil {
						//释放子进程的资源
						_, _ = p.Wait()
					}
				}
			}
		}
	}(restartZP, closeZP)
	///*--------------------------------------------------*/

	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			signal.Stop(ch)
			if app != nil {
				app.Stop()
			}
			closeZP <- struct{}{}
			service.Stop()
			return nil
		case syscall.SIGUSR2:
			if service.Restart() == nil {
				envs := service.RestartEnvs()
				var pid int
				if app != nil {
					pid, err = app.Restart(envs)
				} else {
					pid, err = startProcess(nil, envs)
				}
				if err != nil {
					service.Error(err)
				}
				restartZP <- pid
				//此处不能退出，需要监听子进程的kill信号
			}
		}
	}
}

//启动
func Run(service Service) error {
	return run(service)
}

//停止
func RunStop(pid int) error {
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("fore to stop parent: %s", err)
	}
	return nil
}

//重启
//pid 要重启的服务pid
//envs 重启时，附带的环境变量，xxx=yyy
func RunRestart(pid int) error {
	if err := syscall.Kill(pid, syscall.SIGUSR2); err != nil {
		return fmt.Errorf("fore to restart parent: %s", err)
	}
	return nil
}
