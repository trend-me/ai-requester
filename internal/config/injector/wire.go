//go:build wireinject

package injector

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/wire"
	"github.com/trend-me/ai-requester/internal/config/connections"
	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/ai-requester/internal/delivery/controllers"
	"github.com/trend-me/ai-requester/internal/domain/factories"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/domain/usecases"
	"github.com/trend-me/ai-requester/internal/integration/ai"
	"github.com/trend-me/ai-requester/internal/integration/queue"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
	"google.golang.org/api/option"
)

func NewQueueConnectionAiRequesterConsumer(connection *rabbitmq.Connection) queue.ConnectionAiRequesterConsumer {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueNameAiPromptBuilder,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func NewQueueConnectionAiCallback(connection *rabbitmq.Connection) queue.ConnectionAiCallback {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueAiRequester,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func NewQueueAiRequesterConsumer(controller interfaces.Controller, connectionAiPromptBuilderConsumer queue.ConnectionAiRequesterConsumer) interfaces.QueueAiRequesterConsumer {
	return queue.NewAiPromptBuilderConsumer(connectionAiPromptBuilderConsumer, controller)
}

func GeminiModelConstructor() ai.GeminiModelConstructor {
	return func(ctx context.Context, key string) (ai.GeminiModel, error) {
		client, err := genai.NewClient(context.Background(), option.WithAPIKey(key))
		if err != nil {
			return nil, err
		}
		return client.GenerativeModel(properties.GeminiModel), nil
	}
}

func GeminiKeysGetter() ai.GeminiKeysGetter {
	return properties.AiGeminiKeys
}

func InitializeQueueAiRequesterConsumer() (interfaces.QueueAiRequesterConsumer, error) {
	wire.Build(
		ai.NewGemini,
		GeminiKeysGetter,
		GeminiModelConstructor,
		controllers.NewController,
		factories.NewAiFactory,
		usecases.NewUseCase,
		queue.NewAiRequester,
		NewQueueConnectionAiCallback,
		NewQueueConnectionAiRequesterConsumer,
		connections.ConnectQueue,
		NewQueueAiRequesterConsumer)
	return nil, nil
}
