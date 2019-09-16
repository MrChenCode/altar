package controller

import (
	"altar/application/context/cctx"
	"github.com/gin-gonic/gin"
)

type Topic struct{}

func (_ *Topic) GetTopicInfo(c *cctx.ControllerContext) {
	tid := "5a0bdc104cecc186538b4567"
	res, err := c.Model.Topic.GetTopicDetail(tid, 1)

	if err != nil {
		c.ResponseERR(10000, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code":   0,
		"msg":    "ok",
		"result": res,
	})
}
