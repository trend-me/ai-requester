package factories

import (
	"fmt"

	"github.com/trend-me/ai-requester/internal/config/exceptions"
	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
)

type (
	AiFactory struct {
		Gemini interfaces.Ai
	}
)

func (f *AiFactory) FactoryAi(model string) (ai interfaces.Ai, err error) {
	m := map[string]interfaces.Ai{
		properties.AiModelNameGemini: f.Gemini,
	}

	ai = m[model]
	if ai == nil {
		return nil, exceptions.NewAiFactoryError(fmt.Sprintf("model %s not implemented", model))
	}
	return
}

func NewAiFactory(gemini interfaces.Ai) *AiFactory {
	return &AiFactory{
		Gemini: gemini,
	}
}
