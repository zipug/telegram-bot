package dto

import (
	"bot/internal/core/models"
	"database/sql"
)

type ArticleDbo struct {
	Id          int64          `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Content     string         `db:"content"`
	ProjectId   int64          `db:"project_id"`
	CreatedAt   sql.NullTime   `db:"created_at,omitempty"`
	UpdateAt    sql.NullTime   `db:"updated_at,omitempty"`
	DeleteAt    sql.NullTime   `db:"deleted_at,omitempty"`
}

func (a *ArticleDbo) ToValue() models.Article {
	return models.Article{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description.String,
		Content:     a.Content,
		ProjectId:   a.ProjectId,
	}
}

func ToArticleDbo(a models.Article) ArticleDbo {
	return ArticleDbo{
		Id:          a.Id,
		Name:        a.Name,
		Description: sql.NullString{String: a.Description, Valid: true},
		Content:     a.Content,
		ProjectId:   a.ProjectId,
	}
}
