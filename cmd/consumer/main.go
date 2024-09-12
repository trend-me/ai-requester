package main

import (
	"context"
	"log/slog"

	"github.com/trend-me/ai-requester/internal/config/injector"
	"github.com/trend-me/ai-requester/internal/integration/connections"
)

func main() {
	ctx := context.Background()
	consumer, err := injector.InitializeQueueAiRequesterConsumer()
	if err != nil {
		slog.Error("Error initializing consumer",
			slog.String("error", err.Error()),
		)
		return
	}

	_, _ = consumer.Consume(ctx)

	defer connections.Disconnect()
}
