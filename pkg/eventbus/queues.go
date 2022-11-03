package eventbus

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type QueueName string

func ConfigureQueue(serviceName string, queueName QueueName, exchangeName ExchangeName, routingKey string, durable, autoDelete, exclusive, noWait bool) (*amqp091.Queue, error) {
	if err := ensureConnected(); err != nil {
		return nil, err
	}

	prefixedName := fmt.Sprintf("%s/%s", serviceName, queueName)

	queue, err := currentChannel.QueueDeclare(prefixedName, durable, autoDelete, exclusive, noWait, nil)

	if err != nil {
		return nil, err
	}

	err = currentChannel.QueueBind(queue.Name, routingKey, string(exchangeName), noWait, nil)

	if err != nil {
		return nil, err
	}

	return &queue, nil
}
