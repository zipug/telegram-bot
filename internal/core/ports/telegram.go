package ports

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"context"
)

type TelegramUsersService interface {
	AddTelegramUser(user models.Telegram, project_id int64) error
}

type TelegramUsersRepository interface {
	AddTelegramUser(ctx context.Context, user dto.TelegramDbo, project_id int64) (int64, error)
}
