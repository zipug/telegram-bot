package telegramusers

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"bot/internal/core/ports"
	"context"
	"errors"
)

type TelegramUsersService struct {
	ctx  context.Context
	repo ports.TelegramUsersRepository
}

func NewTelegramUsersService(ctx context.Context, repo ports.TelegramUsersRepository) *TelegramUsersService {
	return &TelegramUsersService{ctx: ctx, repo: repo}
}

func (t *TelegramUsersService) AddTelegramUser(user models.Telegram) error {
	dbo := dto.ToDbo(user)
	id, err := t.repo.AddTelegramUser(t.ctx, dbo)
	if err != nil {
		return err
	}
	if id == -1 {
		return errors.New("could not add telegram user")
	}
	return nil
}
