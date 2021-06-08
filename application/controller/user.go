package controller

import (
	"altar/application/context/cctx"
	"regexp"
	"strconv"
)

type User struct{}

func (_ *User) Login(ctx *cctx.ControllerContext) {
	userName := ctx.Query("user_name")
	phone := ctx.Query("phone")
	code, _ := strconv.Atoi(ctx.Query("code"))
	status := VerifyMobileFormat(phone)
	if status == false{
		ctx.ResponseERR(1005, "手机号错误")
		return
	}
	if code != 111 {
		ctx.ResponseERR(1005, "验证码错误")
		return
	}
	user, err := ctx.Model.User.Login(phone, userName)
	if err != nil {
		ctx.ResponseERR(1000, err.Error())
		return
	}
	ctx.ResponseOK(user)
}


func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147)|(12[0-9]))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
