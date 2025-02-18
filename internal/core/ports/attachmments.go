package ports

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"context"
)

type AttachmentsRepository interface {
	GetAllAttachmentsByArticleId(ctx context.Context, article_id int64) ([]dto.AttachmentDbo, error)
}

type AttachmentsService interface {
	GetAllAttachmentsByArticleId(article_id int64) ([]models.Attachment, error)
}
