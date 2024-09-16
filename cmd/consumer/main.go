package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	errChan, err := consumer.Consume(ctx)
	if err != nil {
		log.Fatalf("Error initializing consumer: %v", err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	case sig := <-sigChan:
		log.Printf("Received signal: %v. Shutting down...", sig)
	}

	defer connections.Disconnect()

}
