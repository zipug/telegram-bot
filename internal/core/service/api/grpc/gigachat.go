package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
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
	ctx      context.Context
	address  string
	auth_url string
	auth_key string
	bearer   string
	scope    string
	client   *grpc.ClientConn
	pb.UnimplementedChatServiceServer
}

type ChatService interface {
	Chat(context.Context, *pb.ChatRequest) (*pb.ChatResponse, error)
	ChatStream(*pb.ChatRequest, grpc.ServerStreamingServer[pb.ChatResponse]) error
}

func NewGigaChatService(ctx context.Context, address, auth_url, auth_key, scope string) *GigaChatService {
	return &GigaChatService{
		ctx:      ctx,
		address:  address,
		auth_url: auth_url,
		auth_key: auth_key,
		scope:    scope,
	}
}

type AuthPayload struct {
	Mode       string     `json:"mode"`
	Urlencoded Urlencoded `json:"urlencoded"`
}

type Urlencoded struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (g *GigaChatService) AuthGigaChat() error {
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
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	fmt.Println("authorization into the giga", resp.Body)

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var res AuthResponse
	if err := json.Unmarshal(readBytes, &res); err != nil {
		return err
	}

	g.bearer = res.AccessToken

	fmt.Println(res.AccessToken)

	return nil
}

func (g *GigaChatService) unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Add the Authorization header to the metadata
	md := metadata.Pairs("authorization", "Bearer "+g.bearer)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Invoke the RPC call with the updated context
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (g *GigaChatService) Connect() (*grpc.ClientConn, error) {
	if err := g.AuthGigaChat(); err != nil {
		fmt.Println("Error while authorization: ", err)
		return nil, err
	}
	fmt.Println(g.bearer)
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Println("Error while loading system cert pool: ", err)
		return nil, err
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: caCertPool,
	})
	md := metadata.Pairs("authorization", "Bearer "+g.bearer)
	ctx := metadata.NewOutgoingContext(g.ctx, md)
	conn, err := grpc.NewClient(
		g.address,
		grpc.WithTransportCredentials(creds),
		// grpc.WithUnaryInterceptor(unaryInterceptor),
	)
	if err != nil {
		fmt.Println("Error while connecting to gRPC server: ", err)
	}
	defer conn.Close()
	fmt.Println("Connected to gRPC server", conn.Target())
	giga := pb.NewChatServiceClient(conn)
	resp, err := giga.Chat(ctx, &pb.ChatRequest{
		Model: "GigaChat",
		Messages: []*pb.Message{
			{
				Role:    "user",
				Content: "Привет, расскажи о себе",
			},
		},
	})
	if err != nil {
		fmt.Println("Error while calling gRPC server: ", err)
		return nil, err
	}
	alts := resp.Alternatives
	if len(alts) == 0 {
		fmt.Println("No alternatives found")
		return nil, errors.New("no alternatives found")
	}
	msg := alts[0].Message
	fmt.Printf("%s\n", msg.Content)
	return conn, err
}
