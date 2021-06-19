package library

import (
	"crypto/md5"
	"fmt"
	"os"
)

type Func struct {
	*BasicLibrary
}

func (f *Func) GetGame() {
	f.ctx.Log.Info("func_gameid", true, "func_gamename", "斗地主")
	f.ctx.Log.Error("func_getgameinfo", "timeout", "func_response", "dial error")
}

func (f *Func) GetPage(pageId, pageSize int) int {
	return pageId * pageSize
}

func (f *Func) ExistsDir (path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (f *Func)Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}