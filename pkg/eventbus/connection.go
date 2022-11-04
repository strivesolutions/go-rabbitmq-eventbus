package eventbus

import (
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OnConnectedFunc func()
type LogFunc func(msg string)

var currentConnectionString string
var currentConnection *amqp.Connection
var currentChannel *amqp.Channel
var connectionError chan error

func ensureConnected() error {
	if currentConnection == nil || currentChannel == nil {
		return Connect(currentConnectionString)
	}
	return nil
}

func statusString(isClosed bool) string {
	if isClosed {
		return "closed"
	} else {
		return "open"
	}
}

func Status() string {
	if currentConnection == nil || currentChannel == nil {
		return "not configured"
	}

	return fmt.Sprintf("Connection: %s\t Channel: %s", statusString(currentConnection.IsClosed()), statusString(currentChannel.IsClosed()))
}

// func establishConnection(connectionString string, retryTime time.Duration, infoLog, warnLog, errorLog LogFunc) {
// 	infoLog("Estabilishing connection to event bus")
// 	err := Connect(connectionString)

// 	if err != nil {
// 		warnLog(fmt.Sprintf("Failed to connect to event bus, will try again in %s", retryTime))
// 		errorLog(fmt.Sprint(err))
// 		time.Sleep(retryTime)
// 		establishConnection(connectionString, retryTime, infoLog, warnLog, errorLog)
// 	}
// }

// func Connect2(connectionString string, retryTime time.Duration, infoLog, warnLog, errorLog LogFunc, onConnected OnConnectedFunc) {
// 	establishConnection(connectionString, retryTime, infoLog, warnLog, errorLog)
// 	onConnected()
// 	return nil
// }

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

	go func() {
		<-currentConnection.NotifyClose(make(chan *amqp.Error))
		connectionError <- errors.New(("connection closed"))
	}()

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
