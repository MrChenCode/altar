package core

import (
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/basic"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/context"
	"sync"
)

type Core struct {
	ctx *context.BasicContext

	controllerPool sync.Pool
	modelPool      sync.Pool
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
		}
	}
	core.modelPool.New = func() interface{} {
		return basic.NewModel()
	}

	return core
}
