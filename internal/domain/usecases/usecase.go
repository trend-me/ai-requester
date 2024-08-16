package usecases

import (
	"context"
	"log/slog"

	"github.com/trend-me/ai-requester/internal/domain/builders"
	"github.com/trend-me/ai-requester/internal/domain/factories"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/domain/models"
	"github.com/trend-me/ai-requester/internal/domain/parsers"
)

type UseCase struct {
	queueAiCallback interfaces.QueueAiCallback
	aiFactory       *factories.AiFactory
}

func (u UseCase) Handle(ctx context.Context, request *models.Request) error {
	slog.InfoContext(ctx, "useCase.Handle",
		slog.String("details", "process started"))

	ai, err := u.aiFactory.FactoryAi(request.Model)
	if err != nil {
		return err
	}

	//todo: fetch prompt_road_map from api

	aiResponse, err := ai.Prompt(ctx, request.Prompt)
	if err != nil {
		return err
	}

	response, err := parsers.ParseAiResponseToJSON(aiResponse)
	if err != nil {
		return err
	}

	//todo: validate response with api

	request.Metadata = builders.BuildMetadata(request.Metadata, response)

	err = u.queueAiCallback.Publish(ctx, request)
	if err != nil {
		return err
	}

	slog.DebugContext(ctx, "useCase.Handle",
		slog.String("details", "process finished"))
	return nil
}

func NewUseCase(queueAiCallback interfaces.QueueAiCallback, aiFactory *factories.AiFactory) interfaces.UseCase {
	return &UseCase{
		queueAiCallback: queueAiCallback,
		aiFactory:       aiFactory,
	}
}
