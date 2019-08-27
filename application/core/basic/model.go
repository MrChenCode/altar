package basic

import (
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/context"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/model"
)

type Model struct {
	rcx  *context.RequestContext
	Book *model.BookModel
}

func NewModel(rcx *context.RequestContext) *Model {
	m := &Model{
		rcx: rcx,
	}
	m.Book = &model.BookModel{RequestContext: m.rcx}
	return m
}

func (m *Model) Reset(rcx *context.RequestContext) {
	m.rcx = rcx
}

//反射扫描是否存在初始化model函数Init
//初始化仅调用一次
//func callModelInit(x *Model) {
//	v := reflect.ValueOf(x)
//	fn := v.Elem().NumField()
//	for i := 0; i < fn; i++ {
//		vv := v.Elem().Field(i)
//		m := vv.MethodByName("Init")
//		if !m.IsValid() {
//			continue
//		}
//
//		if m.Type().NumIn() != 0 {
//			panic(fmt.Sprintf("%s %s Too many parameters.", vv.String(), m.Type().String()))
//		}
//		m.Call(nil)
//	}
//}
