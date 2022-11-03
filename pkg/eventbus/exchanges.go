package eventbus

type ExchangeName string
type ExchangeType string

const (
	Fanout  ExchangeType = "fanout"
	Direct  ExchangeType = "direct"
	Topic   ExchangeType = "topic"
	Headers ExchangeType = "headers"
)

func ConfigureExchange(exchangeName ExchangeName, t ExchangeType, durable, autoDelete, internal, noWait bool) error {
	if err := ensureConnected(); err != nil {
		return err
	}

	return currentChannel.ExchangeDeclare(string(exchangeName), string(t), durable, autoDelete, internal, noWait, nil)
}
