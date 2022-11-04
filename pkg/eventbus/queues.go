package eventbus

import (
	"github.com/rabbitmq/amqp091-go"
)

type QueueName string

func ConfigureQueue(queueName QueueName, exchangeName ExchangeName, routingKey string, durable, autoDelete, exclusive, noWait bool) (*amqp091.Queue, error) {
	queue, err := channel.QueueDeclare(string(queueName), durable, autoDelete, exclusive, noWait, nil)

	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(queue.Name, routingKey, string(exchangeName), noWait, nil)

	if err != nil {
		return nil, err
	}

	return &queue, nil
}
