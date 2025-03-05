package dto

import (
	"bot/internal/core/models"
	"database/sql"
)

type StatisticDto struct {
	Id          int64          `db:"id"`
	BotId       int64          `db:"bot_id"`
	TelegramId  int64          `db:"telegram_id"`
	Question    sql.NullString `db:"question"`
	ArticleId   sql.NullInt64  `db:"article_id"`
	ArticleName string         `db:"article_name"`
	IsResolved  bool           `db:"is_resolved"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	ParentId    sql.NullInt64  `db:"parent_id,omitempty"`
}

func (d StatisticDto) ToValue() models.Statistic {
	return models.Statistic{
		BotId:       d.BotId,
		TelegramId:  d.TelegramId,
		ArticleId:   d.ArticleId.Int64,
		Question:    d.Question.String,
		ArticleName: d.ArticleName,
		IsResolved:  d.IsResolved,
		ParentId:    d.ParentId.Int64,
	}
}

func ToStatisticDbo(s models.Statistic) StatisticDto {
	return StatisticDto{
		BotId:       s.BotId,
		TelegramId:  s.TelegramId,
		ArticleId:   sql.NullInt64{Int64: s.ArticleId, Valid: true},
		Question:    sql.NullString{String: s.Question, Valid: true},
		ArticleName: s.ArticleName,
		IsResolved:  s.IsResolved,
		ParentId:    sql.NullInt64{Int64: s.ParentId, Valid: true},
	}
}
