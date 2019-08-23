package core

import (
	"github.com/gin-gonic/gin"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/controller"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/basic"
)

type HandlerFunc func(*basic.Controller)

func (core *Core) Router(engine *gin.Engine) {
	book := &controller.Book{}

	engine.GET("/bookinfo", core.handle(book.BookInfo))
}

func (core *Core) handle(handler HandlerFunc) gin.HandlerFunc {
	return func(ginc *gin.Context) {
		c := core.pool.Get().(*basic.Controller)

		c.Reset(ginc)

		handler(c)
		core.pool.Put(c)
	}
}
