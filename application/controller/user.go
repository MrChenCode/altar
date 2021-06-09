package controller

import (
	"altar/application/context/cctx"
	"regexp"
	"strconv"
)

type User struct{}

func (_ *User) Login(ctx *cctx.ControllerContext) {
	userName := ctx.Query("user_name")
	password := ctx.Query("password")
	//status := VerifyMobileFormat(phone)
	status, _ := strconv.Atoi( ctx.Query("status"))
	user, err := ctx.Model.User.Login(password, userName, status)
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
