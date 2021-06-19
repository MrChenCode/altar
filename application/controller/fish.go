package controller

import (
	"altar/application/context/cctx"
	"strconv"
)

type Fish struct{}

func (_ *Fish) GetFishList(ctx *cctx.ControllerContext) {
	pageId, _ := strconv.Atoi(ctx.Query("page_id"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	UserId ,_ := strconv.Atoi(ctx.Query("user_id"))
	if pageId < 0 {
		pageId = 0
	}
	if pageSize < 0 {
		pageSize = 5
	}
	info , code ,err := ctx.Model.Fish.GetFishList(pageId, pageSize,UserId)
	if err != nil {
		ctx.ResponseERR(code, err.Error())
		return
	}
	ctx.ResponseOK(info)
}

func ( _ *Fish) AddFishInfo (ctx *cctx.ControllerContext) {
	userId := ctx.PostForm("user_id")
	title := ctx.PostForm("title")
	weight := ctx.PostForm("weight")
	lenght := ctx.PostForm("lenght")
	address := ctx.PostForm("address")
	err := ctx.Model.Fish.UploadImgAndFishInfo(ctx.Context, title, weight, lenght, address, userId)
	if err != nil {
		ctx.ResponseERR(1000, err.Error())
		return
	}
	ctx.ResponseOK(nil)
}
