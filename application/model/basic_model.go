//model基础配置

package model

import (
	"altar/application/context/rctx"
	"altar/application/library"
)

type BasicModel struct {
	ctx     *rctx.RequestContext
	library *library.BasicLibrary

	Book        *BookModel
	Game        *GameModel
	Topic       *TopicModel
	Cloudsshelf *CloudsshelfModel
	User        *UserModel
	Fish        *FishModel
}

//初始化model
func NewModel(ctx *rctx.RequestContext, lib *library.BasicLibrary) *BasicModel {
	m := &BasicModel{ctx: ctx, library: lib}

	m.Book = &BookModel{m}
	m.Game = &GameModel{m}
	m.Topic = &TopicModel{m}
	m.Cloudsshelf = &CloudsshelfModel{m}
	m.User = &UserModel{m}
	m.Fish = &FishModel{m}

	return m
}
