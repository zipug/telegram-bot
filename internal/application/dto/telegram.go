package dto

import (
	"bot/internal/core/models"
	"database/sql"
)

type TelegramDbo struct {
	TelegramId int64          `db:"telegram_id"`
	FirstName  sql.NullString `db:"first_name,omitempty"`
	LastName   sql.NullString `db:"last_name,omitempty"`
	Username   sql.NullString `db:"username,omitempty"`
	ChatId     int64          `db:"chat_id,omitempty"`
	CreatedAt  sql.NullTime   `db:"created_at,omitempty"`
}

func (a *TelegramDbo) ToValue() models.Telegram {
	return models.Telegram{
		TelegramId: a.TelegramId,
		FirstName:  a.FirstName.String,
		LastName:   a.LastName.String,
		Username:   a.Username.String,
		ChatId:     a.ChatId,
		CreatedAt:  a.CreatedAt.Time,
	}
}

func ToDbo(a models.Telegram) TelegramDbo {
	return TelegramDbo{
		TelegramId: a.TelegramId,
		FirstName:  sql.NullString{String: a.FirstName, Valid: true},
		LastName:   sql.NullString{String: a.LastName, Valid: true},
		Username:   sql.NullString{String: a.Username, Valid: true},
		ChatId:     a.ChatId,
		CreatedAt:  sql.NullTime{Time: a.CreatedAt, Valid: true},
	}
}
