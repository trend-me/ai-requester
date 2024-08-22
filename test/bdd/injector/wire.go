//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/ai-requester/internal/delivery/controllers"
	"github.com/trend-me/ai-requester/internal/domain/factories"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/domain/usecases"
	"github.com/trend-me/ai-requester/internal/integration/apis"
	"github.com/trend-me/ai-requester/internal/integration/connections"
	"github.com/trend-me/ai-requester/internal/integration/queues"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
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

func InitializeQueueAiRequesterConsumerMock(geminiMock interfaces.Ai) (interfaces.QueueAiRequesterConsumer, error) {
	wire.Build(
		urlApiPromptRoadMapConfigGetter,
		urlApiValidationGetter,
		apis.NewValidation,
		apis.NewApiPromptRoadMapConfig,
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
