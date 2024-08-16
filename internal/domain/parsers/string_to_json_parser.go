package parsers

import (
	"encoding/json"

	"github.com/trend-me/ai-requester/internal/config/exceptions"
)

func ParseAiResponseToJSON(s string) (map[string]interface{}, error) {
	var result map[string]interface{}

	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, exceptions.NewAiResponseError("failed to parse AI response to JSON: " + err.Error())
	}

	return result, nil
}
