package ports

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"context"
)

type StatisticsService interface {
	AddStatisticRecord(
		record models.Statistic,
	) error
}

type StatisticsRepository interface {
	AddStatisticRecord(
		ctx context.Context,
		record dto.StatisticDto,
	) (int64, error)
}
