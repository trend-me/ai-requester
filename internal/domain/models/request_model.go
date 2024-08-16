package models

type Request struct {
	PromptRoadMapConfigName        string
	PromptRoadMapStep              int
	PromptRoadMapConfigExecutionId string
	OutputQueue                    string
	Model                          string
	Prompt                         string
	Metadata                       map[string]any
}
