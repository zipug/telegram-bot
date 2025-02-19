package postgres

import (
	"bot/internal/application/dto"
	pu "bot/pkg/postgres_utils"
	"context"
	"errors"
)

var ErrCreateArticle = errors.New("could not create article")

func (repo *PostgresRepository) CreateArticle(ctx context.Context, article dto.ArticleDbo) (int64, error) {
	rows, err := pu.Dispatch[dto.ArticleDbo](
		ctx,
		repo.db,
		`
		INSERT INTO articles (name, description, content, project_id)
		VALUES ($1::text, $2::text, $3::text, $4::bigint)
		RETURNING *;
		`,
		article.Name,
		article.Description,
		article.Content,
		article.ProjectId,
	)
	if err != nil {
		return -1, err
	}
	if len(rows) == 0 {
		return -1, ErrCreateArticle
	}
	return rows[0].Id, nil
}
