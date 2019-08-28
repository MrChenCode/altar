package core

import (
	"altar/application/controller"
	"altar/application/core/basic"
	"github.com/gin-gonic/gin"
	"time"
)

type HandlerFunc func(*basic.Controller)

func (core *Core) Router(engine *gin.Engine) {
	book := &controller.Book{}
	game := &controller.Game{}

	engine.GET("/bookinfo", core.handle(book.BookInfo))
	engine.GET("/gameinfo", core.handle(game.GameInfo))
}

func (core *Core) handle(handler HandlerFunc) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		start := time.Now()
		c := core.controllerPool.Get().(*basic.Controller)
		//初始化公共请求资源
		c.RequestContext.Reset(ginctx)

		//初始化basicController
		c.Reset(ginctx)

		//初始化model
		c.Model.Reset(c.RequestContext)

		//执行具体方法
		handler(c)

		infokv, errkv := c.GetLog()
		requestID := c.G.RequestID

		core.controllerPool.Put(c)

		latency := time.Now().Sub(start)

		infokv = append([]interface{}{
			"method", ginctx.Request.Method,
			"httpcode", ginctx.Writer.Status(),
			"path", ginctx.Request.URL.Path,
			"rawquery", ginctx.Request.URL.RawQuery,
			"http_errinfo", ginctx.Errors.ByType(gin.ErrorTypePrivate).String(),
			"request_time", latency.Seconds(),
		}, infokv...)

		c.WriteLogInfo(requestID, infokv...)
		if errkv != nil {
			c.WriteLogError(requestID, errkv...)
		}
	}
}
