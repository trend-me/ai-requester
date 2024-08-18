package exceptions

func NewValidationError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Validation Error",
		message:   messages,
	}
}

func NewUnknownError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    false,
		ErrorType: "Unknown Error",
		message:   messages,
	}
}

func NewAiFactoryError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Metadata Validation Error",
		message:   messages,
	}
}

func NewQueueError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Queue Error",
		message:   messages,
	}
}

func NewAiResponseError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "AI Response Error",
		message:   messages,
	}
}

func NewAiError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "AI Error",
		message:   messages,
	}
}

func NewAiResponseValidationError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Metadata Validation Error",
		message:   messages,
	}
}

func NewPromptRoadMapNotFoundError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Prompt Road Map Not Found Error",
		message:   messages,
	}
}

func NewGetPromptRoadMapConfigError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Get Prompt Road Map Config Error",
		message:   messages,
	}
}

func NewPayloadValidatorNotFoundError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Prompt Road Map Config Error",
		message:   messages,
	}
}

func NewPayloadValidatorError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Payload Validator Error",
		message:   messages,
	}
}
