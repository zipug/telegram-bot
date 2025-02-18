package postgres

import (
	"bot/internal/application/dto"
	pu "bot/pkg/postgres_utils"
	"context"
	"errors"
)

var ErrStatisticAdd = errors.New("could not add new statistic record")

func (repo *PostgresRepository) AddStatisticRecord(ctx context.Context, record dto.StatisticDto) (int64, error) {
	var sql string
	var params []interface{}
	if record.ArticleId.Int64 == 0 {
		sql = `
			INSERT INTO statistics (bot_id, telegram_id, question, article_id, article_name, is_resolved)
			VALUES ($1::bigint, $2::bigint, $3::text, NULL, $4::text, $5::boolean)
			RETURNING id;
		`
		params = []interface{}{
			record.BotId,
			record.TelegramId,
			record.Question,
			record.ArticleName,
			record.IsResolved,
		}
	} else {
		sql = `
			INSERT INTO statistics (bot_id, telegram_id, question, article_id, article_name, is_resolved)
			VALUES ($1::bigint, $2::bigint, $3::text, $4::bigint, $5::text, $6::boolean)
			RETURNING id;
		`
		params = []interface{}{
			record.BotId,
			record.TelegramId,
			record.Question,
			record.ArticleId,
			record.ArticleName,
			record.IsResolved,
		}
	}
	rows, err := pu.Dispatch[dto.StatisticDto](
		ctx,
		repo.db,
		sql,
		params...,
	)
	if err != nil {
		return -1, err
	}
	if len(rows) == 0 {
		return -1, ErrStatisticAdd
	}
	rec := rows[0]
	return rec.Id, nil
}
