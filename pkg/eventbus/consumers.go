package eventbus

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var MaxRequeueCount = 3
var requeues = make(map[string]int)

type MessageFunc func(m Message)
type OnConsumerCreatedFunc func(queueName QueueName)

type Message struct {
	delivery amqp091.Delivery
	queue    string
}

func (m Message) Queue() string {
	return m.queue
}

func (m Message) Json(v any) error {
	return json.Unmarshal(m.delivery.Body, &v)
}

func (m Message) Ack() error {
	requeues[m.queue] = 0
	return m.delivery.Ack(false)
}

func (m Message) Retry() error {
	count := requeues[m.queue]

	if count <= MaxRequeueCount {
		delay := time.Millisecond * time.Duration(count*100)
		logger.Warning(fmt.Sprintf("Requeuing message from %s in %s", m.queue, delay))
		time.Sleep(delay)
		requeues[m.queue] = count + 1
		return m.delivery.Reject(true)
	} else {
		logger.Error(fmt.Errorf("Message from %s has exceeded max requeue count (%d) and will be rejected permanently", m.queue, MaxRequeueCount))
		requeues[m.queue] = 0
		return m.delivery.Reject(false)
	}
}

func (m Message) Abort() error {
	requeues[m.queue] = 0
	return m.delivery.Reject(false)

}

func CreateConsumer(
	queueName QueueName,
	autoAck, exclusive,
	noWait bool,
	onMessage MessageFunc,
	onConsumerCreated OnConsumerCreatedFunc,
) error {
	const noLocal = false // Not supported by RabbitMQ

	delivery, err := channel.Consume(string(queueName), "", autoAck, exclusive, noLocal, noWait, nil)

	if err != nil {
		return err
	}

	if onConsumerCreated != nil {
		onConsumerCreated(queueName)
	}

	var forever chan struct{}
	go func() {
		for x := range delivery {
			message := Message{delivery: x, queue: string(queueName)}
			onMessage(message)
		}
	}()
	<-forever

	return nil
}
