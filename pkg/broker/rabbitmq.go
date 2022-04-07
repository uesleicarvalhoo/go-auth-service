package broker

import (
	"fmt"
	"strings"

	"github.com/streadway/amqp"
	"github.com/uesleicarvalhoo/go-auth-service/internal/schemas"
	"github.com/uesleicarvalhoo/go-auth-service/pkg/logger"
)

type RabbitMQClient struct {
	connection *amqp.Connection
}

func NewRabbitMqClient(cfg Config) (*RabbitMQClient, error) {
	con, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port), // nolint: nosprintfhostport
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{
		connection: con,
	}, nil
}

func (mq *RabbitMQClient) End() {
	mq.connection.Close()
}

func (mq *RabbitMQClient) Start(eventChannel <-chan schemas.Event) {
	ch, err := mq.connection.Channel()
	if err != nil {
		logger.Fatal("Couldn't start to consume events,", err)
	}

	defer ch.Close()

	for !mq.connection.IsClosed() {
		event := <-eventChannel
		exchange := strings.ToLower(fmt.Sprintf("%s.events", event.Service))

		err := ch.Publish(exchange, event.Action, false, false, amqp.Publishing{
			Body: event.Data,
		})
		if err != nil {
			logger.Info(err)
		}
	}

	logger.Info("Stopping to send messages to broker, connection was closed.")
}
