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
	"github.com/trend-me/ai-requester/internal/integration/api"
	"github.com/trend-me/ai-requester/internal/integration/queue"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
	"google.golang.org/api/option"
)

func newQueueConnectionAiRequesterConsumer(connection *rabbitmq.Connection) queue.ConnectionAiRequesterConsumer {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueNameAiPromptBuilder,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func newQueueConnectionAiCallback(connection *rabbitmq.Connection) queue.ConnectionAiCallback {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueAiRequester,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func newQueueAiRequesterConsumer(controller interfaces.Controller, connectionAiPromptBuilderConsumer queue.ConnectionAiRequesterConsumer) interfaces.QueueAiRequesterConsumer {
	return queue.NewAiPromptBuilderConsumer(connectionAiPromptBuilderConsumer, controller)
}

func geminiModelConstructor() ai.GeminiModelConstructor {
	return func(ctx context.Context, key string) (ai.GeminiModel, error) {
		client, err := genai.NewClient(context.Background(), option.WithAPIKey(key))
		if err != nil {
			return nil, err
		}
		return client.GenerativeModel(properties.GeminiModel), nil
	}
}

func geminiKeysGetter() ai.GeminiKeysGetter {
	return properties.AiGeminiKeys
}

func urlApiValidationGetter() api.UrlApiValidation {
	return properties.UrlApiValidation
}

func urlApiPromptRoadMapConfigGetter() api.UrlApiPromptRoadMapConfig {
	return properties.UrlApiPromptRoadMapConfig
}

func InitializeQueueAiRequesterConsumer() (interfaces.QueueAiRequesterConsumer, error) {
	wire.Build(
		urlApiPromptRoadMapConfigGetter,
		urlApiValidationGetter,
		api.NewValidation,
		api.NewApiPromptRoadMapConfig,
		ai.NewGemini,
		geminiKeysGetter,
		geminiModelConstructor,
		controllers.NewController,
		factories.NewAiFactory,
		usecases.NewUseCase,
		queue.NewAiRequester,
		newQueueConnectionAiCallback,
		newQueueConnectionAiRequesterConsumer,
		connections.ConnectQueue,
		newQueueAiRequesterConsumer)
	return nil, nil
}
