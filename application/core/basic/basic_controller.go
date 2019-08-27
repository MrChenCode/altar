package basic

import (
	"github.com/gin-gonic/gin"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/context"
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
