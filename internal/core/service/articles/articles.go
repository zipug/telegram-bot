package articles

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"bot/internal/core/ports"
	"context"
	"errors"
)

type ArticlesService struct {
	ctx  context.Context
	repo ports.ArticleRepository
}

func NewArticlesService(ctx context.Context, repo ports.ArticleRepository) *ArticlesService {
	return &ArticlesService{ctx: ctx, repo: repo}
}

func (s *ArticlesService) CreateArticle(article models.Article) error {
	articleDbo := dto.ToArticleDbo(article)
	article_id, err := s.repo.CreateArticle(s.ctx, articleDbo)
	if err != nil {
		return err
	}
	if article_id == -1 {
		return errors.New("could not create article")
	}
	return nil
}
