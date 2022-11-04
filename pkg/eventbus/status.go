package eventbus

import "fmt"

func statusString(isClosed bool) string {
	if isClosed {
		return "closed"
	} else {
		return "open"
	}
}

func IsConnected() bool {
	return connection != nil && channel != nil && !connection.IsClosed() && !channel.IsClosed()
}

func StatusText() string {
	if connection == nil || channel == nil {
		return "not configured"
	}

	return fmt.Sprintf("Connection: %s\t Channel: %s", statusString(connection.IsClosed()), statusString(channel.IsClosed()))
}
