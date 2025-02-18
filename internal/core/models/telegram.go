package models

import "time"

type Telegram struct {
	TelegramId int64
	FirstName  string
	LastName   string
	Username   string
	ChatId     int64
	CreatedAt  time.Time
}
