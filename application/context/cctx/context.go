//这是controller使用的上下文资源
//他包含了请求上下文、gin框架基础上下文、以及model对象

package cctx

import (
	"altar/application/context/rctx"
	"altar/application/library"
	"altar/application/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ControllerContext struct {
	//公共请求上下文资源
	*rctx.RequestContext

	//gin请求资源
	*gin.Context

	//model
	Model *model.BasicModel

	//lib
	Library *library.BasicLibrary
}

func (c *ControllerContext) Reset(ginctx *gin.Context) {
	//重置请求上下文
	c.RequestContext.Reset(ginctx)

	//重置gin请求资源
	c.Context = ginctx
}

func (c *ControllerContext) ResponseOK(v interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":   1006,
		"msg":    "ok",
		"result": v,
	})
}

func (c *ControllerContext) ResponseERR(code int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"msg":    msg,
		"result": nil,
	})
}
