package eventbus

import (
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConfigureFunc func()
type LogFunc func(msg string)
type LogError func(err error)

var currentConnectionString string
var currentConnection *amqp.Connection
var currentChannel *amqp.Channel
var logger *loggingConfig
var configure ConfigureFunc

type loggingConfig struct {
	logInfo    LogFunc
	logWarning LogFunc
	logError   LogError
}

func connect() error {
	if currentConnection != nil {
		Dispose()
	}

	connection, err := amqp.Dial(currentConnectionString)
	if err != nil {
		return err
	}

	currentConnection = connection

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	currentChannel = channel

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

func Connect2(connectionString string, logInfo, logWarning LogFunc, logError LogError, configureCallback func()) {
	currentConnectionString = connectionString
	configure = configureCallback

	logger = &loggingConfig{
		logInfo:    logInfo,
		logWarning: logWarning,
		logError:   logError,
	}

	logger.logInfo("Event bus: connecting...")
	loopUntilConnected()
	logger.logInfo("Event bus: connected!")
	configure()

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
	configure()
}

/// To do: remove
func Connect(connectionString string) error {
	if currentConnection != nil {
		Dispose()
	}

	if connectionString == "" {
		return errors.New("connection string is empty")
	}

	currentConnectionString = connectionString

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		return err
	}

	currentConnection = connection

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	currentChannel = channel

	return nil
}

func Dispose() {
	if currentChannel != nil {
		currentChannel.Close()
	}

	if currentConnection != nil {
		currentConnection.Close()
	}
}
