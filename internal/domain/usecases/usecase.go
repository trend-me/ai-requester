package usecases

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/trend-me/ai-requester/internal/config/exceptions"
	"github.com/trend-me/ai-requester/internal/domain/builders"
	"github.com/trend-me/ai-requester/internal/domain/factories"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/domain/models"
	"github.com/trend-me/ai-requester/internal/domain/parsers"
)

type UseCase struct {
	apiPromptRoadMapConfig interfaces.ApiPromptRoadMapConfig
	apiValidation          interfaces.ApiValidation
	queueAiCallback        interfaces.QueueAiCallback
	aiFactory              *factories.AiFactory
}

func (u UseCase) Handle(ctx context.Context, request *models.Request) error {
	slog.InfoContext(ctx, "useCase.Handle",
		slog.String("details", "process started"))

	ai, err := u.aiFactory.FactoryAi(request.Model)
	if err != nil {
		return err
	}

	promptRoadMap, err := u.apiPromptRoadMapConfig.GetPromptRoadMap(ctx, request.PromptRoadMapConfigName, request.PromptRoadMapStep)
	if err != nil {
		return err
	}

	aiResponse, err := ai.Prompt(ctx, request.Prompt)
	if err != nil {
		return err
	}

	response, err := parsers.ParseAiResponseToJSON(aiResponse)
	if err != nil {
		return err
	}

	u.validateMetadata(ctx, promptRoadMap, request)

	request.Metadata = builders.BuildMetadata(request.Metadata, response)

	err = u.queueAiCallback.Publish(ctx, request)
	if err != nil {
		return err
	}

	slog.DebugContext(ctx, "useCase.Handle",
		slog.String("details", "process finished"))
	return nil
}

func (u UseCase) validateMetadata(ctx context.Context, promptRoadMap *models.PromptRoadMap, request *models.Request) error {
	payload, err := json.Marshal(request.Metadata)
	if err != nil {
		return exceptions.NewValidationError(err.Error())
	}

	payloadValidationExecutionResponse, err := u.apiValidation.ExecutePayloadValidator(ctx, promptRoadMap.ResponseValidationName, payload)
	if err != nil {
		return err
	}

	bPayloadValidationExecutionResponse, _ := json.Marshal(payloadValidationExecutionResponse)
	slog.InfoContext(ctx, "useCase.Handle",
		slog.String("details", "ai resopnse validation"),
		slog.String("result", string(bPayloadValidationExecutionResponse)))

	if payloadValidationExecutionResponse.Failures != nil && len(*payloadValidationExecutionResponse.Failures) > 0 {
		return exceptions.NewAiResponseValidationError(*payloadValidationExecutionResponse.Failures)
	}
	return nil
}

func NewUseCase(
	queueAiCallback interfaces.QueueAiCallback,
	aiFactory *factories.AiFactory,
	apiPromptRoadMapConfig interfaces.ApiPromptRoadMapConfig,
	apiValidation interfaces.ApiValidation) interfaces.UseCase {
	return &UseCase{
		queueAiCallback:        queueAiCallback,
		aiFactory:              aiFactory,
		apiPromptRoadMapConfig: apiPromptRoadMapConfig,
		apiValidation:          apiValidation,
	}
}
