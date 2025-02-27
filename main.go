package main

import (
	"bot/internal/application"
	"bot/internal/common/service/config"
	logger "bot/internal/common/service/logger/zerolog"
	"bot/internal/core/service/api"
	"bot/internal/core/service/articles"
	"bot/internal/core/service/attachments"
	"bot/internal/core/service/minio"
	"bot/internal/core/service/statistics"
	telegramusers "bot/internal/core/service/telegram_users"
	repo_minio "bot/internal/infrastucture/repository/minio"
	"bot/internal/infrastucture/repository/postgres"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	defer stop()
	cfg := config.NewConfigService()
	fmt.Printf("CURRENT_CONFIG: %v", cfg)
	/*
		gigachatService := grpc.NewGigaChatService(
			ctx,
			"gigachat.devices.sberbank.ru",
			"https://ngw.devices.sberbank.ru:9443/api/v2/oauth",
			"MzIwM2RkNDEtMjJhZS00MzAyLTgzZDEtYWEzODQzZDg0ZDIxOjhkNTY2MWM0LTkxY2YtNGZiZS04NGY0LTM0NGQ2YmQ1NzM0NQ==",
			"GIGACHAT_API_PERS",
		)
		gigachatService.Connect()
	*/
	minioRepository := repo_minio.NewMinioRepository(cfg)
	minioService := minio.NewMinioService(minioRepository)
	postgresRepository := postgres.NewPostgresRepository(cfg)
	api := api.NewApiService(ctx, cfg.OpenRouterAi.Token, cfg.OpenRouterAi.Model, cfg.OpenRouterAi.URL)
	attachmentsService := attachments.NewAttachmentService(ctx, postgresRepository)
	statisticsService := statistics.NewStatisticsService(ctx, postgresRepository)
	tgUsersService := telegramusers.NewTelegramUsersService(ctx, postgresRepository)
	articlesService := articles.NewArticlesService(ctx, postgresRepository)
	logger := logger.New(cfg.Env)
	app := application.New(
		ctx,
		cfg,
		logger,
		minioService,
		postgresRepository,
		api,
		attachmentsService,
		tgUsersService,
		articlesService,
		statisticsService,
	)
	app.Run()
	select {
	case <-ctx.Done():
		fmt.Println("shutting down...")
	}
}
