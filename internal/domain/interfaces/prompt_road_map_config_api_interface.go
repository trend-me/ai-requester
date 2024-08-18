package interfaces

import (
	"context"

	"github.com/trend-me/ai-requester/internal/domain/models"
)

type ApiPromptRoadMapConfig interface {
	GetPromptRoadMap(ctx context.Context, promptRoadMapConfigName string, promptRoadMapStep int) (*models.PromptRoadMap, error)
}
