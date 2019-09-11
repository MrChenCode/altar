package main

import (
	"altar/application/config"
	"altar/application/context"
	"altar/application/logger"
	"altar/application/router"
	"altar/fore"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	//显示帮助
	help bool

	//配置文件路径
	//如果没有设置参数，则扫描默认的配置文件路径，会根据项目编译目录为根目录, 扫描顺序如下：
	//		./altar.ini(默认)
	//		./ini/altar.ini
	//		/etc/altar.ini(linux)
	inifile string

	//检测配置文件
	initest bool

	//打印info
	version bool

	//以下参数由编译时指定
	//编译时间
	BuildTime string

	//编译作者
	BuildAuthor string

	//编译根目录
	BuildPath string

	//编译类型，列表如下：
	//	make(由makefile编译, 默认)
	//	run(直接运行go run)
	//	mds(通过mds编译)
	BuildType string
)

func init() {
	var mds bool
	flag.BoolVar(&help, "h", false, "帮助")
	flag.StringVar(&inifile, "c", "./altar.ini", "设置`配置文件`路径")
	flag.BoolVar(&version, "v", false, "显示详细编译信息")
	flag.BoolVar(&initest, "t", false, "检测配置文件")
	flag.BoolVar(&mds, "mds", false, "以mds运行")

	flag.Usage = usageHelp

	if mds {
		BuildType = "mds"
	}

	if BuildPath != "" {
		config.RootPath = BuildPath
	}
}

type altar struct {
	conf   *config.Config
	engine *gin.Engine
	router *router.Router
	log    *logger.Logger
	logwf  *logger.Logger
}

func main() {
	//解析命令行参数
	flag.Parse()

	//处理部分参数输出指令
	outFlag()

	a := &altar{}

	//如果是windows，把编译模式改为run
	if runtime.GOOS == "windows" {
		BuildType = "run"
	}

	switch BuildType {
	case "make":
		if err := fore.Run(a); err != nil {
			log.Fatal(err.Error())
		}
	default:
		if err := a.Init(); err != nil {
			log.Fatal(err.Error())
		}
		if err := a.Start(); err != nil {
			log.Fatal(err.Error())
		}
		httpServers := a.HttpServers()
		if len(httpServers) > 0 {
			if err := httpServers[0].ListenAndServe(); err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}

//初始化操作
func (a *altar) Init() error {
	conf, err := getConfig(false)
	if err != nil {
		return err
	}
	a.conf = conf

	//处理重启、停止指令
	args := flag.Args()
	if len(args) > 0 {
		pid, err := a.GetPid()
		if err != nil {
			return err
		}
		for _, cmd := range args {
			cmd = strings.ToLower(strings.TrimSpace(cmd))
			switch cmd {
			case "restart":
				if err := fore.RunRestart(pid); err != nil {
					return err
				}
				os.Exit(0)
			case "stop":
				if err := fore.RunStop(pid); err != nil {
					return err
				}
				os.Exit(0)
			}
		}
		return fmt.Errorf("Altar Error: 无效的命令 %s\n", strings.Join(args, " "))
	}
	return nil
}

//start启动
func (a *altar) Start() error {
	loginfo, err := logger.NewConfig(
		filepath.Join(a.conf.LogPath, a.conf.LogFileName),
		a.conf.LogCatTime,
		a.conf.LogRetainDay,
	)
	if err != nil {
		return err
	}
	logwf, err := logger.NewConfig(
		filepath.Join(a.conf.LogPath, a.conf.LogFileName+".wf"),
		a.conf.LogCatTime,
		a.conf.LogRetainDay,
	)
	if err != nil {
		return err
	}

	a.log = loginfo
	a.logwf = logwf

	ctx, err := context.NewContext(a.conf, loginfo, logwf)
	if err != nil {
		return err
	}

	r := router.NewRouter(ctx)
	a.router = r

	if a.conf.Debug == config.DebugOnline {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	gin.DisableConsoleColor()

	e := gin.New()
	e.Use(router.Recovery(logwf))

	r.Router(e)
	a.engine = e

	return nil
}

//返回需要启动的http server
func (a *altar) HttpServers() []*http.Server {
	s := &http.Server{
		Addr:         a.conf.HttpServerAddr,
		Handler:      a.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return []*http.Server{s}
}

//启动完成的通知
func (a *altar) RunSuccess() error {
	a.log.Infow("",
		"msg", "start",
		"pid", os.Getgid(),
		"http_listen", a.conf.HttpServerAddr,
		"debug", a.conf.Debug,
	)
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(a.conf.PidFile, []byte(pid), 0644); err != nil {
		return err
	}
	return nil
}

//重启时，需要传递给子进程的环境变量
func (a *altar) RestartEnvs() []string {
	return nil
}

//重启的通知
func (a *altar) Restart() error {
	a.router.Restart()
	if a.log != nil {
		a.log.Infow("", "msg", "restart", "pid", os.Getgid(), "ppid", os.Getppid())
		_ = a.log.Sync()
	}

	if a.logwf != nil {
		_ = a.logwf.Sync()
	}
	return nil
}

//服务停止的通知
func (a *altar) Stop() {
	if a.log != nil {
		a.log.Infow("", "msg", "stop", "pid", os.Getgid(), "ppid", os.Getppid())
		_ = a.log.Sync()
	}
	if a.logwf != nil {
		_ = a.logwf.Sync()
	}
}

func (a *altar) Error(err error) {
	if a.logwf != nil {
		a.logwf.Errorw("", "msg", "service", "errinfo", err.Error(), "pid", os.Getgid())
	}
}

//根据配置文件pid路径，获取pid
func (a *altar) GetPid() (int, error) {
	oldpid, err := ioutil.ReadFile(a.conf.PidFile)
	if err != nil {
		return 0, err
	}
	if len(oldpid) == 0 {
		return 0, errors.New("no pid found")
	}
	pid, err := strconv.Atoi(string(oldpid))
	if err != nil {
		return 0, err
	}
	if pid <= 0 {
		return 0, errors.New("pid is invalid")
	}
	return pid, nil
}

//获取配置文件对象
func getConfig(out bool) (*config.Config, error) {
	//检测配置文件
	file := iniPath()

	if file == "" {
		return nil, fmt.Errorf("altar: no valid configuration file was found, use -c to set the configuration file path")
	}

	err := config.SyntaxCheck(file)
	if err != nil {
		return nil, fmt.Errorf("altar error: the configuration file %s syntax is error, %s\n", file, err.Error())
	} else if out {
		_, _ = fmt.Fprintf(os.Stdout, "altar: the configuration file %s syntax is ok\n", file)
	}
	conf, err := config.NewConfig(file)
	if err != nil {
		return nil, fmt.Errorf("altar error: configuration file %s test is error, %s\n", file, err.Error())
	} else if out {
		_, _ = fmt.Fprintf(os.Stdout, "altar: configuration file %s test is successful\n", file)
	}
	return conf, nil
}

//返回ini的绝对路径
func iniPath() string {
	//先检查是否是绝对路径, 已经是绝对路径的，直接使用
	if filepath.IsAbs(inifile) {
		return inifile
	}

	//首先获取项目根路径
	var path string
	if BuildPath != "" {
		path = BuildPath
	} else {
		path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		BuildPath = path
	}

	//如果是相对路径，先判断根据当前启动文件路径计算
	inis := []string{
		"./altar.ini",
		"./ini/altar.ini",
		filepath.Join(path, inifile),
		filepath.Join(path, "altar.ini"),
		filepath.Join(path, "ini", "altar.ini"),
	}
	//如果是非windows系统，则尝试扫描/etc/altar.ini
	if runtime.GOOS != "windows" {
		inis = append(inis, "/etc/altar.ini")
	}
	for _, inif := range inis {
		fi, err := os.Stat(inif)
		//如果存在，且不是目录，则使用这个
		if err == nil {
			if !fi.IsDir() {
				return inif
			}
		}
	}
	//如果都找不到，则返回空
	return ""
}

func outFlag() {
	switch {
	case help:
		flag.Usage()
	case version:
		_, _ = fmt.Fprintf(os.Stdout, "Altar %s (Operator %s) %s/%s (build %s)\n",
			BuildTime, BuildAuthor, runtime.GOOS, runtime.GOARCH, runtime.Version())
	case initest:
		_, err := getConfig(true)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, err.Error())
		}
	default:
		return
	}
	os.Exit(0)
}

//显示帮助
func usageHelp() {
	_, _ = fmt.Fprintf(os.Stdout, `Altar(祭坛)
使用命令: %s [-chtv] args

参数说明:
`, os.Args[0])

	flag.PrintDefaults()
}
