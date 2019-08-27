package core

import (
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/basic"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/context"
	"sync"
)

type Core struct {
	ctx *context.BasicContext
	//model          *basic.Model
	controllerPool sync.Pool
	modelPool      sync.Pool
	libraryPool    sync.Pool
}

func NewCore(ctx *context.BasicContext) *Core {
	core := &Core{
		ctx: ctx,
		//model:  basic.NewModel(ctx),
	}
	core.controllerPool.New = func() interface{} {
		return &basic.Controller{
			RequestContext: &context.RequestContext{BasicContext: core.ctx},
		}
	}
	var emptyRequestContext = &context.RequestContext{}
	core.modelPool.New = func() interface{} {
		return basic.NewModel(emptyRequestContext)
	}

	return core
}
