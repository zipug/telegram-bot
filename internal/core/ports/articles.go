package ports

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"context"
)

type ArticlesService interface {
	CreateArticle(models.Article) error
}

type ArticleRepository interface {
	CreateArticle(context.Context, dto.ArticleDbo) (int64, error)
}
