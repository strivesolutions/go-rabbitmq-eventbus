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
	return channel.ExchangeDeclare(string(exchangeName), string(t), durable, autoDelete, internal, noWait, nil)
}
