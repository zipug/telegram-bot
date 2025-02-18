package ports

import "bot/internal/core/models"

type ApiService interface {
	Search(project_id, url, question string) ([]models.Answer, error)
}
