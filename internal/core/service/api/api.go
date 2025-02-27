package api

import (
	"context"
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
