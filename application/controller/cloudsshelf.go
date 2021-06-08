package controller

import (
	"altar/application/context/cctx"
	"github.com/gin-gonic/gin"
)

type Cloudsshelf struct {}

func (_*Cloudsshelf) Getbooklist  (c *cctx.ControllerContext){
	//获取分页
	page := c.Query("page")
	pagesize := c.Query("pagesize")
	if page == ""  ||  pagesize == "" {
		c.ResponseERR(10000, "page or pagesize fail")
		return
	}
	bookList := c.Model.Cloudsshelf.Getbooklist(1,20)

	c.JSON(200, gin.H{
		"code":   0,
		"msg":    "ok",
		"result": bookList,
	})
}