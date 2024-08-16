package usecases

import (
	"context"
	"log/slog"

	"github.com/trend-me/ai-requester/internal/domain/factories"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/domain/models"
)

type UseCase struct {
	queueAiCallback                interfaces.QueueAiCallback
	aiFactory 					  *factories.AiFactory
}

func (u UseCase) Handle(ctx context.Context, request *models.Request) error {
	slog.InfoContext(ctx, "useCase.Handle",
		slog.String("details", "process started"))


	slog.DebugContext(ctx, "useCase.Handle",
		slog.String("details", "process finished"))
	return nil
}
 

func NewUseCase(queueAiCallback interfaces.QueueAiCallback, aiFactory *factories.AiFactory) interfaces.UseCase {
	return &UseCase{
		queueAiCallback:                queueAiCallback,
		aiFactory:		                aiFactory,
	}
}
