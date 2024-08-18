package mocks

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/trend-me/ai-requester/internal/integration/ai"
)

type GeminiModelMock struct {
	response string
	prompt   string
	err      error
}

func (g *GeminiModelMock) GetPrompt() string {
	return g.prompt
}

func (g *GeminiModelMock) Clean() {
	g.prompt = ""
	g.response = ""
	g.err = nil
}

func (g *GeminiModelMock) SetResponse(response string) {
	g.response = response
}

func (g *GeminiModelMock) SetError(err error) {
	g.err = err
}

func (g *GeminiModelMock) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
	for _, part := range parts {
		if text, ok := part.(genai.Text); ok {
			g.prompt = string(text)
		}
	}

	return &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text(g.response),
					},
				},
			},
		},
	}, nil
}

func NewGeminiModelMock() ai.GeminiModel {
	return &GeminiModelMock{}
}
