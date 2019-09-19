package util

import (
	"altar/application/config"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//把v转换成int类型，如果转换失败，会返回0
func Int(v interface{}) int {
	n, err := strconv.Atoi(fmt.Sprint(v))
	if err != nil {
		return 0
	}
	return n
}

//把v转换int64类型, 转换失败会返回0
func Int64(v interface{}) int64 {
	n, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
	if err != nil {
		return 0
	}
	return n
}

//返回v到string拷贝
func String(v interface{}) string {
	return fmt.Sprint(v)
}

//返回v到string的拷贝，但是会剔除前后的空白字符
//包括\t、\n、\v、\f、\r、空格、删除符号、连续空格符(nbsp)
func TrimString(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

//把v转换成float32，转换失败会返回0
func Float32(v interface{}) float32 {
	n, err := strconv.ParseFloat(String(v), 32)
	if err != nil {
		return 0
	}
	return float32(n)
}

//把v转换成float64，转换失败会返回0
func Float64(v interface{}) float64 {
	n, err := strconv.ParseFloat(String(v), 64)
	if err != nil {
		return 0
	}
	return n
}

//判断一个interface是否为true
//
//nil--false
//整型、浮点型、复数 为0/0.000--false
//空字符串--false(剔除空格制表符等)
//chan/func/map/ptr/pointer/interface/slice 非nil返回true
//map/slice/array 非nil长度为0返回false
//非指针struct返回true
func Bool(v interface{}) bool {
	if v == nil {
		return false
	}
	f := reflect.ValueOf(v)
	if !f.IsValid() {
		return false
	}
	kind := f.Kind()
	switch kind {
	case reflect.Bool:
		return f.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return f.Int() != 0
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return f.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		return f.Complex() != 0
	case reflect.String:
		return strings.TrimSpace(f.String()) != ""
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.UnsafePointer, reflect.Slice:
		x := f.IsNil()
		if !x {
			if kind == reflect.Map || kind == reflect.Slice {
				return f.Len() != 0
			}
		}
		return !x
	case reflect.Ptr, reflect.Interface:
		if f.IsNil() {
			return false
		}
		return Bool(f.Elem().Interface())
	case reflect.Array:
		return f.Len() != 0
	case reflect.Struct:
		return true
	default:
		return true
	}
}

//对字符串s，从start位置截取，截取length长度, 此函数按照byte计算，如果要按照字符计算，请使用MbSubstr
//如果start为正数>=0，则从头开始数，如果start为负数<0，则从尾部开始数, 如果超过s长度，返回空字符串
//如果length > 0 则为截取的长度，如果为0，默认截取到末尾, 如果为负数返回空字符串，如果超过s长度，返回有效截取的字符串
func Substr(s string, start, length int) string {
	n := len(s)
	if n == 0 || length < 0 {
		return ""
	}
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
		endIndex = startIndex + length
		if endIndex > n {
			endIndex = n
		}
	} else {
		//length < 0 情况开始已经判断
		endIndex = n
	}
	return string(s[startIndex:endIndex])
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
	if s == "" {
		return ""
	}
	pre := Substr(s, 0, 7)
	if pre == "https:/" || pre == "http://" {
		return s
	}
	return config.CdnDomainRoute + s
}
