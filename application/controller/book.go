package controller

import (
	"altar/application/context/cctx"
	"github.com/gin-gonic/gin"
)

type Book struct{}

func (_ *Book) BookInfo(ctx *cctx.ControllerContext) {
	ctx.Log.Info("controller_bookid", 8000, "controller_bookname", "盘龙")
	ctx.Log.Error("http_get_error", "无效的bookname")
	ctx.Log.Info("controller_getmodel", 1)
	res, _ := ctx.Redis.Get("altar_redis_test").Result()
	ctx.JSON(200, gin.H{
		"code":   0,
		"msg":    ctx.Model.Book.GetBookInfo(),
		"result": res,
	})
}
