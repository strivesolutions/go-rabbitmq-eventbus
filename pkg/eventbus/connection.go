package eventbus

import (
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Options struct {
	// ConnectionString is required.
	ConnectionString string

	// OnConnectionEstablished is called after the event bus is connected or reconnected and may be nil.
	OnConnectionEstablished OnConnectFunc
}

type OnConnectFunc func()

var connection *amqp.Connection
var channel *amqp.Channel
var options *Options

func connect() error {
	if connection != nil {
		Dispose()
	}

	c, err := amqp.Dial(options.ConnectionString)
	if err != nil {
		return err
	}

	connection = c

	ch, err := c.Channel()
	if err != nil {
		return err
	}

	channel = ch

	return nil
}

func loopUntilConnected() {
	err := connect()

	if err != nil {
		retryTime := time.Second * 2
		logger.Warning(fmt.Sprintf("Failed to connect to event bus, will try again in %s", retryTime))
		logger.Error(err)
		time.Sleep(retryTime)
		loopUntilConnected()
	}
}

func Connect(eventBusOptions Options) error {
	if eventBusOptions.ConnectionString == "" {
		return errors.New("connection string is required")
	}

	options = &eventBusOptions

	logger.Info("Event bus: connecting...")
	loopUntilConnected()
	logger.Info("Event bus: connected!")

	if options.OnConnectionEstablished != nil {
		options.OnConnectionEstablished()
	}

	go keepAlive()

	return nil
}

func keepAlive() {
	for {
		if !IsConnected() {
			logger.Warning("Event bus: connection was lost!")
			Reconnect(false)
		}
		time.Sleep(time.Second * 5)
	}
}

func Reconnect(once bool) error {
	logger.Warning("Event bus: reconnecting...")

	if once {
		return connect()
	} else {
		loopUntilConnected()
	}

	logger.Info("Event bus: reconnected!")

	if options.OnConnectionEstablished != nil {
		options.OnConnectionEstablished()
	}

	return nil
}

func Dispose() {
	if channel != nil {
		channel.Close()
	}

	if connection != nil {
		connection.Close()
	}
}
