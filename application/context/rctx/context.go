//每次请求时，使用的上下文，他应该是一个对象池

package rctx

import (
	"altar/application/context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	idIndex int64
)

//每次请求的上下文, 每一个http请求，都是不一样的上下文对象
type RequestContext struct {
	//基础上下文
	*context.Context

	//公共上下行
	G *g

	//log
	Log *log
}

//获取缓冲的日志数据(k-v)
func (rcx *RequestContext) GetLog() ([]interface{}, []interface{}) {
	return rcx.Log.infoKV, rcx.Log.errorKV
}

//初始化
func (rcx *RequestContext) Reset(ginctx *gin.Context) {
	if rcx.Log == nil {
		rcx.Log = &log{}
	} else {
		rcx.Log.reset()
	}

	//初始化G公共上行
	if rcx.G == nil {
		rcx.G = &g{}
	} else {
		rcx.G.reset()
	}

	//初始化__id
	rcx.G.RequestID = ginctx.Query("__id")
	if rcx.G.RequestID == "" {
		rcx.G.RequestID = getId()
	}

	var k string
	for i := 1; i <= 20; i++ {
		k = fmt.Sprintf("p%d", i)
		v := ginctx.Query(k)
		switch i {
		case 1:
			rcx.G.P1 = v
		case 2:
			rcx.G.P2 = v
		case 3:
			rcx.G.P3 = v
		case 4:
			rcx.G.P4 = v
		case 5:
			rcx.G.P5 = v
		case 6:
			rcx.G.P6 = v
		case 7:
			rcx.G.P7 = v
		case 8:
			rcx.G.P8 = v
		case 9:
			rcx.G.P9 = v
		case 10:
			rcx.G.P10 = v
		case 11:
			rcx.G.P11, _ = strconv.Atoi(v)
		case 12:
			rcx.G.P12, _ = strconv.Atoi(v)
		case 13:
			rcx.G.P13, _ = strconv.Atoi(v)
		case 14:
			rcx.G.P14, _ = strconv.Atoi(v)
		case 15:
			rcx.G.P15 = v
		case 16:
			rcx.G.P16 = v
		case 17:
			rcx.G.P17 = v
		case 18:
			rcx.G.P18 = v
		case 19:
			if v == "1" {
				rcx.G.P19 = true
			}
		case 20:
			rcx.G.P20 = v
		}
	}
}

func getId() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), atomic.AddInt64(&idIndex, 1))
}
