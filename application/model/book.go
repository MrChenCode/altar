package model

import (
	"altar/application/core/context"
)

type BookModel struct {
	*context.RequestContext
	Game *GameModel
}

func (b *BookModel) GetBookInfo() string {
	b.Log.Info("bookid", 75699, "bookname", "星辰变")
	b.Log.Error("res", "无效的bookid", "response", "timeout")
	return b.Game.GetGameInfo()
}
