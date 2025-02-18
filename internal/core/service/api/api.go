package api

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type ApiService struct {
	ctx context.Context
}

func NewApiService(ctx context.Context) *ApiService {
	return &ApiService{ctx: ctx}
}

func (s *ApiService) Search(project_id, url, question string) ([]models.Answer, error) {
	client := http.Client{}

	ctx, cancel := context.WithTimeout(s.ctx, 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("id", project_id)
	q.Add("query", question)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res dto.SearchDto
	if err := json.Unmarshal(readBytes, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
