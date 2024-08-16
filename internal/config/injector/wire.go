//go:build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/trend-me/ai-requester/internal/config/connections"
	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/ai-requester/internal/delivery/controllers"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/domain/usecases"
	"github.com/trend-me/ai-requester/internal/integration/queue"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
)

func NewQueueAiPromptBuilderConsumerConnection(connection *rabbitmq.Connection) queue.ConnectionAiRequesterConsumer {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueNameAiPromptBuilder,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func NewQueueAiRequesterConnection(connection *rabbitmq.Connection) queue.ConnectionAiCallback {
	return rabbitmq.NewQueue(
		connection,
		properties.QueueAiRequester,
		rabbitmq.ContentTypeJson,
		properties.CreateQueueIfNX(),
		true,
		true,
	)
}

func NewConsumer(controller interfaces.Controller, connectionAiPromptBuilderConsumer queue.ConnectionAiRequesterConsumer) interfaces.QueueAiRequesterConsumer {
	return queue.NewAiPromptBuilderConsumer(connectionAiPromptBuilderConsumer, controller)
}


func InitializeConsumer() (interfaces.QueueAiRequesterConsumer, error) {
	wire.Build(controllers.NewController,
		usecases.NewUseCase,
		queue.NewAiRequester,
		NewQueueAiRequesterConnection,
		NewQueueAiPromptBuilderConsumerConnection,
		connections.ConnectQueue,
		NewConsumer)
	return nil, nil
}
