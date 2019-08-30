//model基础配置

package model

import (
	"altar/application/context/rctx"
	"altar/application/library"
)

type BasicModel struct {
	*rctx.RequestContext
	Library *library.BasicLibrary

	Book *BookModel
	Game *GameModel
}

//初始化model
func NewModel(ctx *rctx.RequestContext, lib *library.BasicLibrary) *BasicModel {
	m := &BasicModel{RequestContext: ctx, Library: lib}

	m.Book = &BookModel{m}
	m.Game = &GameModel{m}

	return m
}

//每次请求重置model requestContext资源
//此处代码固定，不要修改
//func (m *Model) Reset(ctx *rctx.RequestContext) {
//	*m.rctx = *ctx
//}
