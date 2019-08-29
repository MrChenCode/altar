package library

import "altar/application/context/rctx"

type Library struct {
	*rctx.RequestContext

	Func *Func
}

func NewLibrary(ctx *rctx.RequestContext) *Library {
	l := &Library{RequestContext: ctx}

	l.Func = &Func{l}

	return l
}

//重置每次请求上下文
//固定代码，请勿修改
//func (l *Library) Reset(ctx *rctx.RequestContext) {
//	*l.rctx = *ctx
//}
