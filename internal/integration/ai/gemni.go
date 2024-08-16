package ai

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/generative-ai-go/genai"
	"github.com/trend-me/ai-requester/internal/config/exceptions"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
)

const model = "gemini-pro"

type (
	GeminiClient interface {
		GenerativeModel(model string) GeminiModel
	}
	GeminiModel interface {
		GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
	}

	GeminiClientConstructor func(ctx context.Context, key string) (client GeminiClient, err error)

	GeminiKeysGetter func() []string
	gemini           struct {
		geminiKeys              GeminiKeysGetter
		geminiClientConstructor GeminiClientConstructor
		client                  GeminiClient
	}
)

func (g *gemini) createClient(ctx context.Context, key string) (err error) {
	g.client, err = g.geminiClientConstructor(ctx, key)
	if err != nil {
		slog.WarnContext(ctx, "gemini.createClient",
			slog.String("details", "error during client creation"),
			slog.String("error", err.Error()))
	}
	return
}

func (g *gemini) generateContent(ctx context.Context, prompt string) (resp *genai.GenerateContentResponse, err error) {
	model := g.client.GenerativeModel(model)
	resp, err = model.GenerateContent(ctx, genai.Text(prompt))
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
	geminiClientConstructor GeminiClientConstructor,
) interfaces.Ai {
	return &gemini{
		geminiKeys:              geminiKeys,
		geminiClientConstructor: geminiClientConstructor,
	}
}
