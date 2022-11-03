package eventbus

import "fmt"

type ExchangeName string
type ExchangeType string

const (
	Fanout  ExchangeType = "fanout"
	Direct  ExchangeType = "direct"
	Topic   ExchangeType = "topic"
	Headers ExchangeType = "headers"
)

func ConfigureExchange(serviceName, exchangeName ExchangeName, t ExchangeType, durable, autoDelete, internal, noWait bool) error {
	if err := ensureConnected(); err != nil {
		return err
	}

	prefixedName := fmt.Sprintf("%s:%s", serviceName, exchangeName)

	return currentChannel.ExchangeDeclare(prefixedName, string(t), durable, autoDelete, internal, noWait, nil)
}
