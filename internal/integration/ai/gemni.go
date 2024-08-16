package ai

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/generative-ai-go/genai"
	"github.com/trend-me/ai-requester/internal/config/exceptions"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
)

type (
	GeminiModel interface {
		GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
	}

	GeminiModelConstructor func(ctx context.Context, key string) (GeminiModel, error)

	GeminiKeysGetter func() []string
	gemini           struct {
		geminiKeys             GeminiKeysGetter
		geminiModelConstructor GeminiModelConstructor
		model                  GeminiModel
	}
)

func (g *gemini) createClient(ctx context.Context, key string) (err error) {
	g.model, err = g.geminiModelConstructor(ctx, key)
	if err != nil {
		slog.WarnContext(ctx, "gemini.createClient",
			slog.String("details", "error during client creation"),
			slog.String("error", err.Error()))
	}
	return
}

func (g *gemini) generateContent(ctx context.Context, prompt string) (resp *genai.GenerateContentResponse, err error) {
	resp, err = g.model.GenerateContent(ctx, genai.Text(prompt))
	return
}

func (g *gemini) Prompt(ctx context.Context, prompt string) (result string, err error) {
	var resp *genai.GenerateContentResponse
	for _, key := range g.geminiKeys() {
		err = g.createClient(ctx, key)
		if err != nil {
			continue
		}

		resp, err = g.generateContent(ctx, prompt)
		if err != nil {
			continue
		}
		break
	}

	if err != nil {
		slog.WarnContext(ctx, "gemini.Prompt",
			slog.String("details", "error during content generation"),
			slog.String("error", err.Error()))

		err = exceptions.NewAiError(err.Error())
		return
	}

	part := resp.Candidates[0].Content.Parts[0]
	str := fmt.Sprintf("%s", part)
	result = str
	return
}

func NewGemini(
	geminiKeys GeminiKeysGetter,
	geminiModelConstructor GeminiModelConstructor,
) interfaces.Ai {
	return &gemini{
		geminiKeys:             geminiKeys,
		geminiModelConstructor: geminiModelConstructor,
	}
}
