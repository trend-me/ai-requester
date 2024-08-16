package interfaces

import (
	"context"
)

type QueueAiRequesterConsumer interface {
	Consume(ctx context.Context) (chan error, error)
}
