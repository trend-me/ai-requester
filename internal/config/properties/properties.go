package properties

import (
	"context"
	"os"
	"strconv"
)

type contextKey = string

const ctxReceiveCount          contextKey = "ctxReceiveCount"
const (
	QueueNameAiPromptBuilder            = "ai-requester"
	QueueAiRequester                    = "ai-requester"
	ModelGemini = "gemini"
)

func CreateQueueIfNX() bool {
	return os.Getenv("CREATE_QUEUE_IF_NX") == "true"
}

func QueueConnectionUser() string {
	return os.Getenv("QUEUE_CONNECTION_USER")
}

func QueueConnectionPort() string {
	return os.Getenv("QUEUE_CONNECTION_PORT")
}

func QueueConnectionHost() string {
	return os.Getenv("QUEUE_CONNECTION_HOST")
}

func QueueConnectionPassword() string {
	return os.Getenv("QUEUE_CONNECTION_PASSWORD")
}

func GetMaxReceiveCount() int {
	i, _ := strconv.Atoi(os.Getenv("MAX_RECEIVE_COUNT"))
	return i
}

func SetCtxRetryCount(ctx context.Context, receiveCount int) context.Context {
	return context.WithValue(ctx, ctxReceiveCount, receiveCount)
}

func GetCtxRetryCount(ctx context.Context) int {
	i, _ := ctx.Value(ctxReceiveCount).(int)
	return i
}
