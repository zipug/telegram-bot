package postgres

import (
	"bot/internal/application/dto"
	pu "bot/pkg/postgres_utils"
	"context"
	"errors"
	"strings"
)

var ErrTgUserAdd = errors.New("could not add telegram user")

func (repo *PostgresRepository) AddTelegramUser(ctx context.Context, user dto.TelegramDbo) (int64, error) {
	rows, err := pu.Dispatch[dto.TelegramDbo](
		ctx,
		repo.db,
		`
		INSERT INTO telegram_users (telegram_id, first_name, last_name, username, chat_id)
		VALUES ($1::bigint, $2::text, $3::text, $4::text, $5::bigint)
		RETURNING telegram_id;
		`,
		user.TelegramId,
		user.FirstName,
		user.LastName,
		user.Username,
		user.ChatId,
	)
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
			return user.TelegramId, nil
		}
		return -1, err
	}
	if len(rows) == 0 {
		return -1, ErrTgUserAdd
	}
	usr := rows[0]
	if usr.TelegramId != user.TelegramId {
		return -1, ErrTgUserAdd
	}
	return usr.TelegramId, nil
}
