package library

import "altar/application/context/rctx"

type libraryContext struct {
	ctx *rctx.RequestContext
}

type BasicLibrary struct {
	*libraryContext

	Func *Func
}

func NewLibrary(ctx *rctx.RequestContext) *BasicLibrary {
	libctx := &libraryContext{ctx: ctx}
	l := &BasicLibrary{libraryContext: libctx}

	l.Func = &Func{libctx}

	return l
}

//重置每次请求上下文
//固定代码，请勿修改
//func (l *Library) Reset(ctx *rctx.RequestContext) {
//	*l.rctx = *ctx
//}
