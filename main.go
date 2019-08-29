package main

import (
	"altar/application/config"
	"altar/application/context"
	"altar/application/router"
	"altar/logger"
	"github.com/gin-gonic/gin"
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
	ctx := context.NewController(c, log, logwf)

	if err := ctx.Init(); err != nil {
		panic(err)
	}

	r := router.NewRouter(ctx)

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	e := gin.New()
	e.Use(gin.Recovery())

	r.Router(e)
	_ = e.Run(":8888")
}
