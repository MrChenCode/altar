package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/basic"
)

type Book struct{}

func (_ *Book) BookInfo(ctx *basic.Controller) {
	ctx.Loginfo("gid", 101)
	ctx.Loginfo("gname", "三国战纪")

	ctx.JSON(200, gin.H{
		"code": 0,
		"msg": "",
	})
}
