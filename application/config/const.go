package config

//书籍类型
type BookType = int

const (
	//默认类型，网络书
	BookTypeDefault BookType = 0

	//epub出版物
	BookTypeEpub BookType = 1

	//comic漫画
	BookTypeComic BookType = 2

	//audio音频
	BookTypeAudio BookType = 3
)
