package eventbus

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type DataRequestPayload struct {
	CorrelationId string                 `json:"correlationId"`
	Data          map[string]interface{} `json:"data"`
}

/// Data must be serializable as json
func CreateDataRequestPayload(correlationId string, data interface{}) (*DataRequestPayload, error) {
	b, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(b, &m)

	if err != nil {
		return nil, err
	}

	return &DataRequestPayload{
		CorrelationId: correlationId,
		Data:          m,
	}, nil
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
