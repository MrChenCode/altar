package controller

import (
	"altar/application/context/cctx"
	"github.com/gin-gonic/gin"
)

type Index struct{}

func (_ *Index) Index(ctx *cctx.ControllerContext) {
	ctx.JSON(200,gin.H{"msg": "welcome to altar"})
}
