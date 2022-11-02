package eventbus

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var currentConnectionString string
var currentConnection *amqp.Connection
var currentChannel *amqp.Channel

func ensureConnected() error {
	if currentConnection == nil || currentChannel == nil {
		return Connect(currentConnectionString)
	}
	return nil
}

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

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	currentConnection = connection
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
