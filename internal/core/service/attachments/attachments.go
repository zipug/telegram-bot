package attachments

import (
	"bot/internal/core/models"
	"bot/internal/core/ports"
	"context"
)

type AttachmentsService struct {
	ctx  context.Context
	repo ports.AttachmentsRepository
}

func NewAttachmentService(ctx context.Context, repo ports.AttachmentsRepository) *AttachmentsService {
	return &AttachmentsService{
		ctx:  ctx,
		repo: repo,
	}
}

func (s *AttachmentsService) GetAllAttachmentsByArticleId(article_id int64) ([]models.Attachment, error) {
	dbo, err := s.repo.GetAllAttachmentsByArticleId(s.ctx, article_id)
	if err != nil {
		return nil, err
	}
	var res []models.Attachment
	for _, a := range dbo {
		res = append(res, a.ToValue())
	}
	return res, nil
}
