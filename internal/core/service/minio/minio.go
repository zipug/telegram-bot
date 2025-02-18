package minio

import (
	"bot/internal/core/models"
	"bot/internal/core/ports"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
)

type MinioService struct {
	repo ports.MinioRepository
}

func NewMinioService(repo ports.MinioRepository) *MinioService {
	return &MinioService{repo: repo}
}

func (m *MinioService) GetClient() *minio.Client {
	return m.repo.GetClient()
}

func (m *MinioService) UploadFile(ctx context.Context, file models.File) (models.MinioResponse, error) {
	return m.repo.UploadFile(ctx, file)
}

func (m *MinioService) UploadManyFiles(ctx context.Context, files []models.File) (map[string]models.MinioResponse, error) {
	urls, errs := m.repo.UploadManyFiles(ctx, files)
	if len(errs) == 0 {
		return urls, nil
	}
	err_array := make([]string, 0, len(errs))
	for _, err := range errs {
		err_array = append(err_array, fmt.Sprintf("error: %v, file: %s, bucket: %s", err.Error, err.FileName, err.Bucket))
	}
	return urls, errors.New(strings.Join(err_array, "\n"))
}

func (m *MinioService) GetFileUrl(ctx context.Context, object_id, bucket string) (models.MinioResponse, error) {
	return m.repo.GetFileUrl(ctx, object_id, bucket)
}

func (m *MinioService) GetManyFileUrls(ctx context.Context, object_ids []string, bucket string) (map[string]models.MinioResponse, error) {
	urls, errs := m.repo.GetManyFileUrls(ctx, object_ids, bucket)
	if len(errs) == 0 {
		return urls, nil
	}
	err_array := make([]string, 0, len(errs))
	for _, err := range errs {
		err_array = append(err_array, fmt.Sprintf("error: %v, file: %s, bucket: %s", err.Error, err.FileName, err.Bucket))
	}
	return urls, errors.New(strings.Join(err_array, "\n"))
}

func (m *MinioService) DeleteFile(ctx context.Context, object_id, bucket string) error {
	return m.repo.DeleteFile(ctx, object_id, bucket)
}

func (m *MinioService) DeleteManyFiles(ctx context.Context, object_ids []string, bucket string) error {
	errs := m.repo.DeleteManyFiles(ctx, object_ids, bucket)
	err_array := make([]string, 0, len(errs))
	for _, err := range errs {
		err_array = append(err_array, fmt.Sprintf("error: %v, file: %s, bucket: %s", err.Error, err.FileName, err.Bucket))
	}
	return errors.New(strings.Join(err_array, "\n"))
}
