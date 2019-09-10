package library

import "altar/application/context/rctx"

type BasicLibrary struct {
	ctx *rctx.RequestContext

	Func *Func
}

func NewLibrary(ctx *rctx.RequestContext) *BasicLibrary {
	l := &BasicLibrary{ctx: ctx}

	l.Func = &Func{l}

	return l
}
