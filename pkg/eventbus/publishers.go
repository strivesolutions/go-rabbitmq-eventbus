package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type DataRequestPayload struct {
	CorrelationId string      `json:"correlationId"`
	Data          interface{} `json:"data"`
}

func PublishJson(exchangeName ExchangeName, routingKey string, payload interface{}, publishTimeout time.Duration) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	return channel.PublishWithContext(ctx, string(exchangeName), routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
