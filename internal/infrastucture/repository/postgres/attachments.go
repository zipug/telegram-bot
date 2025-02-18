package postgres

import (
	"bot/internal/application/dto"
	pu "bot/pkg/postgres_utils"
	"context"
)

func (repo *PostgresRepository) GetAllAttachmentsByArticleId(ctx context.Context, article_id int64) ([]dto.AttachmentDbo, error) {
	rows, err := pu.Dispatch[dto.AttachmentDbo](
		ctx,
		repo.db,
		`
		SELECT a.id, a.name, a.description, a.object_id, a.mimetype, a.user_id
    FROM attachments a
		LEFT JOIN attachments_articles aa ON a.id = aa.attachment_id
		WHERE aa.article_id = $1::bigint;
		`,
		article_id,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
