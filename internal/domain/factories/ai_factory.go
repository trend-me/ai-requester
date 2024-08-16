package factories

import (
	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
)



type (
	AiFactory struct {
		Gemini interfaces.Ai
	}
)

func (f AiFactory) FactoryAi(model string) interfaces.Ai {
	m := map[string]interfaces.Ai{
		properties.ModelGemini: f.Gemini,
	}
	return m[model]
}
