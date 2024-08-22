package connections

import (
	"log/slog"
	"sync"

	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
)

var once sync.Once
var connection *rabbitmq.Connection

func ConnectQueue() (*rabbitmq.Connection, error) {
	var err error
	once.Do(func() {
		connection = &rabbitmq.Connection{}
		err = connection.Connect(
			properties.QueueConnectionUser(),
			properties.QueueConnectionPassword(),
			properties.QueueConnectionHost(),
			properties.QueueConnectionPort(),
		)
	})
	return connection, err
}

func disconnectQueue() {
	if err := connection.Disconnect(); err != nil {
		slog.Error("Error disconnecting queue",
			slog.String("error", err.Error()),
		)
	}
}
func Disconnect() {
	disconnectQueue()
}
