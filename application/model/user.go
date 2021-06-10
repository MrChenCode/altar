package model

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type UserModel struct {
	*BasicModel
}
type UserInfo struct {
	UserName string `json:"user_name"`
	UserId   int    `json:"user_id"`
}

func (u *UserModel) Login(password, name string, status int) (*UserInfo, error) {
	user := &UserInfo{}
	u.ctx.Log.Info("model_gameid", 1001, "model_gamename", "三国战纪")
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	md5Str := []byte(password)
	has := md5.Sum(md5Str)
	password = fmt.Sprintf("%x", has)
	if status == 1 {
		//注册
		insertSql := "INSERT INTO user (user_name, password, creat_time )VALUES( ?,?,?)"
		result, err := u.User.ctx.Mysql.Exec(insertSql, name, password, timeStr)
		if err != nil {
			return nil, errors.New("服务错误")
		}
		lastId, _ := result.LastInsertId()
		user.UserName = name
		user.UserId = int(lastId)
	} else {
		sql := "select id, user_name, password from  user where password = ? and user_name = ?  limit 1"
		result, err := u.User.ctx.Mysql.QueryRow(sql, password, name)
		if err != nil {
			return user, err
		}
		if len(result) == 0 {
			return user, errors.New("用户名字或者密码错误")
		}
		user.UserName = result["user_name"]
		user.UserId, _ = strconv.Atoi(result["id"])
	}

	return user, nil
}
