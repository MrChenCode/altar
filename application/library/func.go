package library

import (
	"altar/application/config"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type Func struct {
	*BasicLibrary
}

func (f *Func) GetGame() {
	f.ctx.Log.Info("func_gameid", true, "func_gamename", "斗地主")
	f.ctx.Log.Error("func_getgameinfo", "timeout", "func_response", "dial error")
}

func Int(v interface{}) int {
	n, err := strconv.Atoi(fmt.Sprint(v))
	if err != nil {
		return 0
	}
	return n
}

func Float64(v interface{}) float64 {
	n, err := strconv.ParseFloat(String(v), 64)
	if err != nil {
		return 0
	}
	return n
}

func String(v interface{}) string {
	return fmt.Sprint(v)
}

func TrimString(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

//对字符串s，从start位置截取，截取length长度, 此函数按照byte计算，如果要按照字符计算，请使用MbSubstr
//如果start为正数>=0，则从头开始数，如果start为负数<0，则从尾部开始数, 如果超过s长度，返回空字符串
//如果length > 0 则为截取的长度，如果为0，默认截取到末尾, 如果为负数返回空字符串，如果超过s长度，返回有效截取的字符串
func Substr(s string, start, length int) string {
	if s == "" || length < 0 {
		return s
	}
	b := *(*[]byte)(unsafe.Pointer(&s))
	n := len(b)
	var startIndex, endIndex int
	if start >= 0 {
		startIndex = start
	} else {
		startIndex = n - (start * -1)
	}
	if startIndex >= n {
		return ""
	}
	if length > 0 {
		endIndex = start + length
		if endIndex > n {
			endIndex = n
		}
	} else {
		//length < 0 情况开始已经判断
		endIndex = n
	}
	return string(b[startIndex:endIndex])
}

func GetBookType(bookid interface{}) config.BookType {
	idsub := Int(Substr(String(bookid), -3, 0))
	switch {
	case idsub > 800:
		return config.BookTypeEpub
	case idsub == 699:
		return config.BookTypeComic
	case idsub == 698:
		return config.BookTypeAudio
	default:
		return config.BookTypeDefault
	}
}

func FbkImg(v interface{}) string {
	s := TrimString(v)
	pre := Substr(s, 0, 7)
	if pre == "https:/" || pre == "http://" {
		return s
	}
	return "https://img.xmkanshu.com/novel/" + s
}
