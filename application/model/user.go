package model

import (
	"errors"
	"strconv"
)

type UserModel struct {
	*BasicModel
}
type UserInfo struct {
	UserName string `json:"user_name""`
	Phone int `json:"phone"`
}

func (u *UserModel) Login(phone , name string) (*UserInfo, error) {
	user := &UserInfo{}
	u.ctx.Log.Info("model_gameid", 1001, "model_gamename", "三国战纪")
	sql := "select id, user_name, phone from  user where phone = ? and user_name = ?  limit 1"
	result, err := u.User.ctx.Mysql.QueryRow(sql, phone, name)
	if err != nil {
		return user, err
	}
	if len(result) == 0 {
		return user, errors.New("用户名字或者密码错误")
	}
	user.UserName = result["user_name"]
	user.Phone, _ = strconv.Atoi(result["phone"])
	return user, nil

}
