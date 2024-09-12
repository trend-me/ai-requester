package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/trend-me/ai-requester/internal/config/injector"
	"github.com/trend-me/ai-requester/internal/integration/connections"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file",
			slog.String("error", err.Error()),
		)
		return
	}

	ctx := context.Background()
	consumer, err := injector.InitializeQueueAiRequesterConsumer()
	if err != nil {
		slog.Error("Error initializing consumer",
			slog.String("error", err.Error()),
		)
		return
	}

	defer connections.Disconnect()

	for {
		_, _ = consumer.Consume(ctx)
	}

}
