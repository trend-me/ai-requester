package queue

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/ai-requester/internal/domain/interfaces"
)

type (
	ConnectionAiRequesterConsumer interface {
		Consume(ctx context.Context, handler func(delivery amqp.Delivery) error) (chan error, error)
	}

	aiRequesterConsumer struct {
		queue      ConnectionAiRequesterConsumer
		controller interfaces.Controller
	}
)

func (a aiRequesterConsumer) Consume(ctx context.Context) (chan error, error) {
	return a.queue.Consume(ctx, a.controller.Handle)
}

func NewAiPromptBuilderConsumer(queue ConnectionAiRequesterConsumer, controller interfaces.Controller) interfaces.QueueAiRequesterConsumer {
	return &aiRequesterConsumer{queue: queue, controller: controller}
}
