package models

type Statistic struct {
	BotId       int64
	TelegramId  int64
	ArticleId   int64
	Question    string
	ArticleName string
	IsResolved  bool
	ParentId    int64
}
