package statistics

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"bot/internal/core/ports"
	"context"
	"errors"
)

type StatisticsService struct {
	ctx  context.Context
	repo ports.StatisticsRepository
}

func NewStatisticsService(ctx context.Context, repo ports.StatisticsRepository) *StatisticsService {
	return &StatisticsService{ctx: ctx, repo: repo}
}

func (t *StatisticsService) AddStatisticRecord(record models.Statistic) error {
	dbo := dto.ToStatisticDbo(record)
	id, err := t.repo.AddStatisticRecord(t.ctx, dbo)
	if err != nil {
		return err
	}
	if id == -1 {
		return errors.New("could not add statistic record")
	}
	return nil
}
