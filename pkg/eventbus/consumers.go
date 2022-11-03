package eventbus

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type DeliveryFunc func(d amqp091.Delivery)
type OnConsumerCreatedFunc func(queueName QueueName)

func CreateConsumer(
	serviceName,
	queueName QueueName,
	autoAck, exclusive,
	noWait bool,
	onDelivery DeliveryFunc,
	onConsumerCreated OnConsumerCreatedFunc,
) error {
	if err := ensureConnected(); err != nil {
		return err
	}

	prefixedName := fmt.Sprintf("%s:%s", serviceName, queueName)

	const noLocal = false // Not supported by RabbitMQ

	delivery, err := currentChannel.Consume(prefixedName, "", autoAck, exclusive, noLocal, noWait, nil)

	if err != nil {
		return err
	}

	if onConsumerCreated != nil {
		onConsumerCreated(queueName)
	}

	var forever chan struct{}
	go func() {
		for x := range delivery {
			onDelivery(x)
		}
	}()
	<-forever

	return nil
}
