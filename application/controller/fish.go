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
