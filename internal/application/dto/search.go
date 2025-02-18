package dto

import "bot/internal/core/models"

type SearchDto struct {
	Data []models.Answer `json:"data"`
}
