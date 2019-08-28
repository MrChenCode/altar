package controller

import (
	"altar/application/core/basic"
	"github.com/gin-gonic/gin"
)

type Book struct{}

func (_ *Book) BookInfo(ctx *basic.Controller) {
	ctx.Log.Info("controller_bookid", 8000, "controller_bookname", "盘龙")
	ctx.Log.Error("http_get_error", "无效的bookname")
	ctx.Log.Info("controller_getmodel", 1)
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  ctx.Model.Book.GetBookInfo(),
	})
}
