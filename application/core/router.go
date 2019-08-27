package core

import (
	"github.com/gin-gonic/gin"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/controller"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/basic"
	"time"
)

type HandlerFunc func(*basic.Controller)

func (core *Core) Router(engine *gin.Engine) {
	book := &controller.Book{}

	engine.GET("/bookinfo", core.handle(book.BookInfo))
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
		c.Model = core.modelPool.Get().(*basic.Model)
		c.Model.Reset(c.RequestContext)

		//执行具体方法
		handler(c)
		core.controllerPool.Put(c)

		latency := time.Now().Sub(start)

		infokv, errkv := c.GetLog()
		*infokv = append([]interface{}{
			"method", c.Request.Method,
			"httpcode", c.Writer.Status(),
			"path", c.Request.URL.Path,
			"rawquery", c.Request.URL.RawQuery,
			"http_errinfo", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"request_time", latency.Seconds(),
		}, *infokv...)

		c.WriteLogInfo(c.G.RequestID, *infokv...)
		if *errkv != nil {
			c.WriteLogError(c.G.RequestID, *errkv...)
		}
	}
}
