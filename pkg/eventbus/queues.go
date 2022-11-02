package eventbus

import "github.com/rabbitmq/amqp091-go"

func ConfigureQueue(queueName string, exchangeName string, routingKey string, durable, autoDelete, exclusive, noWait bool) (*amqp091.Queue, error) {
	if err := ensureConnected(); err != nil {
		return nil, err
	}

	queue, err := currentChannel.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)

	if err != nil {
		return nil, err
	}

	err = currentChannel.QueueBind(queue.Name, routingKey, exchangeName, noWait, nil)

	if err != nil {
		return nil, err
	}

	return &queue, nil
}
