package basic

import (
	"altar/application/core/context"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	//公共请求上下文资源
	*context.RequestContext

	//gin请求资源
	*gin.Context

	//model
	Model *Model
}

func (c *Controller) Reset(ginctx *gin.Context) {
	c.Context = ginctx
}
