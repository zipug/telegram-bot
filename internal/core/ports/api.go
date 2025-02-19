package ports

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
)

type ApiService interface {
	Search(project_id, url, question string) ([]models.Answer, error)
	AISearch(question string) (*dto.AIResponseDto, error)
}
