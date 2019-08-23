package core

import (
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/basic"
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/context"
	"gitlab.baidu-shucheng.com/shaohua/bloc/logger"
	"sync"
)

type Core struct {
	ctx    *context.BasicContext
	model  *basic.Model
	logger *logger.Logger
	pool   sync.Pool
}

func NewCore(ctx *context.BasicContext) *Core {
	core := &Core{
		ctx:    ctx,
		logger: ctx.Logger,
		model:  basic.NewModel(ctx),
	}
	core.pool.New = func() interface{} {
		return &basic.Controller{
			BasicContext: core.ctx,
			Model:        core.model,
		}
	}

	return core
}
