package controller

import (
	"altar/application/core/basic"
	"github.com/gin-gonic/gin"
)

type Book struct{}

func (_ *Book) BookInfo(ctx *basic.Controller) {
	ctx.Log.Info("get_bookid", 75699, "get_bookname", "吞噬星空")
	ctx.Log.Error("http_get_error", "无效的bookname")
	ctx.Log.Info("method", "PPP")
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  ctx.Model.Book.GetBookInfo(),
	})
}
