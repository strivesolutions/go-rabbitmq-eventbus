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
	return currentConnection != nil && currentChannel != nil && !currentConnection.IsClosed() && !currentChannel.IsClosed()
}

func StatusText() string {
	if currentConnection == nil || currentChannel == nil {
		return "not configured"
	}

	return fmt.Sprintf("Connection: %s\t Channel: %s", statusString(currentConnection.IsClosed()), statusString(currentChannel.IsClosed()))
}
