package eventbus

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Logger struct {
	logInfo    LogFunc
	logWarning LogFunc
	logError   LogError
}

type Options struct {
	// ConnectionString is required.
	ConnectionString string

	// OnConnectionEstablished is called after the event bus is connected or reconnected and may be nil.
	OnConnectionEstablished OnConnectFunc
}

type OnConnectFunc func()
type LogFunc func(msg string)
type LogError func(err error)

var connection *amqp.Connection
var channel *amqp.Channel
var options *Options

var logger = Logger{
	logInfo:    func(msg string) { log.Print(msg) },
	logWarning: func(msg string) { log.Print(msg) },
	logError:   func(err error) { log.Print(fmt.Sprint(err)) },
}

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
		logger.logWarning(fmt.Sprintf("Failed to connect to event bus, will try again in %s", retryTime))
		logger.logError(err)
		time.Sleep(retryTime)
		loopUntilConnected()
	}
}

func Connect(eventBusOptions Options, customLogger *Logger) {
	options = &eventBusOptions

	if customLogger != nil {
		logger = *customLogger
	}

	logger.logInfo("Event bus: connecting...")
	loopUntilConnected()
	logger.logInfo("Event bus: connected!")

	if options.OnConnectionEstablished != nil {
		options.OnConnectionEstablished()
	}

	go keepAlive()
}

func keepAlive() {
	for {
		if !IsConnected() {
			logger.logWarning("Event bus: connection to was lost!")
			Reconnect()
		}
		time.Sleep(time.Second * 5)
	}
}

func Reconnect() {
	logger.logWarning("Event bus: reconnecting...")
	loopUntilConnected()
	logger.logInfo("Event bus: reconnected!")

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
