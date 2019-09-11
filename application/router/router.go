package router

import (
	"altar/application/context"
	"altar/application/context/cctx"
	"altar/application/context/rctx"
	"altar/application/controller"
	"altar/application/library"
	"altar/application/model"
	"github.com/gin-gonic/gin"
	"sync"
	"sync/atomic"
	"time"
)

type HandlerFunc func(ctx *cctx.ControllerContext)

type Router struct {
	ctx *context.Context

	//controller pool对象
	pool sync.Pool

	//restart是否处于restart状态
	restartd int32
}

func (r *Router) Router(engine *gin.Engine) {
	book := &controller.Book{}
	game := &controller.Game{}
	topic := &controller.Topic{}

	engine.GET("/bookinfo", r.handle(book.BookInfo))
	engine.GET("/gameinfo", r.handle(game.GameInfo))
	engine.GET("/topicinfo", r.handle(topic.GetTopicInfo))
}

func (r *Router) handle(handler HandlerFunc) gin.HandlerFunc {
	return func(ginctx *gin.Context) {
		start := time.Now()
		c := r.pool.Get().(*cctx.ControllerContext)
		//初始化basicController
		c.Reset(ginctx)

		if atomic.LoadInt32(&r.restartd) == 1 {
			c.Header("Connection", "close")
		}
		//执行具体方法
		handler(c)

		infokv, errkv := c.GetLog()
		requestID := c.G.RequestID

		r.pool.Put(c)

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

func NewRouter(ctx *context.Context) *Router {
	r := &Router{
		ctx: ctx,
	}
	r.pool.New = func() interface{} {
		rcx := rctx.NewRequestContext(r.ctx)
		lib := library.NewLibrary(rcx)
		cx := &cctx.ControllerContext{
			RequestContext: rcx,
			Model:          model.NewModel(rcx, lib),
			Library:        lib,
		}
		return cx
	}

	return r
}

func (r *Router) Restart() {
	atomic.StoreInt32(&r.restartd, 1)
}
