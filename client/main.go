package main

// import (
// 	"fmt"
// 	"os"

// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/sakshamagrawal07/cli-chat-app.git/client/ui"
// )

// func main() {
//     p := tea.NewProgram(ui.InitialModel())
//     if _, err := p.Run(); err != nil {
//         fmt.Printf("Alas, there's been an error: %v", err)
//         os.Exit(1)
//     }
// }

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/sakshamagrawal07/cli-chat-app.git/client/cmd"
	"github.com/sakshamagrawal07/cli-chat-app.git/shared/models"
)

func main() {
	// Connect to WebSocket server
	url := "ws://localhost:8081/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Graceful shutdown on Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Ask for username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username) // remove newline and any extra whitespace

	// Goroutine to read from stdin and send to WebSocket
	go func() {
		for {
			fmt.Print("> ")
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading input:", err)
				continue
			}

			input := strings.TrimSpace(text)
			
			// Try to execute as command first
			isCommand, err := cmd.Execute(username, input, conn)
			if err != nil {
				log.Println("Command execution error:", err)
				continue
			}
			
			// If it was a command, don't send as regular message
			if isCommand {
				continue
			}

			// Regular message - send to server (default channel for now)
			msg := models.Message{
				SenderUsername:    username,
				RecipientUsername: "", // Empty for public messages
				Message:           input,
				Channel:          "general", // Default channel, you can make this configurable
			}
			err = conn.WriteJSON(msg)
			if err != nil {
				log.Println("Write error:", err)
				return
			}
		}
	}()

	// Goroutine to read from WebSocket and print
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			var msg models.Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Unmarshal error:", err)
				continue
			}
			fmt.Printf("\r%s: %s\n> ", msg.SenderUsername, msg.Message)
		}
	}()

	// Wait for Ctrl+C
	<-interrupt
	fmt.Println("\nDisconnected from server.")
}