package context

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	LoggerInfoBufferCap  = 20
	LoggerErrorBufferCap = 4
)

var (
	gpool   sync.Pool
	idIndex int64
)

func init() {
	gpool.New = func() interface{} {
		return &g{}
	}
}

//每次请求的上下文
type RequestContext struct {
	//基础上下文
	*BasicContext

	//公共上下行
	G *g

	//log
	Log *log
}

type log struct {
	//暂存log数据
	infoKV  []interface{}
	errorKV []interface{}
}

func (l *log) Info(kvs ...interface{}) {
	if l.infoKV == nil {
		l.infoKV = make([]interface{}, 0, LoggerInfoBufferCap)
	}
	l.infoKV = append(l.infoKV, kvs...)
}

func (l *log) Error(kvs ...interface{}) {
	if l.errorKV == nil {
		l.errorKV = make([]interface{}, 0, LoggerErrorBufferCap)
	}
	l.errorKV = append(l.errorKV, kvs...)
}

func (rcx *RequestContext) GetLog() (*[]interface{}, *[]interface{}) {
	return &rcx.Log.infoKV, &rcx.Log.errorKV
}

//初始化
func (rcx *RequestContext) Reset(ginctx *gin.Context) {
	if rcx.Log == nil {
		rcx.Log = &log{}
	}
	rcx.Log.errorKV = nil
	rcx.Log.infoKV = nil

	//初始化G公共上行
	greq := gpool.Get().(*g)
	rcx.G = greq

	//初始化__id
	greq.RequestID = ginctx.Query("__id")
	if greq.RequestID == "" {
		greq.RequestID = getId()
	}

	var k string
	for i := 1; i <= 20; i++ {
		k = fmt.Sprintf("p%d", i)
		v := ginctx.Query(k)
		switch i {
		case 1:
			greq.P1 = v
		case 2:
			greq.P2 = v
		case 3:
			greq.P3 = v
		case 4:
			greq.P4 = v
		case 5:
			greq.P5 = v
		case 6:
			greq.P6 = v
		case 7:
			greq.P7 = v
		case 8:
			greq.P8 = v
		case 9:
			greq.P9 = v
		case 10:
			greq.P10 = v
		case 11:
			greq.P11, _ = strconv.Atoi(v)
		case 12:
			greq.P12, _ = strconv.Atoi(v)
		case 13:
			greq.P13, _ = strconv.Atoi(v)
		case 14:
			greq.P14, _ = strconv.Atoi(v)
		case 15:
			greq.P15 = v
		case 16:
			greq.P16 = v
		case 17:
			greq.P17 = v
		case 18:
			greq.P18 = v
		case 19:
			if v == "1" {
				greq.P19 = true
			}
		case 20:
			greq.P20 = v
		}
	}

}

//请求公共上下行解析
//P(n)说明具体说明参照wiki：http://wiki.baidu-shucheng.com/pages/viewpage.action?pageId=4915578
type g struct {
	//request __id 每次请求的唯一标识id
	RequestID string
	//用户信息，待添加..

	//p系列公共上下行
	P1, P2, P3, P4, P5, P6, P7, P8, P9, P10 string
	P11, P12, P13, P14                      int
	P15, P16, P17, P18                      string
	P19                                     bool
	P20                                     string
}

func getId() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), atomic.AddInt64(&idIndex, 1))
}
