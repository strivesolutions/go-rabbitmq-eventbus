package eventbus

import (
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

func Connect(eventBusOptions Options, customLogger *Logger) {
	options = &eventBusOptions

	if customLogger != nil {
		logger = *customLogger
	}

	logger.Info("Event bus: connecting...")
	loopUntilConnected()
	logger.Info("Event bus: connected!")

	if options.OnConnectionEstablished != nil {
		options.OnConnectionEstablished()
	}

	go keepAlive()
}

func keepAlive() {
	for {
		if !IsConnected() {
			logger.Warning("Event bus: connection to was lost!")
			Reconnect()
		}
		time.Sleep(time.Second * 5)
	}
}

func Reconnect() {
	logger.Warning("Event bus: reconnecting...")
	loopUntilConnected()
	logger.Info("Event bus: reconnected!")

	if options.OnConnectionEstablished != nil {
		options.OnConnectionEstablished()
	}
}

func Dispose() {
	if channel != nil {
		channel.Close()
	}

	if connection != nil {
		connection.Close()
	}
}
