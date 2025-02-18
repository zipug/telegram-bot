package minio

import (
	"bot/internal/common/service/config"
	"bot/internal/core/models"
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioRepository struct {
	mc           *minio.Client
	buckets      []string
	num_workers  int
	url_lifetime time.Duration
}

var (
	ErrPing       = errors.New("could not ping MiniO")
	ErrUploadFile = errors.New("could not upload file")
	ErrUploadMany = errors.New("could not upload many files")
	ErrGetUrl     = errors.New("could not get file url")
	ErrGetMany    = errors.New("could not get many files")
	ErrDeleteFile = errors.New("could not delete file")
	ErrDeleteMany = errors.New("could not delete many files")
)

func NewMinioRepository(cfg *config.TgBotConfig) *MinioRepository {
	repo := &MinioRepository{}
	if err := repo.InvokeConnect(cfg); err != nil {
		e := fmt.Errorf(
			"MINIO: minio://%s:%s@%s:%d\nERROR: %w",
			cfg.MiniO.User,
			cfg.MiniO.Password,
			cfg.MiniO.Host,
			cfg.MiniO.Port,
			err)
		panic(e)
	}
	return repo
}

func (repo *MinioRepository) InvokeConnect(cfg *config.TgBotConfig) error {
	ctx := context.Background()

	client, err := minio.New(
		fmt.Sprintf("%s:%d", cfg.MiniO.Host, cfg.MiniO.Port),
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.MiniO.User, cfg.MiniO.Password, ""),
			Secure: cfg.MiniO.UseSsl,
		},
	)
	if err != nil {
		return err
	}
	repo.mc = client

	if err := repo.PingTest(); err != nil {
		panic(err)
	}

	repo.url_lifetime = cfg.MiniO.UrlLifetime
	repo.buckets = []string{cfg.MiniO.BucketArticles, cfg.MiniO.BucketAttachments, cfg.MiniO.BucketAvatars}
	for _, bucket := range repo.buckets {
		exists, err := repo.mc.BucketExists(ctx, bucket)
		if err != nil {
			return err
		}
		if !exists {
			if err := repo.mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return err
			}
			return err
		}
	}

	return nil
}

func (repo *MinioRepository) PingTest() error {
	max_errs := 5
	errs := 0
	timeout := 1 * time.Second
	for max_errs > 0 {
		if ok := repo.mc.IsOnline(); !ok {
			fmt.Println("could not ping database - MiniO is offline")
			fmt.Printf("retrying in %s\n", timeout)
			max_errs--
			errs++
			time.Sleep(timeout)
		}
		max_errs = 0
		errs = 0
	}
	if errs == 0 {
		return nil
	}
	return fmt.Errorf("%w: minio: %s", ErrPing, repo.mc.EndpointURL().String())
}

func (repo *MinioRepository) GetClient() *minio.Client {
	return repo.mc
}

func (repo *MinioRepository) UploadFile(ctx context.Context, file models.File) (models.MinioResponse, error) {
	obj_id := uuid.New().String()
	reader := bytes.NewReader(file.Data)
	if _, err := repo.mc.PutObject(
		ctx,
		file.Bucket,
		obj_id,
		reader,
		int64(len(file.Data)),
		minio.PutObjectOptions{ContentType: file.ContentType},
	); err != nil {
		return models.MinioResponse{}, err
	}
	return repo.GetFileUrl(ctx, obj_id, file.Bucket)
}

func (repo *MinioRepository) UploadManyFiles(ctx context.Context, files []models.File) (map[string]models.MinioResponse, []models.MinioErr) {
	urls := make(map[string]models.MinioResponse, len(files))
	errs := make([]models.MinioErr, 0, len(files))

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	urlsCh := make(chan models.MinioResponse, len(files))
	errsCh := make(chan models.MinioErr, len(files))

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := repo.UploadFile(cctx, file)
			if err != nil {
				errsCh <- models.MinioErr{Error: err, FileName: file.Name, Bucket: file.Bucket}
			}
			urlsCh <- url
		}()
	}

	wg.Wait()
	close(urlsCh)
	close(errsCh)
	for url := range urlsCh {
		urls[url.ObjectId] = url
	}
	for err := range errsCh {
		errs = append(errs, err)
	}

	return urls, errs
}

func (repo *MinioRepository) GetFileUrl(
	ctx context.Context,
	obj_id,
	bucket string,
) (models.MinioResponse, error) {
	url, err := repo.mc.PresignedGetObject(ctx, bucket, obj_id, repo.url_lifetime, nil)
	if err != nil {
		return models.MinioResponse{Url: "", ObjectId: obj_id}, err
	}
	return models.MinioResponse{Url: url.String(), ObjectId: obj_id}, nil
}

func (repo *MinioRepository) GetManyFileUrls(ctx context.Context, obj_ids []string, bucket string) (map[string]models.MinioResponse, []models.MinioErr) {
	urls := make(map[string]models.MinioResponse, len(obj_ids))
	errs := make([]models.MinioErr, 0, len(obj_ids))
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	urlsCh := make(chan models.MinioResponse, len(obj_ids))
	errsCh := make(chan models.MinioErr, len(obj_ids))

	for _, obj_id := range obj_ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url, err := repo.GetFileUrl(cctx, obj_id, bucket)
			if err != nil {
				errsCh <- models.MinioErr{Error: err, FileName: obj_id, Bucket: bucket}
			}
			urlsCh <- url
		}()
	}

	wg.Wait()
	close(urlsCh)
	close(errsCh)
	for url := range urlsCh {
		urls[url.ObjectId] = url
	}
	for err := range errsCh {
		errs = append(errs, err)
	}

	return urls, errs
}

func (repo *MinioRepository) DeleteFile(ctx context.Context, obj_id, bucket string) error {
	if err := repo.mc.RemoveObject(ctx, bucket, obj_id, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}

func (repo *MinioRepository) DeleteManyFiles(ctx context.Context, obj_ids []string, bucket string) []models.MinioErr {
	errs := make([]models.MinioErr, 0, len(obj_ids))
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	errsCh := make(chan models.MinioErr, len(obj_ids))

	for _, obj_id := range obj_ids {
		wg.Add(1)
		go func() {
			if err := repo.DeleteFile(cctx, obj_id, bucket); err != nil {
				errsCh <- models.MinioErr{Error: err, FileName: obj_id, Bucket: bucket}
			}
		}()
	}

	wg.Wait()
	close(errsCh)
	for err := range errsCh {
		errs = append(errs, err)
	}

	return errs
}
