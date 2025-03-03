package postgres

import (
	"bot/internal/application/dto"
	pu "bot/pkg/postgres_utils"
	"context"
	"errors"
	"strings"
)

var ErrTgUserAdd = errors.New("could not add telegram user")

func (repo *PostgresRepository) AddTelegramUser(ctx context.Context, user dto.TelegramDbo, project_id int64) (int64, error) {
	tx := repo.db.MustBegin()
	rows, err := pu.DispatchTx[dto.TelegramDbo](
		ctx,
		tx,
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
			tx.Rollback()
			return user.TelegramId, nil
		}
		tx.Rollback()
		return -1, err
	}
	if len(rows) == 0 {
		tx.Rollback()
		return -1, ErrTgUserAdd
	}
	usr := rows[0]
	if usr.TelegramId != user.TelegramId {
		tx.Rollback()
		return -1, ErrTgUserAdd
	}
	rows_giga, err := pu.DispatchTx[dto.GigaChatDbo](
		ctx,
		tx,
		`
		INSERT INTO telegram_dialogs (telegram_id, project_id, dialog)
		VALUES ($1::bigint, $2::bigint, '[]'::jsonb)
		RETURNING telegram_id;
		`,
		user.TelegramId,
		project_id,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	if len(rows_giga) == 0 {
		tx.Rollback()
		return -1, ErrTgUserAdd
	}
	if rows_giga[0].TelegramId != user.TelegramId {
		tx.Rollback()
		return -1, ErrTgUserAdd
	}
	tx.Commit()
	return usr.TelegramId, nil
}
