package main

import (
	"altar/application/core"
	"altar/application/core/config"
	"altar/application/core/context"
	"altar/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"path/filepath"
)

func main() {
	//var mds bool
	//flag.BoolVar(&mds, "mds", false, "mds running...")
	//flag.Parse()

	c, err := config.NewConfig("./bloc.ini")
	if err != nil {
		panic(err)
	}

	io.MultiWriter()
	log, err := logger.NewConfig(filepath.Join(c.LogPath, "bloc.log"), logger.CAT_DAY, 7)
	if err != nil {
		panic(err)
	}
	logwf, err := logger.NewConfig(filepath.Join(c.LogPath, "bloc.log.wf"), logger.CAT_DAY, 7)
	if err != nil {
		panic(err)
	}
	ctx := context.NewBasicController(c, log, logwf)

	if err := ctx.Init(); err != nil {
		panic(err)
	}

	cores := core.NewCore(ctx)

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(gin.Recovery())

	cores.Router(r)
	r.Run(":8888")
}