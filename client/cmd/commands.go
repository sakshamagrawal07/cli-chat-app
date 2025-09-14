package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"
)

// Execute processes commands with WebSocket connection
// Returns (isCommand bool, error) - isCommand indicates if input was recognized as a command
/**
 * Execute processes commands from the user input.
 * It checks for whisper commands, channel commands, and help commands.	
 * @param username The username of the user executing the command.
 * @param input The command input from the user.
*/
func Execute(username string, input string, conn *websocket.Conn) (bool, error) {
	if input == "" {
		return false, nil
	}

	// Convert to lowercase for case-insensitive matching
	lowerInput := strings.ToLower(input)

	// Check for whisper command patterns
	if strings.HasPrefix(lowerInput, "whisper ") {
		return handleWhisperCommand(username, input, conn)
	}

	// Check for channel command patterns
	if strings.HasPrefix(lowerInput, "send ") && strings.Contains(lowerInput, " to ") {
		return handleChannelCommand(username, input, conn)
	}

	// Check for help command               
	if lowerInput == "help" || lowerInput == "commands" {
		ShowHelp()
		return true, nil
	}

	// Not a recognized command
	return false, nil
}

// handleWhisperCommand processes whisper commands
// Expected format: "whisper <username> <message in quotes>"
func handleWhisperCommand(username, input string, conn *websocket.Conn) (bool, error) {
	// Regex to match: whisper <username> "message"
	re := regexp.MustCompile(`(?i)^whisper\s+(\S+)\s+"([^"]*)"$`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 3 {
		fmt.Println("Usage: whisper <username> \"your message here\"")
		return true, nil
	}

	recipientUsername := matches[1]
	message := matches[2]

	err := DirectMessage(username, recipientUsername, message, conn)
	return true, err
}

// handleChannelCommand processes channel commands
// Expected format: "send "message" to <channel>"
func handleChannelCommand(username, input string, conn *websocket.Conn) (bool, error) {
	// Regex to match: send "message" to <channel>
	re := regexp.MustCompile(`(?i)^send\s+"([^"]*)"\s+to\s+(\S+)$`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 3 {
		fmt.Println("Usage: send \"your message here\" to <channel>")
		return true, nil
	}

	message := matches[1]
	channel := matches[2]

	err := SendChannelMessage(username, message, channel, conn)
	return true, err
}

// ShowHelp displays available commands
func ShowHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  whisper <username> \"message\" - Send a direct message to a user")
	fmt.Println("  send \"message\" to <channel> - Send a message to a specific channel")
	fmt.Println("  help - Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  whisper alice \"Hey, how are you doing?\"")
	fmt.Println("  send \"Hello everyone!\" to general")
	fmt.Println("  send \"Let's discuss the new feature\" to dev")
}
