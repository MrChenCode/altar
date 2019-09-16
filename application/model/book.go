package model

type BookModel struct {
	*BasicModel
}

func (b *BookModel) GetBookInfo(bookid, channelid int) string {
	return ""
}
