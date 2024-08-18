package interfaces

import (
	"context"

	"github.com/trend-me/ai-requester/internal/domain/models"
)

type ApiValidation interface {
	ExecutePayloadValidator(ctx context.Context, payloadValidatorId string, payload []byte) (*models.PayloadValidatorExecutionResponse, error)
}
