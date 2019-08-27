package model

import (
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/context"
)

type BookModel struct {
	*context.RequestContext
}

func (b *BookModel) GetBookInfo() string {
	return "abc"
}
