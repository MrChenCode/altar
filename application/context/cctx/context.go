//这是controller使用的上下文资源
//他包含了请求上下文、gin框架基础上下文、以及model对象

package cctx

import (
	"altar/application/context/rctx"
	"altar/application/library"
	"altar/application/model"
	"github.com/gin-gonic/gin"
)

type ControllerContext struct {
	//公共请求上下文资源
	*rctx.RequestContext

	//gin请求资源
	*gin.Context

	//model
	Model *model.Model

	//lib
	Library *library.Library
}

func (c *ControllerContext) Reset(ginctx *gin.Context) {
	//重置请求上下文
	c.RequestContext.Reset(ginctx)

	//重置gin请求资源
	c.Context = ginctx
}
