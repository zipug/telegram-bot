package ports

import (
	"context"

	pb "github.com/zipug/protos/gen/go/gigachat"
)

type GigaChatService interface {
	Chat(message *pb.Message, telegram_id, project_id int64) (*pb.ChatResponse, error)
	Close()
}

type GigaChatRepository interface {
	AddNewDialogMessage(ctx context.Context, telegram_id, project_id int64, content []byte) error
	GetAllDialogMessages(ctx context.Context, telegram_id, project_id int64) (string, error)
}
