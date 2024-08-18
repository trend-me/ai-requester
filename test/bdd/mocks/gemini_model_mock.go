package mocks

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/trend-me/ai-requester/internal/integration/ai"
)

type GeminiModelMock struct {
	response string
	err      error
}

func (g *GeminiModelMock) SetResponse(response string) {
	g.response = response
}

func (g *GeminiModelMock) SetError(err error) {
	g.err = err
}

func (g *GeminiModelMock) GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error) {
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
