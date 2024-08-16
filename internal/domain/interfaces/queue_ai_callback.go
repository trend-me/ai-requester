package interfaces

import (
	"context"

	"github.com/trend-me/ai-requester/internal/domain/models"
)

type QueueAiCallback interface {
	Publish(ctx context.Context, prompt string, request *models.Request) error
}
