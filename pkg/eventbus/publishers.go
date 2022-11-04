package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type DataRequestPayload struct {
	CorrelationId string      `json:"correlationId"`
	Data          interface{} `json:"data"`
}

func TryParseDataPayload[T ~map[string]interface{}](m Message) (string, T, error) {
	var payload DataRequestPayload
	err := m.Json(&payload)

	if err != nil {
		return "", nil, err
	}

	if payload.Data == nil {
		return payload.CorrelationId, nil, nil
	}

	data, ok := payload.Data.(map[string]interface{})

	if !ok {
		return payload.CorrelationId, nil, errors.New("data is not a json map")
	}

	return payload.CorrelationId, data, nil

}

func ensureConnected() error {
	if !IsConnected() {
		// no loop here, only try to reconnect once when publishing
		if err := Reconnect(true); err != nil {
			return err
		}
	}
	return nil
}

func PublishJson(exchangeName ExchangeName, routingKey string, payload interface{}, publishTimeout time.Duration) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := ensureConnected(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), publishTimeout)
	defer cancel()

	return channel.PublishWithContext(ctx, string(exchangeName), routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
