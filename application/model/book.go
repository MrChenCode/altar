package model

type BookModel struct {
	*modelContext
}

func (b *BookModel) GetBookInfo() string {
	b.ctx.Log.Info("model_bookid", 2001, "model_bookname", "吞噬星空")
	b.ctx.Log.Info("model_response", "false", "model_bookres", 0)
	b.ctx.Log.Error("model_bookinfo", "null", "model_bookid", 2001, "model_bookrequest", "timeout")
	return "吞噬星空"
}
