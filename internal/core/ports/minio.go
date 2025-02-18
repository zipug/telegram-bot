package ports

import (
	"bot/internal/core/models"
	"context"

	"github.com/minio/minio-go/v7"
)

type MinioService interface {
	UploadFile(context.Context, models.File) (models.MinioResponse, error)
	UploadManyFiles(context.Context, []models.File) (map[string]models.MinioResponse, error)
	GetFileUrl(context.Context, string, string) (models.MinioResponse, error)
	GetManyFileUrls(context.Context, []string, string) (map[string]models.MinioResponse, error)
	DeleteFile(context.Context, string, string) error
	DeleteManyFiles(context.Context, []string, string) error
	GetClient() *minio.Client
}

type MinioRepository interface {
	UploadFile(context.Context, models.File) (models.MinioResponse, error)
	UploadManyFiles(context.Context, []models.File) (map[string]models.MinioResponse, []models.MinioErr)
	GetFileUrl(context.Context, string, string) (models.MinioResponse, error)
	GetManyFileUrls(context.Context, []string, string) (map[string]models.MinioResponse, []models.MinioErr)
	DeleteFile(context.Context, string, string) error
	DeleteManyFiles(context.Context, []string, string) []models.MinioErr
	GetClient() *minio.Client
}
