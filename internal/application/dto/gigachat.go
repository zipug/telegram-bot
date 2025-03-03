package dto

import (
	"bot/internal/core/models"
	"database/sql"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type GigaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GigaChatDbo struct {
	TelegramId int64          `db:"telegram_id"`
	ProjectId  int64          `db:"project_id"`
	Dialog     sql.NullString `db:"dialog,omitempty"`
}

func (d *GigaChatDbo) ToValue() models.GigaChat {
	return models.GigaChat{
		TelegramId: d.TelegramId,
		ProjectId:  d.ProjectId,
		Dialog:     d.Dialog.String,
	}
}

func ToGigaChatDbo(m models.GigaChat) GigaChatDbo {
	return GigaChatDbo{
		TelegramId: m.TelegramId,
		ProjectId:  m.ProjectId,
		Dialog:     sql.NullString{String: m.Dialog, Valid: true},
	}
}
