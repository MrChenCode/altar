package controller

import (
	"altar/application/context/cctx"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type User struct {}


func (_ *User) Login(ctx *cctx.ControllerContext) {
	var result []interface{}
	userName := ctx.Query("user_name")
	phone, _ := strconv.Atoi(ctx.Query("phone"))
	code, _:= strconv.Atoi(ctx.Query("code"))
	if (code != 111) {
		ctx.JSON(200, gin.H{
			"code":   10005,
			"msg":    "验证码不正确",
			"result": result,
		})
	}

	fmt.Print(phone)
	fmt.Print(userName)
}
