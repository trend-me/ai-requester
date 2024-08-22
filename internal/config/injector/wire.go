//go:build wireinject

package injector

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/wire"
	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/ai-requester/internal/delivery/controllers"
	"github.com/trend-me/ai-requester/internal/domain/factories"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/domain/usecases"
	"github.com/trend-me/ai-requester/internal/integration/ais"
	"github.com/trend-me/ai-requester/internal/integration/apis"
	"github.com/trend-me/ai-requester/internal/integration/connections"
	"github.com/trend-me/ai-requester/internal/integration/queues"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
	"google.golang.org/api/option"
)

func newQueueConnectionAiRequesterConsumer(connection *rabbitmq.Connection) queues.ConnectionAiRequesterConsumer {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueAiRequester,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func newQueueConnectionAiCallback(connection *rabbitmq.Connection) queues.ConnectionAiCallback {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueAiCallback,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func newQueueAiRequesterConsumer(controller interfaces.Controller, connectionAiPromptBuilderConsumer queues.ConnectionAiRequesterConsumer) interfaces.QueueAiRequesterConsumer {
	return queues.NewAiPromptBuilderConsumer(connectionAiPromptBuilderConsumer, controller)
}

func geminiModelConstructor() ais.GeminiModelConstructor {
	return func(ctx context.Context, key string) (ais.GeminiModel, error) {
		client, err := genai.NewClient(context.Background(), option.WithAPIKey(key))
		if err != nil {
			return nil, err
		}
		return client.GenerativeModel(properties.GeminiModel), nil
	}
}

func geminiKeysGetter() ais.GeminiKeysGetter {
	return properties.AiGeminiKeys
}

func urlApiValidationGetter() apis.UrlApiValidation {
	return properties.UrlApiValidation
}

func urlApiPromptRoadMapConfigGetter() apis.UrlApiPromptRoadMapConfig {
	return properties.UrlApiPromptRoadMapConfig
}

func newQueueConnectionOutputGetter(connection *rabbitmq.Connection) queues.ConnectionOutputGetter {
	return func(queueName string) queues.ConnectionOutput {
		return rabbitmq.NewQueue(
			connection,
			queueName,
			rabbitmq.ContentTypeJson,
			properties.CreateQueueIfNX(),
			true,
			true,
		)
	}
}

func InitializeQueueAiRequesterConsumer() (interfaces.QueueAiRequesterConsumer, error) {
	wire.Build(
		urlApiPromptRoadMapConfigGetter,
		urlApiValidationGetter,
		apis.NewValidation,
		apis.NewApiPromptRoadMapConfig,
		ais.NewGemini,
		geminiKeysGetter,
		geminiModelConstructor,
		controllers.NewController,
		factories.NewAiFactory,
		newQueueConnectionOutputGetter,
		queues.NewOutput,
		usecases.NewUseCase,
		queues.NewAiRequester,
		newQueueConnectionAiCallback,
		newQueueConnectionAiRequesterConsumer,
		connections.ConnectQueue,
		newQueueAiRequesterConsumer)
	return nil, nil
}
