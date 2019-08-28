package basic

import (
	"altar/application/core/context"
	"altar/application/model"
)

type Model struct {
	rcx  *context.RequestContext
	Book *model.BookModel
}

func NewModel() *Model {
	m := &Model{
		rcx: &context.RequestContext{},
	}
	m.Book = &model.BookModel{RequestContext: m.rcx, Game: &model.GameModel{m.rcx}}
	return m
}

func (m *Model) Reset(rcx *context.RequestContext) {
	*m.rcx = *rcx
}