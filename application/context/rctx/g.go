package rctx

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

//g初始化
func (a *g) reset() {
	a.RequestID = ""
	a.P1, a.P2, a.P3, a.P4, a.P5 = "", "", "", "", ""
	a.P6, a.P7, a.P8, a.P9, a.P10 = "", "", "", "", ""
	a.P11, a.P12, a.P13, a.P14 = 0, 0, 0, 0
	a.P15, a.P16, a.P17, a.P18 = "", "", "", ""
	a.P19, a.P20 = false, ""
}
