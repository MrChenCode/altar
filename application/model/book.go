package model

import (
	"gitlab.baidu-shucheng.com/shaohua/bloc/application/core/context"
)

type BookModel struct {
	*context.BasicContext
}

func (b *BookModel) GetBookInfo() string {
	return "abcd"
}
