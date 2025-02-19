package api

import (
	"bot/internal/application/dto"
	"bot/internal/core/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type ApiService struct {
	ctx          context.Context
	ai_api_token string
	ai_api_model string
	ai_api_url   string
}

func NewApiService(ctx context.Context, api_token, api_model, api_url string) *ApiService {
	return &ApiService{ctx: ctx, ai_api_token: api_token, ai_api_model: api_model, ai_api_url: api_url}
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

func (s *ApiService) AISearch(question string) (*dto.AIResponseDto, error) {
	client := http.Client{}

	ctx, cancel := context.WithTimeout(s.ctx, 60*time.Second)
	defer cancel()

	payload := dto.AIPayloadDto{
		Model: s.ai_api_model,
		Messages: []dto.AIMessage{
			{
				Role:    "user",
				Content: question,
			},
		},
	}
	reqData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.ai_api_url, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.ai_api_token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res dto.AIResponseDto
	if err := json.Unmarshal(readBytes, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
