package ai

import "github.com/trend-me/ai-requester/internal/domain/interfaces"

type (
 gemini struct {}
)

func (g *gemini) Prompt(prompt string) (string, error){
	return "Gemini", nil
}	

func NewGemini() interfaces.Ai{
	return &gemini{}
}