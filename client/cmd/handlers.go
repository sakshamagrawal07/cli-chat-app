package cmd

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sakshamagrawal07/cli-chat-app.git/shared/models"
)

// SendChannelMessage sends a message to a specific channel
func SendChannelMessage(username, message, channel string, conn *websocket.Conn) error {
	return sendMessageToChannel(username, message, channel, conn)
}

// sendMessageToChannel sends a message to a specific channel
func sendMessageToChannel(username, message, channel string, conn *websocket.Conn) error {
	channelMsg := models.Message{
		SenderUsername:    username,
		RecipientUsername: "", // Empty for channel messages
		Message:           message,
		Channel:           channel,
	}

	err := conn.WriteJSON(channelMsg)
	if err != nil {
		return fmt.Errorf("failed to send channel message: %v", err)
	}

	fmt.Printf("Message sent to channel %s: %s\n", channel, message)
	return nil
}

// DirectMessage sends a direct message to a specific user
func DirectMessage(username, recipientUsername, message string, conn *websocket.Conn) error {
	// Create a direct message using the proper Message structure
	directMsg := models.Message{
		SenderUsername:    username,
		RecipientUsername: recipientUsername,
		Message:           message,
		Channel:           "whisper:" + getChannelName(username, recipientUsername),
	}

	// Send the message through WebSocket
	err := conn.WriteJSON(directMsg)
	if err != nil {
		return fmt.Errorf("failed to send direct message: %v", err)
	}

	fmt.Printf("Direct message sent to %s: %s\n", recipientUsername, message)
	return nil
}

func getChannelName(username1 string, username2 string) string {
	if username1 < username2 {
		return username1 + ":" + username2
	}
	return username2 + ":" + username1
}

// Future: Add encryption support
// func EncryptMessage(sender, recipient, message string) (string, error) {
//     // TODO: Implement message encryption
//     return message, nil
// }
