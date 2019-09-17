package controller

import (
	"altar/application/context/cctx"
	"github.com/gin-gonic/gin"
)

type Topic struct{}

func (_ *Topic) GetTopicInfo(c *cctx.ControllerContext) {
	tid := c.Query("tid")
	if tid == "" {
		c.ResponseERR(10000, "Parameter TID is invalid!")
		return
	}
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
