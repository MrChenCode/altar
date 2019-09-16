package controller

import (
	"altar/application/context/cctx"
	"strconv"
)

type Book struct{}

func (_ *Book) BookInfo(ctx *cctx.ControllerContext) {
	bookid, _ := strconv.Atoi(ctx.Query("bookid"))
	//channelid, _ := strconv.Atoi(ctx.Query("channelid"))

	// 是否检测伪下架的书 1检测（检测的意思就是伪下架为真实下架） 默认0 不检测
	//checkPseudoOffline, _ := strconv.Atoi(ctx.Query("heck_pseudo"))

	if bookid <= 0 {
		ctx.ResponseERR(10000, "no data!")
		return
	}

}
