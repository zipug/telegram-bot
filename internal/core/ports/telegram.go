package ports

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"context"
)

type TelegramUsersService interface {
	AddTelegramUser(user models.Telegram) error
}

type TelegramUsersRepository interface {
	AddTelegramUser(ctx context.Context, user dto.TelegramDbo) (int64, error)
}
