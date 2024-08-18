//go:build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/trend-me/ai-requester/internal/config/connections"
	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/ai-requester/internal/delivery/controllers"
	"github.com/trend-me/ai-requester/internal/domain/factories"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/domain/usecases"
	"github.com/trend-me/ai-requester/internal/integration/api"
	"github.com/trend-me/ai-requester/internal/integration/queue"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
)

func newQueueConnectionAiRequesterConsumer(connection *rabbitmq.Connection) queue.ConnectionAiRequesterConsumer {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueAiRequester,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func newQueueConnectionAiCallback(connection *rabbitmq.Connection) queue.ConnectionAiCallback {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueAiCallback,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func newQueueAiRequesterConsumer(controller interfaces.Controller, connectionAiPromptBuilderConsumer queue.ConnectionAiRequesterConsumer) interfaces.QueueAiRequesterConsumer {
	return queue.NewAiPromptBuilderConsumer(connectionAiPromptBuilderConsumer, controller)
}

func urlApiValidationGetter() api.UrlApiValidation {
	return properties.UrlApiValidation
}

func urlApiPromptRoadMapConfigGetter() api.UrlApiPromptRoadMapConfig {
	return properties.UrlApiPromptRoadMapConfig
}

func InitializeQueueAiRequesterConsumerMock(geminiMock interfaces.Ai) (interfaces.QueueAiRequesterConsumer, error) {
	wire.Build(
		urlApiPromptRoadMapConfigGetter,
		urlApiValidationGetter,
		api.NewValidation,
		api.NewApiPromptRoadMapConfig,
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
