package core

import (
	"altar/application/core/basic"
	"altar/application/core/context"
	"sync"
)

type Core struct {
	ctx *context.BasicContext

	controllerPool sync.Pool
}

func NewCore(ctx *context.BasicContext) *Core {
	core := &Core{
		ctx: ctx,
	}
	core.controllerPool.New = func() interface{} {
		return &basic.Controller{
			RequestContext: &context.RequestContext{
				BasicContext: core.ctx,
			},
			Model: basic.NewModel(),
		}
	}

	return core
}
