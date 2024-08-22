package models

import "github.com/trend-me/ai-requester/internal/config/exceptions"

type Request struct {
	PromptRoadMapConfigName        string
	PromptRoadMapStep              int
	PromptRoadMapConfigExecutionId string
	OutputQueue                    string
	Model                          string
	Error                          *exceptions.ErrorType
	Prompt                         string
	Metadata                       map[string]any
}
