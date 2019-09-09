//model基础配置

package model

import (
	"altar/application/context/rctx"
	"altar/application/library"
)

type modelContext struct {
	ctx     *rctx.RequestContext
	library *library.BasicLibrary
}

type BasicModel struct {
	*modelContext
	Book *BookModel
	Game *GameModel
}

//初始化model
func NewModel(ctx *rctx.RequestContext, lib *library.BasicLibrary) *BasicModel {
	mctx := &modelContext{ctx, lib}
	m := &BasicModel{modelContext: mctx}

	m.Book = &BookModel{mctx}
	m.Game = &GameModel{mctx}

	return m
}

//每次请求重置model requestContext资源
//此处代码固定，不要修改
//func (m *Model) Reset(ctx *rctx.RequestContext) {
//	*m.rctx = *ctx
//}
