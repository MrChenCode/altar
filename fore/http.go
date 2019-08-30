package fore

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
)

type App struct {
	//网络包
	net *Net
	//需要启动的http server
	servers []*http.Server
	//监听的server tcp
	listeners []net.Listener
	//httpdown
	http *HTTPDOWN
	sds  []HttpDownServer
	//遇到的错误
	errors chan error
}

func newApp(servers []*http.Server) *App {
	return &App{
		net:       &Net{},
		servers:   servers,
		listeners: make([]net.Listener, 0, len(servers)),
		http:      &HTTPDOWN{},
		sds:       make([]HttpDownServer, 0, len(servers)),
		errors:    make(chan error, len(servers)*2),
	}
}

//监听端口
func (a *App) listen() error {
	for _, s := range a.servers {
		l, err := a.net.Listen("tcp", s.Addr)
		if err != nil {
			return err
		}
		if s.TLSConfig != nil {
			l = tls.NewListener(l, s.TLSConfig)
		}
		a.listeners = append(a.listeners, l)
	}
	return nil
}

//http serve
func (a *App) serve() {
	for i, s := range a.servers {
		a.sds = append(a.sds, a.http.Serve(s, a.listeners[i]))
	}
}

//启动app服务
func (a *App) run() error {
	if err := a.listen(); err != nil {
		return err
	}
	a.serve()
	return nil
}

//stop会同时关闭所有http.Server
func (a *App) Stop() {
	if len(a.servers) == 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(a.sds) * 2)
	for _, s := range a.sds {
		go func(s HttpDownServer) {
			defer wg.Done()
			if err := s.Wait(); err != nil {
				a.errors <- err
			}
		}(s)
		go func(s HttpDownServer) {
			defer wg.Done()
			if err := s.Stop(); err != nil {
				a.errors <- err
			}
		}(s)
	}
	wg.Wait()
	close(a.errors)
}

//重启http服务，返回重启成功后的pid和可能遇到的错误
func (a *App) Restart(envs []string) (pid int, err error) {
	return a.net.StartProcess(envs)
}

func (a *App) ErrorChan() <-chan error {
	return a.errors
}

func HttpServe(servers []*http.Server) (*App, error) {
	app := newApp(servers)
	if err := app.run(); err != nil {
		return nil, err
	}
	return app, nil
}
