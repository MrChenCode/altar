package basic

import (
	"github.com/gin-gonic/gin"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/context"
	"strconv"
)

type Controller struct {
	//基础的数据资源
	*context.BasicContext

	//gin请求资源
	*gin.Context

	//model
	Model *Model

	//请求公共上行
	G *G
}

func (c *Controller) Reset(ginc *gin.Context) {
	c.Context = ginc

	g := &G{}
	c.G = g
	for i := 1; i <= 20; i++ {
		k := "p" + strconv.Itoa(i)
		v := ginc.Query(k)
		switch i {
		case 1:
			g.P1 = v
		case 2:
			g.P2 = v
		case 3:
			g.P3 = v
		case 4:
			g.P4 = v
		case 5:
			g.P5 = v
		case 6:
			g.P6 = v
		case 7:
			g.P7 = v
		case 8:
			g.P8 = v
		case 9:
			g.P9 = v
		case 10:
			g.P10 = v
		case 11:
			g.P11, _ = strconv.Atoi(v)
		case 12:
			g.P12, _ = strconv.Atoi(v)
		case 13:
			g.P13, _ = strconv.Atoi(v)
		case 14:
			g.P14, _ = strconv.Atoi(v)
		case 15:
			g.P15 = v
		case 16:
			g.P16 = v
		case 17:
			g.P17 = v
		case 18:
			g.P18 = v
		case 19:
			if v == "1" {
				g.P19 = true
			}
		case 20:
			g.P20 = v
		}
	}
}

//请求公共上下行解析
//P(n)说明具体说明参照wiki：http://wiki.baidu-shucheng.com/pages/viewpage.action?pageId=4915578
type G struct {
	//用户信息，待添加..

	//p系列公共上下行
	P1  string
	P2  string
	P3  string
	P4  string
	P5  string
	P6  string
	P7  string
	P8  string
	P9  string
	P10 string
	P11 int
	P12 int
	P13 int
	P14 int
	P15 string
	P16 string
	P17 string
	P18 string
	P19 bool
	P20 string
}
