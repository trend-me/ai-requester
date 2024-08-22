package exceptions

func NewValidationError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Validation Error",
		Message:   messages,
	}
}

func NewUnknownError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    false,
		ErrorType: "Unknown Error",
		Message:   messages,
	}
}

func NewAiFactoryError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Metadata Validation Error",
		Message:   messages,
	}
}

func NewQueueError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Queue Error",
		Message:   messages,
	}
}

func NewAiResponseError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "AI Response Error",
		Message:   messages,
	}
}

func NewAiError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "AI Error",
		Message:   messages,
	}
}

func NewAiResponseValidationError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Metadata Validation Error",
		Message:   messages,
	}
}

func NewPromptRoadMapNotFoundError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Prompt Road Map Not Found Error",
		Message:   messages,
	}
}

func NewGetPromptRoadMapConfigError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Get Prompt Road Map Config Error",
		Message:   messages,
	}
}

func NewPayloadValidatorNotFoundError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Prompt Road Map Config Error",
		Message:   messages,
	}
}

func NewPayloadValidatorError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Payload Validator Error",
		Message:   messages,
	}
}
