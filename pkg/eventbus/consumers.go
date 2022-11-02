package eventbus

import "github.com/rabbitmq/amqp091-go"

type DeliveryFunc func(d amqp091.Delivery)

func CreateConsumer(queueName string, autoAck, exclusive, noWait bool, onDelivery DeliveryFunc) error {
	if err := ensureConnected(); err != nil {
		return err
	}

	const noLocal = false // Not supported by RabbitMQ

	delivery, err := currentChannel.Consume(queueName, "", autoAck, exclusive, noLocal, noWait, nil)

	if err != nil {
		return err
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
