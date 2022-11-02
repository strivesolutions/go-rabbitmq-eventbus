package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func PublishJson(exchangeName string, payload interface{}, publishTimeout time.Duration) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	return currentChannel.PublishWithContext(ctx, exchangeName, "", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
