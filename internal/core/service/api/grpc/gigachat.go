package grpc

import (
	"bot/internal/application/dto"
	"bot/internal/core/ports"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "github.com/zipug/protos/gen/go/gigachat"
)

type GigaChatService struct {
	ctx          context.Context
	repo         ports.GigaChatRepository
	grpc_address string
	auth_url     string
	auth_key     string
	bearer       string
	scope        string
	model        string
	client       *grpc.ClientConn
	giga         pb.ChatServiceClient
}

func NewGigaChatService(
	ctx context.Context,
	repo ports.GigaChatRepository,
	grpc_address,
	auth_url,
	auth_key,
	scope,
	model string,
) *GigaChatService {
	service := &GigaChatService{
		ctx:          ctx,
		repo:         repo,
		grpc_address: grpc_address,
		auth_url:     auth_url,
		auth_key:     auth_key,
		scope:        scope,
		model:        model,
	}
	if err := service.connect(); err != nil {
		panic(err)
	}
	return service
}

func (g *GigaChatService) authGigaChat() error {
	client := http.Client{}
	ctx, cancel := context.WithTimeout(g.ctx, 60*time.Second)
	defer cancel()
	payload := strings.NewReader("scope=" + g.scope)
	req, err := http.NewRequestWithContext(ctx, "POST", g.auth_url, payload)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("RqUID", "6f0b1291-c7f3-43c6-bb2e-9f3efb2dc98e")
	req.Header.Add("Authorization", "Basic "+g.auth_key)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res dto.AuthResponse
	if err := json.Unmarshal(readBytes, &res); err != nil {
		return err
	}
	g.bearer = res.AccessToken
	return nil
}

func (g *GigaChatService) connect() error {
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return fmt.Errorf("error while loading system cert pool: %w", err)
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: caCertPool,
	})
	conn, err := grpc.NewClient(
		g.grpc_address,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return fmt.Errorf("error while connecting to gRPC server: %w", err)
	}
	g.client = conn
	g.giga = pb.NewChatServiceClient(g.client)
	return nil
}

func (g *GigaChatService) Close() {
	g.client.Close()
}

func (g *GigaChatService) Chat(message *pb.Message, telegram_id, project_id int64) (*pb.ChatResponse, error) {
	var messages []*pb.Message
	if err := g.authGigaChat(); err != nil {
		return nil, fmt.Errorf("error while authorization: %w", err)
	}
	md := metadata.Pairs("authorization", "Bearer "+g.bearer)
	ctx := metadata.NewOutgoingContext(g.ctx, md)

	allMessages, err := g.repo.GetAllDialogMessages(ctx, telegram_id, project_id)
	if err != nil {
		return nil, err
	}
	if allMessages != "" && allMessages != "[]" {
		var dbmsgs []map[string]string
		if err := json.Unmarshal([]byte(allMessages), &dbmsgs); err != nil {
			return nil, err
		}
		for _, msg := range dbmsgs {
			var role, content string
			for key, value := range msg {
				if key == "role" {
					role = value
				}
				if key == "content" {
					content = value
				}
			}
			messages = append(messages, &pb.Message{
				Role:    role,
				Content: content,
			})
		}
	}

	userMessage := dto.GigaChatMessage{
		Role:    message.Role,
		Content: message.Content,
	}
	jsonUserMessage, err := json.Marshal(userMessage)
	if err != nil {
		return nil, err
	}

	if err := g.repo.AddNewDialogMessage(
		ctx,
		telegram_id,
		project_id,
		jsonUserMessage,
	); err != nil {
		return nil, err
	}

	messages = append(messages, &pb.Message{
		Role:    message.Role,
		Content: message.Content,
	})

	resp, err := g.giga.Chat(ctx, &pb.ChatRequest{
		Model:    g.model,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	assistantMessage := dto.GigaChatMessage{
		Role:    resp.Alternatives[0].Message.Role,
		Content: resp.Alternatives[0].Message.Content,
	}
	jsonAssistantMessage, err := json.Marshal(assistantMessage)
	if err != nil {
		return nil, err
	}

	if err := g.repo.AddNewDialogMessage(
		ctx,
		telegram_id,
		project_id,
		jsonAssistantMessage,
	); err != nil {
		return nil, err
	}

	return resp, nil
}
