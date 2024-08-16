package ai

import (
	"context"

	"github.com/trend-me/ai-requester/internal/domain/interfaces"
)

type (
	gemini struct{}
)

func (g *gemini) Prompt(ctx context.Context, prompt string) (string, error) {
	return "Gemini", nil
}

func NewGemini() interfaces.Ai {
	return &gemini{}
}
