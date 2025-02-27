package api

import (
	"bot/internal/application/dto"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

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
